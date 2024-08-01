package generator

import (
	"go/types"
	"strings"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
	"github.com/sirkon/message"
)

// NewStruct constructs Struct instance.
func NewStruct(
	r *renderer.Go,
	src *types.Named,
	hands map[*types.Var]handlers.Type,
	pointer bool,
	encoderName, decoderName, sizeName string,
) *Struct {
	return &Struct{
		r:           r,
		src:         src,
		hands:       hands,
		pointer:     pointer,
		encoderName: encoderName,
		decoderName: decoderName,
		sizeName:    sizeName,
	}
}

// Struct flat encoding/decoding code generation for structures.
type Struct struct {
	r *renderer.Go

	argName string
	src     *types.Named
	hands   map[*types.Var]handlers.Type

	pointer                  bool
	encoderName, decoderName string
	sizeName                 string
}

// Generate struct flat encoding/decoding.
func (g *Struct) Generate() {
	if g.argName == "" {
		name := []rune(g.src.Obj().Name())
		g.argName = strings.ToLower(string(name[0]))
	}

	g.generateLen(g.r)
	g.generateEncoding(g.r)
	g.generateDecoding(g.r)
}

func (g *Struct) generateLen(r *renderer.Go) {
	r = r.Scope()
	var fnR *gogh.GoFuncRenderer[*renderer.Imports]
	if g.pointer {
		r.L(`// $0 calculate the size of the $1 fields for encoding in bytes.`, g.sizeName, r.Type(g.src))
		fnR = r.M(g.argName, "*"+r.Type(g.src))(g.sizeName)()
	} else {
		fnR = r.F(
			gogh.Public(g.src.Obj().Name(), g.sizeName),
		)(
			g.argName,
			"*"+r.Type(g.src),
		)
	}

	fnR.Returns("int").Body(func(r *renderer.Go) {
		r.L(`if $0 == nil {`, g.argName)
		r.L(`	return 0`)
		r.L(`}`)
		r.N()
		s := g.src.Underlying().(*types.Struct)

		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			r.L(`// len $0($1).`, f.Name(), r.T().Type(f.Type()))
		}
		r.N()

		lens := make([]string, 0, s.NumFields())
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			h := g.hands[f]

			lr := r.Scope()
			lr.Let("src", g.argName+"."+f.Name())
			h.Pre(lr, f.Name())
			lenExpr := h.LenExpr(lr, "src")
			if lenExpr != "" {
				lens = append(lens, lenExpr)
			} else {
				message.Warningf("field %s has no len expression", f.Name())
			}
		}
		r.N()
		r.L(`return $0`, strings.Join(lens, " + "))
	})
}

func (g *Struct) generateEncoding(r *renderer.Go) {
	r = r.Scope()
	var fnR *gogh.GoFuncRenderer[*renderer.Imports]
	if g.pointer {
		r.L(`// $0 encodes content of $1.`, g.encoderName, r.Type(g.src))
		fnR = r.M(g.argName, "*"+r.Type(g.src))(g.encoderName)("dst", "[]byte")
	} else {
		fnR = r.F(
			gogh.Public(g.src.Obj().Name(), g.encoderName),
		)(
			"dst",
			"[]byte",
			g.argName,
			"*"+r.Type(g.src),
		)
	}

	fnR.Returns("[]byte").Body(func(r *renderer.Go) {
		r.L(`if $0 == nil {`, g.argName)
		r.L(`	return dst`)
		r.L(`}`)
		r.L(`if dst == nil {`, g.argName)
		if g.pointer {
			r.L(`	dst = make([]byte, 0, $0.$1())`, g.argName, g.sizeName)
		} else {
			r.L(`	dst = make([]byte, 0, $0($1))`, gogh.Public(g.src.Obj().Name(), g.sizeName), g.argName)
		}
		r.L(`}`)
		r.N()
		s := g.src.Underlying().(*types.Struct)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			r.L(`// Encode $0($1).`, f.Name(), r.T().Type(f.Type()))
			h := g.hands[f]

			rr := r.Scope()
			rr.Let("src", g.argName+"."+f.Name())
			rr.Let("dst", "dst")
			rr.Let("dstType", r.T().Type(f.Type()))
			h.Encoding(rr, "dst", g.argName+"."+f.Name())
			r.N()
		}
		r.L(`return dst`)
	})
}

func (g *Struct) generateDecoding(r *renderer.Go) {
	r = r.Scope()
	var fnR *gogh.GoFuncRenderer[*renderer.Imports]
	if g.pointer {
		r.L(`// $0 decodes content of $1.`, g.decoderName, r.Type(g.src))
		fnR = r.M(g.argName, "*"+r.Type(g.src))(g.decoderName)("src", "[]byte")
	} else {
		r.L(`// $0 decodes content of $1.`, gogh.Public(g.src.Obj().Name(), "Encode"), r.Type(g.src))
		fnR = r.F(gogh.Public(g.src.Obj().Name(), g.decoderName))(
			g.argName,
			"*"+r.Type(g.src),
			"src",
			"[]byte",
		)
	}

	fnR.Returns("err", "error").Body(func(r *renderer.Go) {
		s := g.src.Underlying().(*types.Struct)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			h := g.hands[f]

			r.L(`// Decode $0($1).`, f.Name(), r.T().Type(f.Type()))
			rr := r.Scope()
			rr.Let("dst", g.argName+"."+f.Name())
			rr.Let("dstType", r.T().Type(f.Type()))
			rr.Let("src", "src")
			if !h.Decoding(rr, g.argName+"."+f.Name(), "src") {
				r.L(`src = src[$0:]`, h.LenExpr(r, "src"))
			}
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
