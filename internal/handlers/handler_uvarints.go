package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// ArchUvarint handles uint with uleb128 zigzag encoding/decoding.
func ArchUvarint() *Uvarint {
	chillCheck()

	return &Uvarint{
		bits:   0,
		lenkey: "",
	}
}

// Uvarint16 handles uint16 with uleb128 zigzag encoding/decoding.
func Uvarint16() *Uvarint {
	return &Uvarint{
		bits:   16,
		lenkey: "",
	}
}

// Uvarint32 handles uint32 with uleb128 zigzag encoding/decoding.
func Uvarint32() *Uvarint {
	return &Uvarint{
		bits:   32,
		lenkey: "",
	}

}

// Uvarint64 handles uint64 with uleb128 zigzag encoding/decoding.
func Uvarint64() *Uvarint {
	return &Uvarint{
		bits:   64,
		lenkey: "",
	}
}

// Uvarint handles uintXXX with uleb128 zigzag encoding/decoding.
type Uvarint struct {
	bits   int
	lenkey string
}

// Name to implement TypeHandler.
func (v *Uvarint) Name(r *renderer.Go) string {
	if v.bits == 0 {
		return "uint"
	}

	return r.S(`uint$0`, v.bits)
}

// Pre to implement TypeHandler.
func (v *Uvarint) Pre(r *renderer.Go, src string) {
	key := gogh.Private(dotIsSep("len", src))
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Uint($src)`, uniq)
	v.lenkey = uniq
}

// Len to implement TypeHandler.
func (v *Uvarint) Len() int {
	return -1
}

// LenExpr to implement TypeHandler.
func (v *Uvarint) LenExpr(r *renderer.Go, src string) string {
	return v.lenkey
}

// Encoding to implement TypeHandler.
func (v *Uvarint) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	if v.bits == 64 {
		r.L(`$dst = $bin.AppendUvarint($dst, $src)`)
	} else {
		r.L(`$dst = $bin.AppendUvarint($dst, uint64($src))`)
	}
}

// Decoding to implement TypeHandler.
func (v *Uvarint) Decoding(r *renderer.Go, dst, src string) bool {
	off := r.Uniq("off")
	val := r.Uniq("val", dst)
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r.Let(off, off)

	r.L(`{`)

	if v.bits == 64 {
		r.L(`    var $off uint`)
		r.L(`    $dst, $off = $bin.Uvarint($src)`)
	} else {
		r.Let(val, val)
		r.L(`    $val, $off := $bin.Uvarint($src)`)
	}
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode: $malformedUvarint").Rend(r)
	r.L(`    }`)
	if v.bits != 64 {
		r.L(`$dst = $0($val)`, v.Name(r), v.bits)
	}
	r.L(`    $src = $src[$off:]`)
	r.L(`}`)

	return true
}
