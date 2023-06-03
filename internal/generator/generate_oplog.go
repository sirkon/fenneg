package generator

import (
	"go/types"
	"strings"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

// Generate perform code generation.
func (g *OpLog) Generate() {
	r := g.r

	r.L(`const (`)
	for i := 0; i < g.srcIface().NumMethods(); i++ {
		r.L(`    $0 = $1`, g.methodCode(g.srcIface().Method(i)), i+1)
	}
	r.L(`)`)

	r = r.Scope()
	r.Let("rcvr", r.S("$0 *$1", g.rcvr, r.Type(g.typ)))
	r.Let("self", g.rcvr)

	for i := 0; i < g.srcIface().NumMethods(); i++ {
		m := g.srcIface().Method(i)
		g.genMethod(r.Scope(), m)
	}

	r.N()
	var hnd string
	if g.hnd != g.src {
		hnd = "*" + r.Type(g.hnd)
	} else {
		hnd = r.Type(g.src)
	}

	r.Imports().Errors().Ref("errors")
	r.Imports().Binary().Ref("bin")
	r.Uniq("disp")
	r.Uniq("rec")
	r.Let("src", "rec")
	r.Imports().Errors().Ref("errors")

	r.L(`// $0 dispatches encoded data made with $1`, g.disp, r.Type(g.typ))
	r.L(`func $0(disp $1, rec []byte) error {`, g.disp, hnd)
	r.L(`    if len(rec) < 4 {`)
	er.Return().New("decode branch code: $recordTooSmall").LenReq("4").LenSrc().Rend(r)
	r.L(`    }`)
	r.N()
	r.L(`    branch := $bin.LittleEndian.Uint32(rec[:4])`)
	r.L(`    rec = rec[4:]`)
	r.N()
	r.L(`    switch branch {`)
	for i := 0; i < g.srcIface().NumMethods(); i++ {
		m := g.srcIface().Method(i)
		g.genBranchHandling(r, m)
		r.N()
	}
	r.L(`    default:`)
	er.Return().Newf("invalid branch code %d", "branch").Uint32("invalid-branch-code", "branch").Rend(r)
	r.L(`    }`)
	r.N()
	r.L(`    return nil`)
	r.L(`}`)
}

func (g *OpLog) genMethod(r *renderer.Go, m *types.Func) {
	lenParts := []string{}

	lenParts = append(lenParts, "4")
	ps := m.Type().(*types.Signature).Params()
	r = r.Scope()

	r.L(`// $0 encodes arguments tuple of this method.`, m.Name())
	r.M("$rcvr")(m.Name())(g.params(r, m)).Returns(g.encodeResultType).Body(func(r *renderer.Go) {
		// Generate pre parts.
		for i := 0; i < ps.Len(); i++ {
			p := ps.At(i)
			h := g.hands[p]
			lr := r.Scope()
			lr.Let("src", p.Name())
			h.Pre(lr, p.Name())
			lenParts = append(lenParts, h.LenExpr(r, m.Name()))
		}

		// Make sure buffer name won't clash with anything
		r.Let("dst", r.Uniq("buf"))
		r.Imports().Binary().Ref("bin")

		// Standard part that relies on user-defined allocateBuffer method.
		if g.lenPrefix {
			r.InnerScope(func(r *gogh.GoRenderer[*renderer.Imports]) {
				r.Imports().Varsize().Ref("vars")

				r.Let("len", r.Uniq("bufSize"))
				r.L(`var $dst []byte`)
				r.L(`{`)
				r.L(`    $len := $0`, strings.Join(lenParts, "+"))
				r.L(`    $dst = $self.allocateBuffer($vars.Uint(uint64($len)) + $len)`)
				r.N()
				r.L(`    // Encode record length.`)
				r.L(`    $dst = $bin.AppendUvarint($dst, uint64($len))`)
				r.L(`}`)
				r.N()
			})
		} else {
			r.L(`    $dst := $self.allocateBuffer($0)`, strings.Join(lenParts, " + "))
		}
		r.N()

		// Standard buffer allocation part.
		r.L(`    // Encode branch (method) code.`)
		r.L(`    $dst = $bin.LittleEndian.AppendUint32($dst, uint32($0))`, g.methodCode(m))

		// Arguments handling.
		for i := 0; i < ps.Len(); i++ {
			p := ps.At(i)
			h := g.hands[p]

			lr := r.Scope()
			lr.Let("src", p.Name())
			lr.Let("srcType", r.Type(p.Type()))
			lr.N()
			lr.L(`// Encode $src($srcType).`)
			h.Encoding(lr, "buf", p.Name())
			lr.N()
		}

		// Standard part that relies on user-defined writeBuffer method.
		r.Imports().Errors().Ref("errors")
		if g.missingWriteBuffer {
			r.L(`return $dst`)
			return
		}

		r.L(`return $self.writeBuffer($dst)`)
	})
	r.N()
}

func (g *OpLog) genBranchHandling(r *renderer.Go, m *types.Func) {
	r.L(`case $0:`, g.methodCode(m))
	r = r.Scope()
	r.Let("branch", m.Name())

	r = r.Scope()
	ps := m.Type().(*types.Signature).Params()
	var params gogh.Commas
	for i := 0; i < ps.Len(); i++ {
		p := ps.At(i)
		h := g.hands[p]

		lr := r.Scope()
		param := r.Uniq(p.Name())
		if isptr(p.Type()) {
			params.Add("&" + param)
		} else {
			params.Add(param)
		}
		lr.Let("dst", param)
		lr.Let("dstType", r.Type(uptrd(p.Type())))

		lr.L(`// Decode $0($dstType).`, p.Name())
		lr.L(`var $dst $dstType`)
		if !h.Decoding(lr, param, "rec") {
			lr.L(`$src = $src[$0:]`, h.LenExpr(r, "rec"))
		}
		lr.N()
	}

	r.N()
	r.L(`    if len($src) > 0 {`)
	er.Return().
		New("decode $branch: the record was not emptied after the last argument decoded").
		Int("record-bytes-left", r.S(`len($src)`)).
		Rend(r)
	r.L(`    }`)

	r.N()
	r.L(`    if err := disp.$0($1); err != nil {`, m.Name(), params.String())
	er.Return().Wrap("err", "call $0").Rend(r, m.Name())
	r.L(`    }`)
	r.N()
	r.L(`    return nil`)
}

func (g *OpLog) methodCode(m *types.Func) string {
	return gogh.Private(g.src.Obj().Name(), "code", m.Name())
}

func (g *OpLog) params(r *renderer.Go, m *types.Func) gogh.Params {
	res := gogh.Params{}
	ps := m.Type().(*types.Signature).Params()
	for i := 0; i < ps.Len(); i++ {
		p := ps.At(i)
		res.Add(p.Name(), r.Type(p.Type()))
		r.Uniq(p.Name())
	}

	return res
}
