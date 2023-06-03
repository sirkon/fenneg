package handlers

import (
	"strconv"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// Int8 handles int8.
func Int8() Int {
	return Int(8)
}

// Int16 handles int16.
func Int16() Int {
	return Int(16)
}

// Int32 handles int32.
func Int32() Int {
	return Int(32)
}

// Int64 handles int64.
func Int64() Int {
	return Int(64)
}

// Int handles intXXX.
type Int int

// Name to implement TypeHandler.
func (i Int) Name(*renderer.Go) string {
	return "int" + strconv.Itoa(int(i))
}

// Pre to implement TypeHandler.
func (i Int) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (i Int) Len() int {
	return i.bytes()
}

// LenExpr to implement TypeHandler.
func (i Int) LenExpr(r *renderer.Go, src string) string {
	return strconv.Itoa(i.bytes())
}

// Encoding to implement TypeHandler.
func (i Int) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r.L(`$dst = $bin.LittleEndian.AppendUint$0($dst, uint$0($src))`, i)
}

// Decoding to implement TypeHandler.
func (i Int) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")
	r.Imports().Binary().Ref("bin")

	r.L(`if len($src) < $0 {`, i.bytes())
	er.Return().New("$decode: $recordTooSmall").LenReq(i.bytes()).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`$dst = int$0($bin.LittleEndian.Uint$0($src))`, i)

	return false
}

func (i Int) bytes() int {
	return int(i >> 3)
}
