package fenneg

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/sirkon/errors"
	"github.com/sirkon/gogh"
	"github.com/sirkon/jsonexec"
)

type runContext struct {
	fset     *token.FileSet
	handlers *TypesHandlers
	logger   LoggerType
	args     *arguments
}

type arguments struct {
	ErrorsPath packagePath `help:"Custom structured errors package." default:"github.com/sirkon/errors" short:"e"`

	OpLog   commandOpLog   `cmd:"" help:"Generate code for operation log encoding/decoding."`
	Struct  commandStruct  `cmd:"" help:"Generate code for structure flat encoding/decoding"`
	Version commandVersion `cmd:"" help:"Show version and exit."`
}

// sourcePoint represents a command line argument that looks like
// <path>:<identifier>.
type sourcePoint struct {
	Path string
	ID   string
}

func (s sourcePoint) IsValid() bool {
	return s.Path != "" && s.ID != ""
}

// UnmarshalText to implement encoding.TextUnmarshaler.
func (s *sourcePoint) UnmarshalText(text []byte) error {
	v := string(text)
	parts := strings.Split(v, ":")
	switch len(parts) {
	case 1:
		return errors.Newf("missing ':' in '%s'", v)
	case 2:
	default:
		return errors.Newf("path and identifier separated with ':', got %d separated parts instead", len(parts))
	}

	pkg, err := checkPkg(parts[0])
	if err != nil {
		return errors.Wrap(err, "check path")
	}
	if _, err := parser.ParseExpr(parts[1]); err != nil {
		return errors.Newf("invalid identifier '%s'", parts[1])
	}

	s.Path = pkg
	s.ID = parts[1]
	return nil
}

func checkPkg(pkgpath string) (string, error) {
	var dst struct {
		ImportPath string
	}
	if err := jsonexec.Run(&dst, "go", "list", "--json", pkgpath); err != nil {
		return "", errors.Wrap(err, "get package meta info")
	}

	return dst.ImportPath, nil
}

// identifier represents a command line argument that must be a correct
// Go identifier.
type identifier string

func (id *identifier) UnmarshalText(text []byte) error {
	v := string(text)

	if _, err := parser.ParseExpr(v); err != nil {
		return errors.Newf("invalid identifier '%s'", v)
	}

	if v != gogh.Public(v) {
		return errors.Newf("invalid identifier").
			Pfx("identifier").
			Str("actual", v).
			Str("would-be-acceptable", gogh.Public(v))
	}

	*id = identifier(v)
	return nil
}

func (id identifier) String() string {
	return string(id)
}

type packagePath string

func (p *packagePath) UnmarshalText(text []byte) error {
	pkg, err := checkPkg(string(text))
	if err != nil {
		return err
	}

	*p = packagePath(pkg)
	return nil
}

func (p packagePath) String() string {
	return string(p)
}
