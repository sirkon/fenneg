package handlers

import (
	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// Duration handles time.Duration type.
type Duration struct{}

func (t *Duration) Name(r *renderer.Go) string {
	return "time.Duration"
}

func (t *Duration) Pre(r *renderer.Go, src string) {}

func (t *Duration) Len() int { return 8 }

func (t *Duration) LenExpr(r *renderer.Go, src string) string { return "8" }

func (t *Duration) Encoding(r *renderer.Go, dst, src string) {
	r.Imports().Binary().Ref("binary")
	r.L(`$dst = $binary.LittleEndian.AppendUint64($dst, uint64($src))`)
}

func (t *Duration) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")
	r.Imports().Binary().Ref("binary")
	r.Imports().Add("time").Ref("time")
	r.L(`if len($src) < 8 {`)
	er.Return().New("$decode: $recordTooSmall").LenReq(1).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`$dst = $time.Duration($binary.LittleEndian.Uint64($src[:8]))`)
	return false
}
