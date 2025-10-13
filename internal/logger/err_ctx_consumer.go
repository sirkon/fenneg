package logger

import (
	"bytes"
	"fmt"
	"go/token"
	"strconv"

	"github.com/sirkon/errors"
)

type errorContextConsumer struct {
	buf *bytes.Buffer
	p   string
}

func (e *errorContextConsumer) NextLink() {}

func (e *errorContextConsumer) Flt32(name string, value float32) {
	e.Float32(name, value)
}

func (e *errorContextConsumer) Flt64(name string, value float64) {
	e.Float64(name, value)
}

func (e *errorContextConsumer) Str(name string, value string) {
	//TODO implement me
	panic("implement me")
}

func (e *errorContextConsumer) SetLinkInfo(loc token.Position, descr errors.ErrorChainLinkDescriptor) {
}

func (e *errorContextConsumer) print(name string, value string) {
	e.buf.WriteString(e.p)
	e.buf.WriteString(name)
	e.buf.WriteString(": ")
	_, _ = fmt.Fprintln(e.buf, value)
}

func (e *errorContextConsumer) Bool(name string, value bool) {
	e.print(name, strconv.FormatBool(value))
}

func (e *errorContextConsumer) Int(name string, value int) {
	e.print(name, strconv.Itoa(value))
}

func (e *errorContextConsumer) Int8(name string, value int8) {
	e.print(name, strconv.FormatInt(int64(value), 10))
}

func (e *errorContextConsumer) Int16(name string, value int16) {
	e.print(name, strconv.FormatInt(int64(value), 10))
}

func (e *errorContextConsumer) Int32(name string, value int32) {
	e.print(name, strconv.FormatInt(int64(value), 10))
}

func (e *errorContextConsumer) Int64(name string, value int64) {
	e.print(name, strconv.FormatInt(value, 10))
}

func (e *errorContextConsumer) Uint(name string, value uint) {
	e.print(name, strconv.FormatUint(uint64(value), 10))
}

func (e *errorContextConsumer) Uint8(name string, value uint8) {
	e.print(name, strconv.FormatUint(uint64(value), 10))
}

func (e *errorContextConsumer) Uint16(name string, value uint16) {
	e.print(name, strconv.FormatUint(uint64(value), 10))
}

func (e *errorContextConsumer) Uint32(name string, value uint32) {
	e.print(name, strconv.FormatUint(uint64(value), 10))
}

func (e *errorContextConsumer) Uint64(name string, value uint64) {
	e.print(name, strconv.FormatUint(value, 10))
}

func (e *errorContextConsumer) Float32(name string, value float32) {
	e.print(name, strconv.FormatFloat(float64(value), 'f', 5, 32))
}

func (e *errorContextConsumer) Float64(name string, value float64) {
	e.print(name, strconv.FormatFloat(value, 'f', 5, 64))
}

func (e *errorContextConsumer) String(name string, value string) {
	e.print(name, value)
}

func (e *errorContextConsumer) Any(name string, value interface{}) {
	e.print(name, fmt.Sprintf("%#v", value))
}
