package handlers

import (
	"go/types"

	"github.com/sirkon/gogh"
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// NewAuto creates a handler of the given type.
func NewAuto(typ types.Type) *Auto {
	return &Auto{
		typ: typ,
	}
}

// Auto handles types that passes TypeImplementsEncoder
// and TypeImplementsDecoder checks.
type Auto struct {
	lenkey string
	typ    types.Type
}

// Name to implement TypeHandler.
func (a *Auto) Name(r *renderer.Go) string {
	return r.Type(a.typ)
}

// Pre to implement TypeHandler.
func (a *Auto) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $src.Len()`, uniq)
	a.lenkey = uniq
}

// Len to implement TypeHandler.
func (a *Auto) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (a *Auto) LenExpr(r *renderer.Go, src string) string {
	return a.lenkey
}

// Encoding to implement TypeHandler.
func (a *Auto) Encoding(r *renderer.Go, dst, src string) {
	r.L(`$dst = $src.Encode($dst)`)
}

// Decoding to implement TypeHandler.
func (a *Auto) Decoding(r *renderer.Go, dst, src string) bool {
	r.Let("rest", r.Uniq("recRest"))
	r.L(`if $rest, err := $dst.Decode($src); err != nil {`)
	er.Return().Wrap("err", "$decode").Rend(r)
	r.L(`} else {`)
	r.L(`    $src = $rest`)
	r.L(`}`)

	return true
}
