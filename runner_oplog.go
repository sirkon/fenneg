package fenneg

import (
	"go/types"

	"github.com/sirkon/errors"
	"github.com/sirkon/gogh"
	"github.com/sirkon/message"
	"golang.org/x/tools/go/packages"

	"github.com/sirkon/fenneg/internal/app"
	"github.com/sirkon/fenneg/internal/generator"
	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/logger"
	"github.com/sirkon/fenneg/internal/tdetect"
)

// RunnerOpLog oplog codegen runner.
type RunnerOpLog struct {
	r        *Runner
	src      sourcePoint
	typ      sourcePoint
	hnd      sourcePoint
	dispatch string
	genlen   bool
}

// Run codegen.
func (r *RunnerOpLog) Run() error {

	log := r.r.logger
	gargs, err := r.computeOplogArgs()
	if err != nil {
		return errors.Wrap(err, "detect types")
	}

	aret, wret, ok := checkOpLogTypesPrerequisites(log, gargs.typ)
	if !ok {
		return errors.Newf(
			"encoder type %s:%s prerequisites checks failed",
			r.typ.Path,
			r.typ.ID,
		)
	}
	noWriteBuffer := wret == nil
	if noWriteBuffer {
		wret = aret
	}

	p, gr, err := getRenderer(app.Name, r.r.fset, gargs.typ)
	if err != nil {
		return errors.Wrap(err, "set up renderer")
	}

	h, err := getOpLogArgsHandlers(log, gargs.src.Underlying().(*types.Interface), r.r.handlers)
	if err != nil {
		return errors.Wrap(err, "set up arguments types handlers")
	}

	g := generator.NewOpLog(
		log,
		gargs.src,
		gargs.typ,
		gargs.hnd,
		r.dispatch,
		r.genlen,
		gr,
		h,
		wret,
		noWriteBuffer,
	)
	// All code generation in this call.
	message.AddFatalHook(func() {
		panic("fatal stacktrace")
	})
	g.Generate()

	if err := p.Render(); err != nil {
		return errors.Wrap(err, "render generated source code")
	}

	return nil
}

// DispatchFunc sets dispatch func name.
func (r *RunnerOpLog) DispatchFunc(name string) *RunnerOpLog {
	r.dispatch = name
	return r
}

// DispatchHandler sets handler type for the dispatching type.
func (r *RunnerOpLog) DispatchHandler(pkg, typ string) *RunnerOpLog {
	r.hnd = sourcePoint{
		Path: pkg,
		ID:   typ,
	}
	return r
}

// LengthPrefix sets generation of a length uleb128 prefix before the rest.
func (r *RunnerOpLog) LengthPrefix(on bool) *RunnerOpLog {
	r.genlen = on
	return r
}

// OpLog run operation log generator.
func (r *Runner) OpLog() *RunnerOpLogSource {
	return &RunnerOpLogSource{
		r: &RunnerOpLog{
			r: r,
		},
	}
}

// RunnerOpLogSource oplog runner component to set up
// oplog source interface.
type RunnerOpLogSource struct {
	r *RunnerOpLog
}

// RunnerOpLogType oplog runner component to set up
// oplog encoding type.
type RunnerOpLogType struct {
	r *RunnerOpLog
}

// Source set up oplog source interface.
func (r *RunnerOpLogSource) Source(pkg, typ string) *RunnerOpLogType {
	r.r.src = sourcePoint{
		Path: pkg,
		ID:   typ,
	}
	r.r.hnd = r.r.src
	return &RunnerOpLogType{
		r: r.r,
	}
}

// Type set up oplog encoding type.
func (r *RunnerOpLogType) Type(pkg, typ string) *RunnerOpLog {
	r.r.typ = sourcePoint{
		Path: pkg,
		ID:   typ,
	}
	r.r.dispatch = gogh.Public(typ, "Dispatch")

	return r.r
}

func (r *RunnerOpLog) computeOplogArgs() (
	res struct {
		src *types.Named
		typ *types.Named
		hnd *types.Named

		spkg  *packages.Package
		types map[*types.Var]*types.Named
	},
	err error,
) {
	l := newSouceLoader(r.r.fset)

	spkg, src, err := getType(l, r.src, typeByNameSpecific[*types.Interface], "source")
	if err != nil {
		return res, errors.Wrap(err, "get source interface")
	}
	res.src = src
	res.spkg = spkg
	res.types = l.mapExactTypes(src)

	_, typ, err := getType(l, r.typ, typeByName, "type")
	if err != nil {
		return res, errors.Wrap(err, "get type")
	}
	res.typ = typ

	res.hnd = src
	if !r.hnd.IsValid() {
		return res, nil
	}

	_, hnd, err := getType(l, r.hnd, typeByName, "handler")
	if err != nil {
		return res, errors.Wrap(err, "get handler type")
	}
	res.hnd = hnd

	return res, nil
}

// getOpLogArgsHandlers set up handlers for every argument of the
// source interface.
func getOpLogArgsHandlers(
	log logger.Type,
	iface *types.Interface,
	dispatch *TypesHandlers,
) (map[*types.Var]handlers.Type, error) {
	var errcount int
	res := map[*types.Var]handlers.Type{}

outer:
	for i := 0; i < iface.NumMethods(); i++ {
		ps := iface.Method(i).Type().(*types.Signature).Params()
		for j := 0; j < ps.Len(); j++ {
			p := ps.At(j)
			hnd := dispatch.Handler(p)
			if hnd == nil {
				if errcount == 10 {
					break outer
				}

				log.Pos(p.Pos(), errors.New("argument type is not supported").Stg("type-name", p.Type()))
				errcount++
				continue
			}

			res[p] = hnd
		}
	}

	if errcount > 0 {
		return nil, errors.New("unsupported arguments types detected, cannot continue")
	}

	return res, nil
}

// checkOpLogTypesPrerequisites do
//
//  1. check if allocateBuffer(int) []byte is here
//  2. check if writeBuffer([]byte) <tuple> is here, this is optional.
//
// Result values:
//   - aret is a result tuple of the allocateBuffer.
//   - wret is a result tuple of the writeBuffer if it does exist.
//   - err != nil if checks haven't passed.
func checkOpLogTypesPrerequisites(
	l LoggerType,
	typ *types.Named,
) (aret *types.Tuple, wret *types.Tuple, ok bool) {
	var allocSuccess bool
	var hasAlloc bool
	writeBufferSuccess := true
	for i := 0; i < typ.NumMethods(); i++ {
		m := typ.Method(i)
		switch m.Name() {
		case "allocateBuffer":
			hasAlloc = true
			aret, allocSuccess = checkAllocateBuffer(l, m)
		case "writeBuffer":
			wret, writeBufferSuccess = checkWriteBuffer(l, m)
		default:
			continue
		}
	}

	if !hasAlloc {
		l.Pos(typ.Obj().Pos(), errors.New("missing allocateBuffer method on this type"))
		return nil, nil, false
	}

	return aret, wret, allocSuccess && writeBufferSuccess
}

func checkAllocateBuffer(l LoggerType, f *types.Func) (*types.Tuple, bool) {
	s := f.Type().(*types.Signature)
	if s.Params().Len() != 1 || !tdetect.IsBasic(s.Params().At(0).Type(), types.Int) {
		l.Pos(
			f.Pos(),
			errors.New("allocateBuffer: invalid number of parameters and/or their types").
				Pfx("allocate-buffer-params").
				Str("expected", "int").
				Stg("actual", s.Params()),
		)
		return nil, false
	}

	if s.Results().Len() != 1 || !tdetect.IsByteSlice(s.Results().At(0).Type()) {
		l.Pos(
			f.Pos(),
			errors.New("allocateBuffer: invalid number of result values and/or their types").
				Pfx("allocate-buffer-params").
				Str("expected", "[]byte").
				Stg("actual", s.Results()),
		)
		return nil, false
	}

	return s.Results(), true
}

func checkWriteBuffer(l LoggerType, m *types.Func) (*types.Tuple, bool) {
	s := m.Type().(*types.Signature)
	if s.Params().Len() != 1 || !tdetect.IsByteSlice(s.Params().At(0).Type()) {
		l.Pos(
			m.Pos(),
			errors.New("writeBuffer: invalid number of parameters and/or their types").
				Pfx("allocate-buffer-params").
				Str("expected", "[]byte").
				Stg("actual", s.Params()),
		)
	}

	return s.Results(), true
}
