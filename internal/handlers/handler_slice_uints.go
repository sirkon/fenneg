package handlers

import (
	"strconv"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

func NewSliceUint16() *SliceUint {
	return &SliceUint{bits: 16}
}

func NewSliceUint32() *SliceUint {
	return &SliceUint{bits: 32}
}

func NewSliceUint64() *SliceUint {
	return &SliceUint{bits: 64}
}

type SliceUint struct {
	bits   int
	lenkey string
}

// Name to implement TypeHandler.
func (i *SliceUint) Name(*renderer.Go) string {
	return "[]uint" + strconv.Itoa(i.bits)
}

// Pre to implement TypeHandler.
func (i *SliceUint) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Len($src) + len($src) * $1`, uniq, i.elemBytes())
	i.lenkey = uniq
}

// Len to implement TypeHandler.
func (i *SliceUint) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (i *SliceUint) LenExpr(r *renderer.Go, src string) string {
	return i.lenkey
}

// Encoding to implement TypeHandler.
func (i *SliceUint) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r.L(`$dst = $bin.AppendUvarint($dst, uint64(len($src)))`)
	r.L(`for i := range $src {`)
	r.L(`$dst = $bin.LittleEndian.AppendUint$0($dst, uint$0($src[i]))`, i.bits)
	r.L(`}`)
}

// Decoding to implement TypeHandler.
func (i *SliceUint) Decoding(r *renderer.Go, dst, src string) bool {
	r.Scope()

	off := r.Uniq("off")
	siz := r.Uniq("size")
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r = r.Scope()
	r.Let("siz", siz)
	r.Let("off", off)
	r.Let("bits", i.bits)
	r.Let("len", i.elemBytes())

	r.L(`{`)
	r.L(`    $siz, $off := $bin.Uvarint($src)`)
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode length: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode length: $malformedUvarint").Rend(r)
	r.L(`    }`)
	r.L(`    $src = $src[$off:]`)
	r.L(`    if int($siz)*$len > len($src) {`)
	er.Return().New("$decode content: $recordTooSmall").LenReq("int($siz)").LenSrc().Rend(r)
	r.L(`    }`)
	r.L(`    if $siz > 0 {`)
	r.L(`        $dst = make($0, 0, $siz)`, i.Name(r))
	r.L(`        for range $siz {`)
	r.L(`            $dst = append($dst, uint$bits($bin.LittleEndian.Uint$bits($src)))`)
	r.L(`            $src = $src[$len:]`)
	r.L(`        }`)
	r.L(`    }`)
	r.L(`}`)

	return true
}

func (i *SliceUint) elemBytes() int {
	return i.bits >> 3
}
