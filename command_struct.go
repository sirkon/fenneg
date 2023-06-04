package fenneg

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/sirkon/errors"
	"github.com/sirkon/fenneg/internal/app"
	"github.com/sirkon/fenneg/internal/generator"
	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
	"golang.org/x/tools/go/packages"
)

const errorCannotContinue errors.Const = "cannot continue"

// commandStruct handles structure use case.
type commandStruct struct {
	Source sourcePoint `arg:"" help:"Source structure <path>:<struct> to generate code for its encoding/decoding."`
}

// Run command.
func (c *commandStruct) Run(ctx *runContext) error {
	log := ctx.logger
	l := newSouceLoader(ctx.fset)

	pkg, err := l.loadPkg(c.Source.Path)
	if err != nil {
		return errors.Wrap(err, "load package").Str("package-path", c.Source.Path)
	}

	t := pkg.Types.Scope().Lookup(c.Source.ID)
	if t == nil {
		return errors.New("unknown type").
			Str("package-path", c.Source.Path).
			Str("type-name", c.Source.ID)
	}

	tn, ok := t.Type().(*types.Named)
	if !ok {
		log.Pos(
			t.Pos(),
			errors.New(c.Source.ID+" is not a type").
				Str("package-path", c.Source.Path).
				Str("type-name", c.Source.ID),
		)
		return errorCannotContinue
	}

	s, ok := tn.Underlying().(*types.Struct)
	if !ok {
		log.Pos(
			t.Pos(),
			errors.New(c.Source.ID+" is not a type").
				Str("package-path", c.Source.Path).
				Str("type-name", c.Source.ID),
		)
		return errorCannotContinue
	}

	if !c.checkStruct(ctx, s, l, pkg) {
		return errors.New("structure is not supported")
	}

	prj, err := gogh.New(gogh.FancyFmt, renderer.NewImports)
	if err != nil {
		return errors.Wrap(err, "set up the code renderer")
	}
	p, err := prj.Package(pkg.Name, pkg.PkgPath)
	if err != nil {
		return errors.Wrap(err, "set up rendering package")
	}
	_, file := filepath.Split(ctx.fset.Position(tn.Obj().Pos()).Filename)
	r := p.Go(strings.TrimSuffix(file, ".go")+"_generated.go", gogh.Autogen(app.Name))

	manhands := map[*types.Var]handlers.Type{}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		manhands[f] = ctx.handlers.Handler(f)
	}
	g := generator.NewStruct(r, tn, manhands)

	g.Generate()

	if err := prj.Render(); err != nil {
		return errors.Wrap(err, "render generated source code")
	}

	return nil
}

func (c *commandStruct) checkStruct(
	ctx *runContext,
	s *types.Struct,
	l *souceLoader,
	pkg *packages.Package,
) bool {
	log := ctx.logger
	success := true

	aliases := map[*types.Var]*types.Named{}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		if !f.Exported() {
			log.Pos(f.Pos(), errors.New("unexported fields are not supported").Str("field-name", f.Name()))
			success = false
			continue
		}

		if f.Anonymous() {
			log.Pos(f.Pos(), errors.New("embedded fields are not supported"))
			success = false
			continue
		}

		if ctx.handlers.Handler(f) == nil {
			log.Pos(f.Pos(), errors.New("unsupported type").Stg("unsupported-type", f.Type()))
			success = false
			continue
		}

		aliases[f] = l.digForAliases(pkg, f)
	}
	ctx.handlers.setArgsAliases(aliases)

	return success
}
