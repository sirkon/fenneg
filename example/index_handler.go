package main

import (
	"go/types"
	"strconv"

	"github.com/sirkon/fenneg"
)

// NewHandlerIndex index type handler constructor.
func NewHandlerIndex(indexType *types.Named) *HandlerIndex {
	return &HandlerIndex{
		indexType: indexType,
	}
}

// HandlerIndex to handle example.Index type.
type HandlerIndex struct {
	indexType *types.Named
}

// Name to implement fenneg.TypeHandler.
func (h HandlerIndex) Name(r *fenneg.Go) string {
	return r.Type(h.indexType)
}

// Pre to implement fenneg.TypeHandler.
func (h HandlerIndex) Pre(r *fenneg.Go, src string) {}

// Len to implement fenneg.TypeHandler.
func (h HandlerIndex) Len() int {
	return 16
}

// LenExpr to implement fenneg.TypeHandler.
func (h HandlerIndex) LenExpr(r *fenneg.Go, src string) string {
	return strconv.Itoa(h.Len())
}

// Encoding to implement fenneg.TypeHandler.
func (h HandlerIndex) Encoding(r *fenneg.Go, dst, src string) {
	r.Imports().Binary().Ref("bin")
	r.L(`$dst = $bin.LittleEndian.AppendUint64($dst, $src.Term)`)
	r.L(`$dst = $bin.LittleEndian.AppendUint64($dst, $src.Index)`)
}

// Decoding to implement fenneg.TypeHandler.
func (h HandlerIndex) Decoding(r *fenneg.Go, dst, src string) bool {
	r.Imports().Binary().Ref("bin")
	r.Imports().Errors().Ref("errors")

	r.L(`if len($src) < 16 {`)
	fenneg.ReturnError().New("$decode: $recordTooSmall").LenReq(16).LenSrc().Rend(r)
	r.L(`}`)
	r.L(`$dst.Term = $bin.LittleEndian.Uint64($src)`)
	r.L(`$dst.Index = $bin.LittleEndian.Uint64($src[8:])`)

	return false
}
