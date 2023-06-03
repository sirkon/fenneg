package generator

import (
	"go/types"
	"strings"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// NewStruct constructs Struct instance.
func NewStruct(
	r *renderer.Go,
	src *types.Named,
	hands map[*types.Var]handlers.Type,
) *Struct {
	return &Struct{
		r:     r,
		src:   src,
		hands: hands,
	}
}

// Struct flat encoding/decoding code generation for structures.
type Struct struct {
	r *renderer.Go

	argName string
	src     *types.Named
	hands   map[*types.Var]handlers.Type
}

// Generate struct flat encoding/decoding.
func (g *Struct) Generate() {
	if g.argName == "" {
		name := []rune(g.src.Obj().Name())
		g.argName = strings.ToLower(string(name[0]))
	}

	g.generateEncoding(g.r)
	g.generateDecoding(g.r)
}

func (g *Struct) generateEncoding(r *renderer.Go) {
	r = r.Scope()
	r.L(`// $0 encodes content of $1.`, gogh.Public(g.src.Obj().Name(), "Encode"), r.Type(g.src))
	r.F(
		gogh.Public(g.src.Obj().Name(), "Encode"),
	)(
		"dst",
		"[]byte",
		g.argName,
		"*"+r.Type(g.src),
	).Returns("[]byte").Body(func(r *renderer.Go) {
		s := g.src.Underlying().(*types.Struct)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			r.L(`// Encode $0($1).`, f.Name(), varTypeName(f))
			h := g.hands[f]

			rr := r.Scope()
			rr.Let("src", g.argName+"."+f.Name())
			rr.Let("dst", "dst")
			rr.Let("dstType", varTypeName(f))
			h.Encoding(rr, "dst", g.argName+"."+f.Name())
			r.N()
		}
		r.L(`return dst`)
	})
}

func (g *Struct) generateDecoding(r *renderer.Go) {
	r = r.Scope()
	r.L(`// $0 decodes content of $1.`, gogh.Public(g.src.Obj().Name(), "Encode"), r.Type(g.src))
	r.F(gogh.Public(g.src.Obj().Name(), "Decode"))(
		g.argName,
		"*"+r.Type(g.src),
		"src",
		"[]byte",
	).Returns("error").Body(func(r *renderer.Go) {
		s := g.src.Underlying().(*types.Struct)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			h := g.hands[f]

			r.L(`// Decode $0($1).`, f.Name(), varTypeName(f))
			rr := r.Scope()
			rr.Let("dst", g.argName+"."+f.Name())
			rr.Let("dstType", varTypeName(f))
			rr.Let("src", "src")
			h.Decoding(rr, g.argName+"."+f.Name(), "src")
			r.N()
		}

		r.L(`    if len(src) > 0 {`)
		er.Return().
			New("the record was not emptied after the last argument decoded").
			Int("record-bytes-left", r.S(`len(src)`)).
			Rend(r)
		r.L(`    }`)
		r.N()
		r.L(`return nil`)
	})
}

func varTypeName(v *types.Var) string {
	n := v.Type().String()
	lastIndex := strings.LastIndexByte(n, '/')
	if lastIndex < 0 {
		return n
	}

	return n[lastIndex+1:]
}
