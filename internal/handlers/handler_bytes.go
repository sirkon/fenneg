package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// Bytes handles []byte type.
type Bytes struct {
	lenkey string
}

// Name to satisfy Handler.
func (b *Bytes) Name(*renderer.Go) string {
	return "[]byte"
}

// Pre to satisfy Handler.
func (b *Bytes) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Len($src) + len($src)`, uniq)
	b.lenkey = uniq
}

// Len to satisfy Handler.
func (b *Bytes) Len() int {
	return -1
}

// LenExpr to satisfy Handler.
func (b *Bytes) LenExpr(r *renderer.Go, src string) string {
	return b.lenkey
}

// Encoding to satisfy Handler.
func (b *Bytes) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r.L(`$dst = $bin.AppendUvarint($dst, uint64(len($src)))`)
	r.L(`$dst = append($dst, $src...)`)
}

// Decoding to satisfy Handler.
func (b *Bytes) Decoding(r *renderer.Go, dst, src string) bool {
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
	r.L(`    if uint64(len($src)) < $siz {`)
	er.Return().New("$decode content: $recordTooSmall").LenReq("$siz").LenSrc().Rend(r)
	r.L(`    }`)
	r.L(`    $dst = $src[:$siz]`)
	r.L(`    $src = $src[$siz:]`)
	r.L(`}`)

	return true
}
