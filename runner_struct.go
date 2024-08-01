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

type options struct {
	FileSuffix    string
	StructPointer bool
	EncoderName   string
	DecoderName   string
	SizeName      string
}

type Option func(opts *options)

func WithFileSuffix(f string) Option {
	return func(opt *options) {
		if f != "" {
			opt.FileSuffix = f
		}
	}
}

func StructPointer() Option {
	return func(opt *options) {
		opt.StructPointer = true
	}
}

func WithEncoderName(e string) Option {
	return func(opt *options) {
		if e != "" {
			opt.EncoderName = e
		}
	}
}

func WithDecoderName(d string) Option {
	return func(opt *options) {
		if d != "" {
			opt.DecoderName = d
		}
	}
}

func WithSizeName(siz string) Option {
	return func(opt *options) {
		if siz != "" {
			opt.SizeName = siz
		}
	}
}

// Struct run struct generator.
func (r *Runner) Struct(pkg, typ string, optsFn ...Option) error {
	opts := &options{
		FileSuffix:  "generated",
		EncoderName: "Encode",
		DecoderName: "Decode",
		SizeName:    "Len",
	}
	for _, opt := range optsFn {
		opt(opts)
	}

	loader := newSouceLoader(r.fset)
	p, err := loader.loadPkg(pkg)
	if err != nil {
		return errors.Wrap(err, "load package").Str("package-path", pkg)
	}

	t := p.Types.Scope().Lookup(typ)
	if t == nil {
		return errors.New("unknown type").
			Str("package-path", pkg).
			Str("type-name", typ)
	}

	log := r.logger

	tn, ok := t.Type().(*types.Named)
	if !ok {
		log.Pos(
			t.Pos(),
			errors.New(typ+" is not a type").
				Str("package-path", pkg).
				Str("type-name", typ),
		)
		return errorCannotContinue
	}

	s, ok := tn.Underlying().(*types.Struct)
	if !ok {
		log.Pos(
			t.Pos(),
			errors.New(typ+" is not a struct").
				Str("package-path", pkg).
				Str("type-name", typ),
		)
		return errorCannotContinue
	}

	if !r.checkStruct(s, p, loader) {
		return errorCannotContinue
	}

	prj, err := gogh.New(gogh.FancyFmt, renderer.NewImports)
	if err != nil {
		return errors.Wrap(err, "set up the code renderer")
	}
	rpkg, err := prj.Package(p.Name, p.PkgPath)
	if err != nil {
		return errors.Wrap(err, "set up rendering package")
	}
	_, file := filepath.Split(r.fset.Position(tn.Obj().Pos()).Filename)
	rr := rpkg.Go(strings.TrimSuffix(file, ".go")+"_"+strings.ToLower(opts.FileSuffix)+".go", gogh.Autogen(app.Name))

	manhands := map[*types.Var]handlers.Type{}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		manhands[f] = r.handlers.Handler(f)
	}
	g := generator.NewStruct(
		rr, tn, manhands, opts.StructPointer, opts.EncoderName, opts.DecoderName, opts.SizeName,
	)

	g.Generate()

	if err := prj.Render(); err != nil {
		return errors.Wrap(err, "render generated source code")
	}

	return nil
}

func (r *Runner) checkStruct(s *types.Struct, p *packages.Package, loader *souceLoader) bool {
	log := r.logger
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

		if r.handlers.Handler(f) == nil {
			log.Pos(f.Pos(), errors.New("unsupported type").Stg("unsupported-type", f.Type()))
			success = false
			continue
		}

		aliases[f] = loader.digForAliases(p, f)
	}
	r.handlers.setArgsAliases(aliases)

	return success
}
