package handlers

import (
	"strconv"

	"github.com/sirkon/gogh"
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// Float32 handles float32.
func Float32() Float {
	return Float(32)
}

// Float64 handles float64.
func Float64() Float {
	return Float(64)
}

// Float handles floatX.
type Float int

// Name to implement TypeHandler.
func (f Float) Name(*renderer.Go) string {
	return "float" + strconv.Itoa(int(f))
}

// Pre to implement TypeHandler.
func (f Float) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (f Float) Len() int {
	return f.bytes()
}

// LenExpr to implement TypeHandler.
func (f Float) LenExpr(r *renderer.Go, src string) string {
	return strconv.Itoa(f.bytes())
}

// Encoding to implement TypeHandler.
func (f Float) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Errors().Ref("errors")
	r.Imports().Add("math").Ref("math")
	r.Imports().Binary().As("bin")
	key := r.Uniq(gogh.Private("key", src))

	r.L(`$dst = $bin.LittleEndian.AppendUint$0($dst, math.Float$0bits($src))`, f, key)
}

// Decoding to implement TypeHandler.
func (f Float) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")
	r.Imports().Add("math").Ref("math")
	r.Imports().Binary().As("bin")
	key := r.Uniq(gogh.Private("key", src))

	r.L(`if len($src) >= $0 {`, f.bytes())
	r.L(`    $0 := $bin.LittleEndian.Uint$1($src)`, key, f)
	r.L(`    $dst = $math.Float$0frombits($1)`, f, key)
	r.L(`} else {`)
	er.Return().New("$decode: $recordTooSmall").LenReq(f.bytes()).LenSrc().Rend(r)
	r.L(`}`)

	return false
}

func (f Float) bytes() int {
	return int(f) >> 3
}
