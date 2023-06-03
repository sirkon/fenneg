package generator

import (
	"go/types"
	"strings"

	"github.com/facette/natsort"
	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/logger"
	"github.com/sirkon/fenneg/internal/renderer"
)

// NewOpLog construct new generator for interfaces instance.
func NewOpLog(
	logger logger.Type,
	src *types.Named,
	typ *types.Named,
	hnd *types.Named,
	disp string,
	lenPrefix bool,
	r *renderer.Go,
	thands map[*types.Var]handlers.Type,
	encoderResultType *types.Tuple,
	missingWriteBuffer bool,
) *OpLog {
	res := &OpLog{
		logger:             logger,
		src:                src,
		typ:                typ,
		hnd:                hnd,
		disp:               disp,
		lenPrefix:          lenPrefix,
		r:                  r,
		hands:              thands,
		encodeResultType:   encoderResultType,
		missingWriteBuffer: missingWriteBuffer,
	}
	res.rcvr = res.computeReceiverName(r)

	return res
}

// OpLog flat encoding/decoding code generation for operation log.
type OpLog struct {
	logger    logger.Type
	src       *types.Named
	typ       *types.Named
	hnd       *types.Named
	disp      string
	lenPrefix bool
	r         *renderer.Go

	encodeResultType   *types.Tuple
	missingWriteBuffer bool

	hands map[*types.Var]handlers.Type
	rcvr  string
}

type tuple struct {
	method string
	name   string
}

// computeReceiverName if to fulfil the need to make sure
// the receiver name will not clash with any of arguments.
//
// Preferred names are (in this exactly order):
//
//   1. Receiver name used in its existing methods.
//   2. First letter 𝒜 of the type in a low case.
//   3. Just x
//
// We check them first in each method.
func (g *OpLog) computeReceiverName(r *renderer.Go) (recvName string) {
	defer func() {
		r.Uniq(recvName)
	}()

	for i := 0; i < g.typ.NumMethods(); i++ {
		m := g.typ.Method(i)
		recv := m.Type().(*types.Signature).Recv()
		if recv != nil && recv.Name() != "" {
			return recv.Name()
		}
	}

	fl := []rune(strings.ToLower(g.typ.Obj().Name()))[0]
	fls := string(fl)
	if g.checkReceiverName(fl) {
		return fls
	}

	if g.checkReceiverName('x') {
		return "x"
	}

	var choices []string
	for i := 0; i < g.src.NumMethods(); i++ {
		ps := g.typ.Method(i).Type().(*types.Signature).Params()
		ir := r.Scope()
		for i := 0; i < ps.Len(); i++ {
			p := ps.At(i)
			ir.Uniq(p.Name())
		}
		choices = append(choices, ir.Uniq(fls))
	}

	natsort.Sort(choices)
	return choices[len(choices)-1]
}

func (g *OpLog) checkReceiverName(letter rune) bool {
	prefix := string(letter)

	for i := 0; i < g.src.NumMethods(); i++ {
		ps := g.typ.Method(i).Type().(*types.Signature).Params()
		for j := 0; j < ps.Len(); j++ {
			p := ps.At(j)
			if strings.HasPrefix(p.Name(), prefix) {
				return false
			}
		}
	}

	return true
}

func (g *OpLog) srcIface() *types.Interface {
	return g.src.Underlying().(*types.Interface)
}

func uptrd(t types.Type) types.Type {
	v, ok := t.(*types.Pointer)
	if ok {
		return v
	}

	return t
}

func isptr(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
}
