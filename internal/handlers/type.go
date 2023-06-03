package handlers

import "github.com/sirkon/fenneg/internal/renderer"

// Type a doubler of fenneg.TypeHandler to avoid cyclic imports.
type Type interface {
	Name(r *renderer.Go) string
	Pre(r *renderer.Go, src string)
	Len() int
	LenExpr(r *renderer.Go, src string) string
	Encoding(r *renderer.Go, dst, src string)
	Decoding(r *renderer.Go, dst, src string) bool
}
