package handlers

import (
	"go/types"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// SlicesUniform []T types support for supported types T with fixed length encoding.
type SlicesUniform struct {
	handler  Type
	itemType types.Type

	lenkey string
}

// NewSlicesUniform construct SliceUniform handler.
func NewSlicesUniform(handler Type, itemType types.Type) *SlicesUniform {
	return &SlicesUniform{
		handler:  handler,
		itemType: itemType,
	}
}

// Name to implement TypeHandler.
func (s *SlicesUniform) Name(r *renderer.Go) string {
	return "[]" + s.handler.Name(r)
}

// Pre to implement TypeHandler.
func (s *SlicesUniform) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Len($src) + len($src)*$1`, uniq, s.handler.Len())
	s.lenkey = uniq
}

// Len to implement TypeHandler.
func (s *SlicesUniform) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (s *SlicesUniform) LenExpr(r *renderer.Go, src string) string {
	return s.lenkey
}

// Encoding to implement TypeHandler.
func (s *SlicesUniform) Encoding(r *renderer.Go, dst, src string) {
	r = r.Scope()

	v := r.Uniq("v", src)
	r.L(`for _, $0 := range $src {`, v)
	r.Let("src", v)
	s.handler.Encoding(r, dst, v)
	r.L(`}`)
}

// Decoding to implement TypeHandler.
func (s *SlicesUniform) Decoding(r *renderer.Go, dst, src string) bool {
	r = r.Scope()

	off := r.Uniq("off")
	siz := r.Uniq("size")
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r.Let("siz", siz)
	r.Let("off", off)

	r.L(`{`)
	r.L(`    $siz, $off := $bin.Uvarint($src)`)
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode length: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode length: $malformedUvarint").Rend(r)
	r.L(`    }`)
	r.L(`    $src = $src[$off:]`)
	r.L(`    $dst = make([]$0, $siz)`, r.Type(s.itemType))
	it := r.Uniq("iter", dst)
	r.Let("dst", r.S("$dst[$0]", it))
	r.L(`    for $0 := 0; $0 < int($siz); $0++ {`, it)
	s.handler.Decoding(r, dst, src)
	r.L(`        $src = $src[$0:]`, s.handler.LenExpr(r, src))
	r.L(`    }`)
	r.L(`}`)

	return true
}
