package handlers

import (
	"strconv"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// BytesArray handles [N]byte
type BytesArray int

// Name to implement TypeHandler.
func (b BytesArray) Name(r *renderer.Go) string {
	return r.S("[$0]byte", b.Len())
}

// Pre to implement TypeHandler.
func (b BytesArray) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (b BytesArray) Len() int {
	return int(b)
}

// LenExpr to implement TypeHandler.
func (b BytesArray) LenExpr(r *renderer.Go, src string) string {
	return strconv.Itoa(int(b))
}

// Encoding to implement TypeHandler.
func (b BytesArray) Encoding(r *renderer.Go, dst, src string) {
	r.L(`$dst = append($dst, $src[:]...)`)
}

// Decoding to implement TypeHandler.
func (b BytesArray) Decoding(r *renderer.Go, dst, src string) bool {
	r.L(`if len($src) < $0 {`, b.Len())
	er.Return().New("$decode: $recordTooSmall").LenReq(int(b)).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`copy($dst[:$0], $src)`, b.Len())

	return false
}
