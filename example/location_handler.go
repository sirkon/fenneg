package main

import (
	"strings"

	"github.com/sirkon/fenneg"
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/gogh"
)

type LocationHandler struct {
	lenkey string
}

func (l *LocationHandler) Name(r *fenneg.Go) string {
	return "Location"
}

func (l *LocationHandler) Pre(r *fenneg.Go, src string) {
	key := gogh.Private(dotIsSep("len", src))
	uniq := r.Uniq(key)
	r.Imports().Varsize().Ref("vsize")
	r.L(`$0 := $vsize.Uint(uint(len($src.File))) + $vsize.Uint(uint($src.Line))`, uniq)
	l.lenkey = uniq
}

func (l *LocationHandler) Len() int {
	return -1
}

func (l *LocationHandler) LenExpr(r *fenneg.Go, src string) string {
	return l.lenkey
}

func (l *LocationHandler) Encoding(r *fenneg.Go, dst, src string) {
	r.L(`if !$src.Valid {`)
	r.L(`    $dst = append($dst, 0);`)
	r.L(`} else {`)
	r.L(`    $dst = $binary.AppendUvarint($dst, uint64(len($src.File)))`)
	r.L(`    $dst = append($dst, $src.File...)`)
	r.L(`    $dst = $binary.AppendUvarint($dst, uint64($src.Line))`)
	r.L(`}`)
}

func (l *LocationHandler) Decoding(r *fenneg.Go, dst, src string) bool {
	r.Imports().Binary().Ref("binary")
	r.L(`if len($src) < 1 {`)
	er.Return().New("$decode length: $recordTooSmall").Rend(r)
	r.L(`}`)
	r.L(`if $src[0] == 0 {`)
	r.L(`    $dst.Valid = false`)
	r.L(`} else {`)
	r = r.Scope()

	off := r.Uniq("off")
	siz := r.Uniq("size")
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")
	r = r.Scope()
	r.Let("siz", siz)
	r.Let("off", off)

	r.L(`    $siz, $off := $bin.Uvarint($src)`)
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode File: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode File: $malformedUvarint").Rend(r)
	r.L(`    }`)
	r.L(`    $src = $src[$off:]`)
	r.L(`    if int($siz) > len($src) {`)
	er.Return().New("$decode content: $recordTooSmall").LenReq("int($siz)").LenSrc().Rend(r)
	r.L(`    }`)
	r.L(`    $dst.File = string($src[:$siz])`)
	r.L(`    $src = $src[$siz:]`)
	r.L(`    $siz, $off = $bin.Uvarint($src)`)
	r.L(`    if $off <= 0 {`)
	r.L(`        if $off == 0 {`)
	er.Return().New("$decode Line: $recordTooSmall").Rend(r)
	r.L(`        }`)
	er.Return().New("$decode Line: $malformedUvarint").Rend(r)
	r.L(`    }`)
	r.L(`    $dst.Line = int($siz)`)
	r.L(`    $src = $src[$off:]`)
	r.L(`}`)

	return true
}

func dotIsSep(head string, a ...string) string {
	var buf strings.Builder
	buf.WriteString(strings.ReplaceAll(head, ".", "_"))
	for _, p := range a {
		buf.WriteByte('_')
		buf.WriteString(strings.ReplaceAll(p, ".", "_"))
	}

	return buf.String()
}
