package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// ArchVarint handles int with uleb128 zigzag encoding/decoding.
func ArchVarint() *Varint {
	chillCheck()

	return &Varint{
		bits:   0,
		lenkey: "",
	}
}

// Varint16 handles int16 with uleb128 zigzag encoding/decoding.
func Varint16() *Varint {
	return &Varint{
		bits:   16,
		lenkey: "",
	}
}

// Varint32 handles int32 with uleb128 zigzag encoding/decoding.
func Varint32() *Varint {
	return &Varint{
		bits:   32,
		lenkey: "",
	}

}

// Varint64 handles int64 with uleb128 zigzag encoding/decoding.
func Varint64() *Varint {
	return &Varint{
		bits:   64,
		lenkey: "",
	}
}

// Varint handles intXXX with uleb128 zigzag encoding/decoding.
type Varint struct {
	bits   int
	lenkey string
}

// Name to implement TypeHandler.
func (v *Varint) Name(r *renderer.Go) string {
	if v.bits == 0 {
		return "int"
	}

	return r.S(`int$0`, v.bits)
}

// Pre to implement TypeHandler.
func (v *Varint) Pre(r *renderer.Go, src string) {
	key := gogh.Private("len", src)
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Int($src)`, uniq)
	v.lenkey = uniq
}

// Len to implement TypeHandler.
func (v *Varint) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (v *Varint) LenExpr(r *renderer.Go, src string) string {
	return v.lenkey
}

// Encoding to implement TypeHandler.
func (v *Varint) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	if v.bits == 64 {
		r.L(`$dst = $bin.AppendVarint($dst, $src)`)
	} else {
		r.L(`$dst = $bin.AppendVarint($dst, int64($src))`)
	}
}

// Decoding to implement TypeHandler.
func (v *Varint) Decoding(r *renderer.Go, dst, src string) bool {
	off := r.Uniq("off")
	val := r.Uniq("val", dst)
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r.Let(off, off)

	r.L(`{`)

	if v.bits == 64 {
		r.L(`    var $off int`)
		r.L(`    $dst, $off = $bin.Varint($src)`)
	} else {
		r.Let(val, val)
		r.L(`    $val, $off := $bin.Varint($src)`)
	}
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode: $malformedVarint").Rend(r)
	r.L(`    }`)
	if v.bits != 64 {
		r.L(`$dst = $0($val)`, v.Name(r))
	}
	r.L(`    $src = $src[$off:]`)
	r.L(`}`)

	return true
}
