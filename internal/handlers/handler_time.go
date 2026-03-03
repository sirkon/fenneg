package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// Time handles time.Time type.
type Time struct{}

func (t *Time) Name(r *renderer.Go) string {
	return "time.Time"
}

func (t *Time) Pre(r *renderer.Go, src string) {}

func (t *Time) Len() int { return 8 }

func (t *Time) LenExpr(r *renderer.Go, src string) string { return "8" }

func (t *Time) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("binary")
	r.L(`$dst = $binary.LittleEndian.AppendUint64($dst, uint64($src.UnixNano()))`)
}

func (t *Time) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")
	r.Imports().Binary().Ref("binary")
	r.Imports().Add("time").Ref("time")
	r.L(`if len($src) < 8 {`)
	er.Return().New("$decode: $recordTooSmall").LenReq(1).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`$dst = $time.Unix(0, int64($binary.LittleEndian.Uint64($src[:8])))`)
	return false
}
