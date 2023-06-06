package fenneg

import (
	"go/token"

	"github.com/sirkon/errors"
)

// NewRunner constructs an entity to run codegen for specific
// entities right from the code.
func NewRunner(errpkg string, handlers *TypesHandlers) (*Runner, error) {
	pkgPath, err := checkPkg(errpkg)
	if err != nil {
		return nil, errors.Wrap(err, "check errors package").Str("errors-pkg-path", errpkg)
	}

	fset := token.NewFileSet()
	return &Runner{
		errors:   pkgPath,
		handlers: handlers,
		fset:     fset,
		logger:   Logger(fset),
	}, nil
}

// Runner an entity to do codegen without a final utility.
type Runner struct {
	errors   string
	handlers *TypesHandlers
	fset     *token.FileSet
	logger   LoggerType
}
