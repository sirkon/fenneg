package fenneg

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"sync"

	"github.com/sirkon/errors"
	"github.com/sirkon/gogh"
	"github.com/sirkon/jsonexec"
	"golang.org/x/tools/go/packages"

	"github.com/sirkon/fenneg/internal/tdetect"
)

type typeGetter func(p *packages.Package, name string) (*types.Named, error)

// typeByName gets type by name.
func typeByName(pkg *packages.Package, name string) (*types.Named, error) {
	s := pkg.Types.Scope().Lookup(name)
	if s == nil {
		return nil, errors.New("type not found")
	}

	v, ok := s.Type().(*types.Named)
	if !ok {
		return nil,
			errors.Newf("%s is not a type", name).
				Pfx("node-type").
				Type("expected", new(types.Named)).
				Type("actual", s.Type())
	}

	return v, nil
}

func typeByNameSpecific[T types.Type](pkg *packages.Package, name string) (*types.Named, error) {
	typ, err := typeByName(pkg, name)
	if err != nil {
		return nil, errors.Wrap(err, "look for the type")
	}

	var tmp T
	if _, ok := typ.Underlying().(T); !ok {
		return nil, errors.New("unexpected type kind").
			Pfx("type-kind").
			Type("expected", tmp).
			Type("actual", typ.Underlying())
	}

	return typ, nil
}

func newSouceLoader(fset *token.FileSet) *souceLoader {
	return &souceLoader{
		fset:     fset,
		pkgCache: map[string]*packages.Package{},
	}
}

type souceLoader struct {
	fset *token.FileSet
	lock sync.Mutex

	pkgCache map[string]*packages.Package
}

// mapExactTypes compute exact types of arguments of the interface methods.
func (l *souceLoader) mapExactTypes(iface *types.Named) map[*types.Var]*types.Named {
	res := map[*types.Var]*types.Named{}

	// Need to find AST of a file iface is defined in.
	var pkg *packages.Package
	for imp, p := range l.pkgCache {
		if imp == iface.Obj().Pkg().Path() {
			pkg = p
			break
		}
	}

	t := iface.Underlying().(*types.Interface)
	for i := 0; i < t.NumMethods(); i++ {
		m := t.Method(i)
		ps := m.Type().(*types.Signature).Params()
		for i := 0; i < ps.Len(); i++ {
			p := ps.At(i)
			at := l.digForAliases(pkg, p)
			if at != nil {
				res[p] = at
			}
		}
	}

	return res
}

func (l *souceLoader) digForAliases(pkg *packages.Package, p *types.Var) *types.Named {
	var file *ast.File
	for _, sntx := range pkg.Syntax {
		// Look for a file with this parameter.
		if isWithin(p.Pos(), sntx) {
			file = sntx
			break
		}
	}

	if file == nil {
		panic(pkg.Fset.Position(p.Pos()).String() + " failed to this file")
	}

	// Check if param type means the same in AST and ATT.
	imports := map[string]string{}
	var srcPkg string
	var srcName string
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ImportSpec:
			// Need to know package names.
			if n.Name != nil {
				imports[n.Name.Name] = n.Path.Value
			} else {
				pval, _ := strconv.Unquote(n.Path.Value)
				ipkg := pkg.Imports[pval]
				imports[ipkg.Name] = pval
			}

		case *ast.Field:
			var found bool
			for _, name := range n.Names {
				if name.Pos() == p.Pos() {
					found = true
					break
				}
			}
			if !found {
				return true
			}

			switch v := n.Type.(type) {
			case *ast.Ident:
				if v.Name == p.Type().String() {
					// Exact match for names, this is the same builtin type.
					return false
				}
				if pkg.PkgPath+"."+v.Name == p.Type().String() {
					// Exact match for names, this is the same user defined type –
					// they have package path prefix in their names.
					return false
				}

				// Names are different. This means there's an alias
				// in the actual source, so we need to check it in the
				// current package, since it is just an Indent rather
				// than a SelectorExpr which the other package's type
				// would have.
				srcName = v.Name

			case *ast.SelectorExpr:
				name := v.Sel.Name
				pkgident := v.X.(*ast.Ident).Name
				pkgpath := imports[pkgident]

				op, ok := p.Type().(*types.Named)
				if !ok {
					// This is an alias exactly since it is a type
					// from a different package, and it must be a
					// named type for a type referenced as a selector
					// expression.
					// Look for the name in the pkgident package then.
					srcPkg = pkgpath
					srcName = name
					return false
				}

				if op.Obj().Pkg().Path() == pkgpath && op.Obj().Name() == name {
					// This is the same name, no aliasing here. Just stop.
					return false
				}

				// It means this is a different type. Look for it.
				srcPkg = pkgpath
				srcPkg = pkgpath
				return false
			}
			return false

		default:
			return true
		}

		return false
	})

	if srcName == "" {
		return nil
	}

	if srcPkg != "" {
		pkg = pkg.Imports[srcPkg]
	}

	res := pkg.Types.Scope().Lookup(srcName).Type()
	switch v := res.(type) {
	case *types.Named:
		var methods []*types.Func
		for i := 0; i < v.NumMethods(); i++ {
			methods = append(methods, v.Method(i))
		}
		return types.NewNamed(
			types.NewTypeName(token.NoPos, pkg.Types, srcName, res),
			v.Underlying(),
			methods,
		)
	default:
		return types.NewNamed(
			types.NewTypeName(token.NoPos, pkg.Types, srcName, res),
			res,
			nil,
		)
	}
}

func (l *souceLoader) loadPkg(pkgpath string) (pkg *packages.Package, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if v, ok := l.pkgCache[pkgpath]; ok {
		return v, nil
	}

	defer func() {
		if err != nil {
			return
		}

		l.pkgCache[pkgpath] = pkg
	}()

	var dst struct {
		ImportPath string
	}
	if err := jsonexec.Run(&dst, "go", "list", "--json", pkgpath); err != nil {
		return nil, errors.Wrap(err, "get package meta info").Str("pkg-path", pkgpath)
	}

	defer func() {
		if err == nil {
			return
		}

		err = errors.Just(err).Str("package-path", dst.ImportPath)
	}()

	mode := packages.NeedImports | packages.NeedTypes | packages.NeedName |
		packages.NeedDeps | packages.NeedSyntax | packages.NeedFiles | packages.NeedModule

	pkgs, err := packages.Load(
		&packages.Config{
			Mode:  mode,
			Fset:  l.fset,
			Tests: false,
		},
		dst.ImportPath,
	)
	if err != nil {
		return nil, errors.Wrap(err, "load package source")
	}

	for _, pkg := range pkgs {
		if pkg.PkgPath == dst.ImportPath {
			return pkg, nil
		}
	}

	return nil, errors.New("failed to load package source")
}

func validateSourceInterface(logger LoggerType, iface *types.Named) bool {
	noerrs := true

	name := iface.Obj().Name()
	if name != gogh.Public(name) {
		// I love them canonical and so you will too LMAO.
		logger.Pos(
			iface.Obj().Pos(),
			errors.New("invalid interface name").
				Pfx("interface-name").
				Str("actual", name).
				Str("would-be-valid", gogh.Public(name)),
		)
		noerrs = false
	}

	def := iface.Underlying().(*types.Interface)
	for i := 0; i < def.NumMethods(); i++ {
		m := def.Method(i)

		if m.Name() != gogh.Public(m.Name()) {
			logger.Pos(
				m.Pos(),
				errors.New("invalid method name").
					Pfx("interface-method-name").
					Str("actual", m.Name()).
					Str("would-be-valid", gogh.Public(m.Name())),
			)
			noerrs = false
		}

		checkRes := true
		s := m.Type().(*types.Signature)
		switch s.Results().Len() {
		case 0:
			logger.Pos(
				m.Pos(),
				errors.New("any method is required to return an error"),
			)
			checkRes = false
		case 1:
			rt := s.Results().At(0)
			if rt.Name() != "" {
				checkRes = false
				logger.Pos(
					rt.Pos(),
					errors.New("the result value must not have a name"),
				)
			}
			if _, ok := rt.Type().(*types.Named); !ok {
				checkRes = false
				logger.Pos(
					rt.Pos(),
					errors.New("the result type must be the error"),
				)
			}
		default:
			logger.Pos(
				s.Results().At(1).Pos(),
				errors.New("any method is required to return one value exactly"),
			)
			checkRes = false
		}

		noerrs = validateMethod(logger, s, checkRes) && noerrs
	}

	return noerrs
}

// validateMethod won't return true if checkRes is false.
func validateMethod(logger LoggerType, s *types.Signature, checkRes bool) bool {
	noerrs := true

	for i := 0; i < s.Params().Len(); i++ {
		p := s.Params().At(i)
		switch {
		case p.Name() == "":
			logger.Pos(
				p.Pos(),
				errors.New("missing parameter name"),
			)
			noerrs = false
		case p.Name() == "err":
			logger.Pos(
				p.Pos(),
				errors.New("err name is not allowed for parameters"),
			)
			noerrs = false
		case p.Name() != gogh.Private(p.Name()):
			logger.Pos(
				p.Pos(),
				errors.New("invalid parameter name").
					Pfx("method-parameter-name").
					Str("actual", p.Name()).
					Str("would-be-valid", gogh.Private(p.Name())),
			)
			noerrs = false
		}

		// Get rid of a pointer and strip aliases.
		tp := p.Type()
		if v, ok := tp.(*types.Pointer); ok {
			tp = v.Elem()
		}
	loop:
		for {
			switch v := tp.(type) {
			case *types.Named:
				if v.Obj().IsAlias() {
					tp = v.Underlying()
				}
			default:
				break loop
			}
		}

		switch v := p.Type().(type) {
		case *types.Basic:
			switch v.Kind() {
			case types.Complex64, types.Complex128:
				logger.Pos(
					p.Pos(),
					errors.New("complex types are not supported"),
				)
			case types.Int:
				logger.Pos(p.Pos(), errorPlatformDependent("int"))
			case types.Uint:
				logger.Pos(p.Pos(), errorPlatformDependent("uint"))
			case types.Uintptr:
				logger.Pos(p.Pos(), errorPlatformDependent("uintptr"))
			}
		case *types.Named:
			if _, ok := v.Obj().Type().(*types.Interface); ok {
				logger.Pos(
					p.Pos(),
					errorMeaningful("interfaces are not known for this"),
				)
				noerrs = false

			}
		case *types.Pointer:
			logger.Pos(
				p.Pos(),
				errorMeaningful("a pointer of a pointer cannot do this"),
			)
			noerrs = false
		case *types.Chan:
			logger.Pos(
				p.Pos(),
				errorMeaningful("a channel won't do this"),
			)
		default:
		}
	}

	if !checkRes {
		return false
	}

	res := s.Results().At(0)
	if !tdetect.IsErrorType(res.Type()) {
		logger.Pos(
			res.Pos(),
			errors.New("the result type must be the error"),
		)
	}

	return noerrs
}

func errorPlatformDependent(what string) *errors.Error {
	return errors.Newf("type %s is platform dependent and not supported", what)
}

func errorMeaningful(descr string) *errors.Error {
	return errors.Newf("types are required to have meaningful zero value and %s", descr)
}

func isWithin(needle token.Pos, hay ast.Node) bool {
	return hay.Pos() <= needle && needle < hay.End()
}
