package fenneg

import (
	"go/token"

	"github.com/sirkon/fenneg/internal/logger"
)

// LoggerType error logging abstraction to log with text position.
type LoggerType interface {
	Pos(pos token.Pos, err error)
	Error(err error)
}

// Logger returns a standard LoggerType implementation that
// seems to be sufficient for most needs.
func Logger(fset *token.FileSet) LoggerType {
	return logger.New(fset)
}

var (
	_ logger.Type = LoggerType(nil)
)
