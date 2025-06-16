package handlers

import (
	"go/types"
	"strconv"

	"github.com/sirkon/errors"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// HandlerStruct support of struct type.
type HandlerStruct struct {
	typ    types.Type
	tuples []fieldTuple
	lenkey string
}

// NewStruct constructs HandlerStruct.
// It should be filled with fields data then.
func NewStruct(typ types.Type) *HandlerStruct {
	return &HandlerStruct{
		typ: typ,
	}
}

func (s *HandlerStruct) Name(r *renderer.Go) string {
	return r.Type(s.typ)
}

func (s *HandlerStruct) Pre(r *renderer.Go, src string) {
	if isFixed(s) {
		return
	}

	key := gogh.Private(dotIsSep("len", src))
	uniq := r.Uniq(key)
	s.lenkey = uniq

	r.L(`var $0 int`, s.lenkey)
	r.L(`{`)
	var lenExprs []string
	for _, tuple := range s.tuples {
		r := r.Scope()
		fieldSrc := r.S("$src.$0", tuple.name)
		r.Let("src", fieldSrc)
		tuple.handler.Pre(r, fieldSrc)
		r.L(`$0 += $1`, s.lenkey, tuple.handler.LenExpr(r, fieldSrc))
		lenExprs = append(lenExprs, tuple.handler.LenExpr(r, fieldSrc))
	}
	r.L(`}`)
	r.N()
}

func (s *HandlerStruct) Len() int {
	var res int
	for _, t := range s.tuples {
		l := t.handler.Len()
		if l <= 0 {
			return -1
		}

		res += l
	}

	return res
}

func (s *HandlerStruct) LenExpr(r *renderer.Go, src string) string {
	l := s.Len()
	if l > 0 {
		return strconv.Itoa(l)
	}

	return s.lenkey
}

func (s *HandlerStruct) Encoding(r *renderer.Go, dst, src string) {
	for _, tpl := range s.tuples {
		r := r.Scope()
		fieldSrc := src + "." + tpl.name
		r.Let("src", fieldSrc)
		tpl.handler.Encoding(r, dst, fieldSrc)
	}
}

func (s *HandlerStruct) Decoding(r *renderer.Go, dst, src string) bool {
	for _, tpl := range s.tuples {
		r := r.Scope()
		fieldDst := dst + "." + tpl.name
		r.Let("dst", fieldDst)
		tpl.handler.Decoding(r, fieldDst, src)
		if isFixed(tpl.handler) {
			r.L(`$src = $src[$0:]`, tpl.handler.Len())
		}
	}

	return isVariadic(s)
}

// AddField add field info.
func (s *HandlerStruct) AddField(name string, typ types.Type, handler Type) *HandlerStruct {
	for _, tuple := range s.tuples {
		if tuple.name == name {
			panic(errors.Newf("field %q has already been added", name))
		}
	}

	s.tuples = append(s.tuples, fieldTuple{
		name:    name,
		typ:     typ,
		handler: handler,
	})

	return s
}

type fieldTuple struct {
	name    string
	typ     types.Type
	handler Type
}
