package logger

import (
	"bytes"
	"fmt"
	"go/token"
	"runtime"
)

// NewTesting creates log for unit testing.
func NewTesting(t TestingMeans, fset *token.FileSet) *TestingLogger {
	logger := &TestingLogger{
		l:    New(fset),
		dest: t,
	}
	logger.l.head = "\r"
	return logger
}

// TestingMeans to hide testing.XXX.
type TestingMeans interface {
	Log(...any)
	Error(...any)
}

// TestingLogger logger for testing.
type TestingLogger struct {
	l    *Logger
	dest TestingMeans
}

// Err log error.
func (l *TestingLogger) Err() TestingPerformer {
	return TestingPerformer{
		l: l,
		w: l.dest.Error,
	}
}

// Wrn log warning.
func (l *TestingLogger) Wrn() TestingPerformer {
	return TestingPerformer{
		l: l,
		w: l.dest.Log,
	}
}

// TestingPerformer handles actual logging.
type TestingPerformer struct {
	l *TestingLogger
	w func(...any)
}

// Pos logs with positional information.
func (p TestingPerformer) Pos(pos token.Pos, err error) {
	var buf bytes.Buffer
	p.l.l.log(&buf, p.l.l.fset.Position(pos).String(), err)
	p.w(buf.String())
}

// Err generic loggin.
func (p TestingPerformer) Error(err error) {
	_, file, line, ok := runtime.Caller(1)
	var prefix string
	if ok {
		prefix = fmt.Sprintf("%s:%d", file, line)
	}

	var buf bytes.Buffer
	p.l.l.log(&buf, prefix, err)
	p.w(buf.String())
}
