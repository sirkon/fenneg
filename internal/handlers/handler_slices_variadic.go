package handlers

import (
	"go/types"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// SlicesVariadic []T types support for supported types T with variable lengths encoding.
type SlicesVariadic struct {
	handler  Type
	itemType types.Type

	lenkey string
}

// NewSlicesVariadic constructs SlicesVariadic handler.
func NewSlicesVariadic(handler Type, itemType types.Type) *SlicesVariadic {
	return &SlicesVariadic{
		handler:  handler,
		itemType: itemType,
	}
}

// Name to implement TypeHandler.
func (s *SlicesVariadic) Name(r *renderer.Go) string {
	return "[]" + s.handler.Name(r)
}

// Pre to implement TypeHandler.
func (s *SlicesVariadic) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	s.lenkey = uniq

	item := r.Uniq("item", src)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Len($src)`, uniq)
	r.L(`for _, $0 := range $src { `, item)
	r = r.Scope()
	r.Let("src", item)
	s.handler.Pre(r, item)
	r.L(`    $0 += $1`, uniq, s.handler.LenExpr(r, item))
	r.L(`}`)
}

// Len to implement TypeHandler.
func (s *SlicesVariadic) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (s *SlicesVariadic) LenExpr(r *renderer.Go, src string) string {
	return s.lenkey
}

// Encoding to implement TypeHandler.
func (s *SlicesVariadic) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r = r.Scope()

	r.L(`$dst = $bin.AppendUvarint($dst, uint64(len($src)))`)

	v := r.Uniq("v", src)
	r.L(`for _, $0 := range $src {`, v)
	r.Let("src", v)
	s.handler.Encoding(r, dst, v)
	r.L(`}`)
}

// Decoding to implement TypeHandler.
func (s *SlicesVariadic) Decoding(r *renderer.Go, dst, src string) bool {
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
	it := r.Uniq("i")
	r.Let("dst", r.S("$dst[$0]", it))
	r.L(`    for $0 := 0; $0 < int($siz); $0++ {`, it)
	s.handler.Decoding(r, dst, src)
	r.L(`    }`)
	r.L(`}`)

	return true
}
