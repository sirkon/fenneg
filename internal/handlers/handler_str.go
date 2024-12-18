package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// StringHandler handles string type.
type StringHandler struct {
	lenkey string
}

// Name to satisfy Handler.
func (b *StringHandler) Name(*renderer.Go) string {
	return "string"
}

// Pre to satisfy Handler.
func (b *StringHandler) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Uint(uint(len($src))) + len($src)`, uniq)
	b.lenkey = uniq
}

// Len to satisfy Handler.
func (b *StringHandler) Len() int {
	return -1
}

// LenExpr to satisfy Handler.
func (b *StringHandler) LenExpr(r *renderer.Go, src string) string {
	return b.lenkey
}

// Encoding to satisfy Handler.
func (b *StringHandler) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r.L(`$dst = $bin.AppendUvarint($dst, uint64(len($src)))`)
	r.L(`$dst = append($dst, $src...)`)
}

// Decoding to satisfy Handler.
func (b *StringHandler) Decoding(r *renderer.Go, dst, src string) bool {
	r = r.Scope()

	off := r.Uniq("off")
	siz := r.Uniq("size")
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r = r.Scope()
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
	r.L(`    if int($siz) > len($src) {`)
	er.Return().New("$decode content: $recordTooSmall").LenReq("int($siz)").LenSrc().Rend(r)
	r.L(`    }`)
	r.L(`    $dst = string($src[:$siz])`)
	r.L(`    $src = $src[$siz:]`)
	r.L(`}`)

	return true
}
