package logger

import (
	"bytes"
	"go/token"
	"io"
	"os"

	"github.com/sirkon/errors"
)

const contextIndent = "        "

// Type an abstraction over a logger.
type Type interface {
	Pos(pos token.Pos, err error)
	Error(err error)
}

// New constructs Logger.
func New(fset *token.FileSet) *Logger {
	return &Logger{
		fset: fset,
	}
}

// Logger generic workload logger with positions.
type Logger struct {
	fset *token.FileSet

	head string
}

// Pos log errors with the given source position.
func (l *Logger) Pos(pos token.Pos, err error) {
	var buf bytes.Buffer
	l.log(&buf, l.fset.Position(pos).String(), err)
	_, _ = io.Copy(os.Stderr, &buf)
}

// Error prints generic errors.
func (l *Logger) Error(err error) {
	var buf bytes.Buffer
	l.log(&buf, "", err)
	_, _ = io.Copy(os.Stderr, &buf)
}

func (l *Logger) log(buf *bytes.Buffer, prefix string, err error) {
	buf.WriteString(l.head)
	if len(prefix) > 0 {
		buf.WriteString(prefix)
		buf.WriteByte(' ')
	}
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	del := errors.GetContextDeliverer(err)
	if del == nil {
		return
	}

	cons := errorContextConsumer{
		buf: buf,
		p:   contextIndent,
	}
	del.Deliver(&cons)
}

var _ Type = new(Logger)
