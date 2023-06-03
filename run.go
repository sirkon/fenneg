package fenneg

import (
	"go/token"
	"go/types"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/sirkon/errors"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
	"golang.org/x/tools/go/packages"
)

// Run is the entry point of a final utility.
func Run(
	appName string,
	handlers *TypesHandlers,
) error {
	var args arguments
	cliParser := kong.Must(
		&args,
		kong.Name(appName),
		kong.Description("Generate records encoder and encoding dispatcher for the given interface."),
		kong.UsageOnError(),
	)

	ctx, err := cliParser.Parse(os.Args[1:])
	if err != nil {
		cliParser.FatalIfErrorf(err)
	}

	renderer.SetStructuredErrorsPkgPath(string(args.ErrorsPath))

	runCtx := &runContext{
		fset:     handlers.fs(),
		handlers: handlers,
		logger:   Logger(handlers.fs()),
		args:     &args,
	}
	if err := ctx.Run(runCtx); err != nil {
		return errors.Wrap(err, "run command")
	}

	return nil
}

func getRenderer(appName string, fset *token.FileSet, typ *types.Named) (*Project, *Go, error) {
	prj, err := gogh.New(gogh.FancyFmt, renderer.NewImports)
	if err != nil {
		return nil, nil, err
	}

	pkg, err := prj.Package("", typ.Obj().Pkg().Path())
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	pos := fset.Position(typ.Obj().Pos())
	_, filename := path.Split(pos.Filename)
	destfile := strings.TrimSuffix(filename, ".go") + "_generated.go"
	r := pkg.Go(destfile, gogh.Autogen(appName))

	return prj, r, nil
}

func getType(l *souceLoader, p sourcePoint, g typeGetter, what string) (*packages.Package, *types.Named, error) {
	pkg, err := l.loadPkg(p.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "load package").Pfx(what).Str("package-name", p.Path)
	}

	res, err := g(pkg, p.ID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "load type by name").Pfx(what).
			Str("package-name", p.Path).
			Str("type-name", p.ID)
	}

	return pkg, res, nil
}
