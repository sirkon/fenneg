package handlers

import (
	"go/types"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// MapsUniform map[K]V types support for supported K and V types.
type MapsUniform struct {
	khandler Type
	vhandler Type
	ktype    types.Type
	vtype    types.Type

	lenkey string
}

func NewMaps(khandler Type, ktype types.Type, vhandler Type, vtype types.Type) *MapsUniform {
	return &MapsUniform{
		khandler: khandler,
		vhandler: vhandler,
		ktype:    ktype,
		vtype:    vtype,
	}
}

func (m *MapsUniform) Name(r *renderer.Go) string {
	return "map[" + m.khandler.Name(r) + "]" + m.vhandler.Name(r)
}

func (m *MapsUniform) Pre(r *renderer.Go, src string) {
	key := gogh.Private(dotIsSep("len", src))
	uniq := r.Uniq(key)
	m.lenkey = uniq

	if isFixed(m.khandler) && isFixed(m.vhandler) {
		r.L(`$0 := len($src)*($1 + $2)`, m.lenkey, m.khandler.Len(), m.vhandler.Len())
		return
	}

	item := r.Uniq("key", src)
	val := r.Uniq("value", src)
	r.Imports().Varsize().Ref("vsize")

	if isFixed(m.khandler) {
		r.L(`$0 := $vsize.MapLen($src) + len($src)*$1`, m.lenkey, m.khandler.Len())
		r.L(`for _, $0 := range $src {`, val)
		r = r.Scope()
		r.Let("src", val)
		m.vhandler.Pre(r, val)
		r.L(`    $0 += $1`, uniq, m.vhandler.LenExpr(r, val))
		r.L(`}`)
		return
	}

	if isFixed(m.vhandler) {
		r.L(`$0 := $vsize.MapLen($src) + len($src)*$1`, m.lenkey, m.vhandler.Len())
		r.L(`for $0 := range $src {`, item)
		r = r.Scope()
		r.Let("src", item)
		m.khandler.Pre(r, item)
		r.L(`    $0 += $1`, uniq, m.khandler.LenExpr(r, item))
		r.L(`}`)
		return
	}

	r.L(`$0 := $vsize.MapLen($src)`)
	r.L(`for _, $0 := range $src{`)
	r = r.Scope()
	r.Let("src", item)
	m.khandler.Pre(r, item)
	r = r.Scope()
	r.Let("src", val)
	m.vhandler.Pre(r, val)
	r.L(`    $0 += $1 + $2`, uniq, m.khandler.LenExpr(r, item), m.vhandler.LenExpr(r, val))
	r.L(`}`)
}

func (m *MapsUniform) Len() int {
	return -1
}

func (m *MapsUniform) LenExpr(r *renderer.Go, src string) string {
	return m.lenkey
}

func (m *MapsUniform) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r = r.Scope()

	r.L(`$dst = $bin.AppendUvarint($dst, uint64(len($src)))`)
	k := r.Uniq("k", src)
	v := r.Uniq("v", src)
	r.L(`for $0, $1 := range $src {`, k, v)
	r = r.Scope()
	r.Let("src", k)
	m.khandler.Encoding(r, dst, k)
	r = r.Scope()
	r.Let("src", v)
	m.vhandler.Encoding(r, dst, v)
	r.L(`}`)
}

func (m *MapsUniform) Decoding(r *renderer.Go, dst, src string) bool {
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
	r.L(`    $dst = make(map[$0]$1, $siz)`, r.Type(m.ktype), r.Type(m.vtype))

	it := r.Uniq("i")
	r.L(`    for $0 := 0; $0 < int($siz); $0++ {`, it)
	kkk := r.Uniq("k")
	vvv := r.Uniq("v")
	r.L(`    var $0 $1`, kkk, m.khandler.Name(r))
	r.L(`    var $0 $1`, vvv, m.vhandler.Name(r))

	// Decode key.
	{
		r := r.Scope()
		r.Let("dst", kkk)
		if !m.khandler.Decoding(r, kkk, src) {
			r.L(`$src = $src[$0:]`, m.khandler.Len())
		}
	}

	// Decode value.
	{
		r := r.Scope()
		r.Let("dst", vvv)
		if !m.vhandler.Decoding(r, vvv, src) {
			r.L(`$src = $src[$0:]`, m.vhandler.Len())
		}
	}

	r.L(`        $dst[$0] = $1`, kkk, vvv)

	r.L(`    }`)
	r.L(`}`)

	return true
}
