package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// Bool tuples bool type.
type Bool struct{}

// Name to implement TypeHandler.
func (Bool) Name(*renderer.Go) string {
	return "bool"
}

// Pre to implement TypeHandler.
func (Bool) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (Bool) Len() int {
	return 1
}

// LenExpr to implement handler
func (Bool) LenExpr(r *renderer.Go, src string) string {
	return "1"
}

// Encoding to implement TypeHandler.
func (Bool) Encoding(r *renderer.Go, dst, src string) {
	r.L(`if $src {`)
	r.L(`    $dst = append($dst, 1)`)
	r.L(`} else {`)
	r.L(`    $dst = append($dst, 0)`)
	r.L(`}`)
}

// Decoding to implement TypeHandler.
func (Bool) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")
	r.L(`if len($src) < 1 {`)
	er.Return().New("$decode: $recordTooSmall").LenReq(1).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`if $src[0] != 0 {`)
	r.L(`    $dst = true`)
	r.L(`}`)

	return false
}
