package fenneg

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/sirkon/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/sirkon/fenneg/internal/handlers"
	"github.com/sirkon/fenneg/internal/tdetect"
)

var builtinSupport = map[string]func() TypeHandler{}
var builtinHandlingSupport []func(p types.Type) TypeHandler

// NewTypesHandlers constructor. User can additional handlers via TypeHandlerByName
// and NewTypeHandler.
func NewTypesHandlers(handlers ...CustomTypeHandler) (*TypesHandlers, error) {
	builtins := maps.Clone(builtinSupport)
	builtinsHandling := slices.Clone(builtinHandlingSupport)
	res := &TypesHandlers{
		fset:     token.NewFileSet(),
		fixed:    builtins,
		handlers: builtinsHandling,
	}

	for i, handler := range handlers {
		if err := handler(handlerPlaceholder{}, res); err != nil {
			return nil, errors.Wrap(err, "process custom handler").Int("invalid-handler-index", i)
		}
	}

	return res, nil
}

// CustomTypeHandler custom option type. Cannot be defined outside of this package
// and only implemented by TypeHandlerByName and NewTypeHandler functions.
type CustomTypeHandler func(handlerPlaceholder, *TypesHandlers) error

type handlerPlaceholder struct {
	handler func(name string) TypeHandler
}

// TypesHandlers this returns type handlers.
type TypesHandlers struct {
	fset     *token.FileSet
	fixed    map[string]func() TypeHandler
	handlers []func(typ types.Type) TypeHandler
}

// TypeHandlerByName adds this custom handler for the given type.
func TypeHandlerByName(name string, handler func() TypeHandler) CustomTypeHandler {
	return func(_ handlerPlaceholder, dispatch *TypesHandlers) error {
		hv := handler()
		if hv == nil {
			return errors.New("zero handler definition is not allowed").Str("type-name", name)
		}

		if _, ok := builtinSupport[name]; ok {
			return errors.New("builtin types redefinition is not allowed").
				Str("type-name", name).
				Type("handler-type-attempted", handler())
		}

		if v, ok := dispatch.fixed[name]; ok {
			return errors.New("duplicate definition for the type").
				Str("type-name", name).
				Type("handler-type-previous", v).
				Type("handler-type-attempted", hv)
		}

		dispatch.fixed[name] = handler
		return nil
	}
}

// NewTypeHandler adds custom handler.
func NewTypeHandler(handler func(p types.Type) TypeHandler) CustomTypeHandler {
	return func(_ handlerPlaceholder, dispatch *TypesHandlers) error {
		dispatch.handlers = append(dispatch.handlers, handler)

		return nil
	}
}

// Handler returns a handler for the given type by its name.
func (h *TypesHandlers) Handler(arg *types.Var) TypeHandler {
	var name string

	typ := arg.Type()

	switch v := typ.(type) {
	case *types.Basic:
		name = v.Name()
	case *types.Slice:
		name = v.String()
	}

	if v, ok := h.fixed[name]; ok {
		return v()
	}

	for _, h := range h.handlers {
		hv := h(typ)
		if hv != nil {
			return hv
		}
	}

	// Type support was not found. May be it is a container of supported types or a pointer to it or a struct?
	switch t := arg.Type().(type) {
	case *types.Slice:
		hh := h.Handler(types.NewVar(arg.Pos(), arg.Pkg(), arg.Name(), t.Elem()))
		if hh == nil {
			break
		}

		return handlers.NewSlices(hh, t.Elem())
	case *types.Map:
		kh := h.Handler(types.NewVar(arg.Pos(), arg.Pkg(), arg.Name(), t.Key()))
		if kh == nil {
			break
		}
		vh := h.Handler(types.NewVar(arg.Pos(), arg.Pkg(), arg.Name(), t.Elem()))
		if vh == nil {
			break
		}

		return handlers.NewMaps(kh, t.Key(), vh, t.Elem())
	case *types.Alias:
		if t.Obj().Pkg().Path() != inttypesPkgPath {
			break
		}

		switch t.Obj().Name() {
		case "VI":
			return handlers.ArchVarint()
		case "VI16":
			return handlers.Varint16()
		case "VI32":
			return handlers.Varint32()
		case "VI64":
			return handlers.Varint64()
		case "VU":
			return handlers.ArchUvarint()
		case "VU16":
			return handlers.Uvarint16()
		case "VU32":
			return handlers.Uvarint32()
		case "VU64":
			return handlers.Uvarint64()
		default:
			break
		}

	case *types.Named:
		s, ok := t.Underlying().(*types.Struct)
		if !ok {
			break
		}

		sh := handlers.NewStruct(t)
		for i := 0; i < s.NumFields(); i++ {
			f := s.Field(i)
			fh := h.Handler(f)
			sh.AddField(f.Name(), f.Type(), fh)
		}

		return sh
	case *types.Struct:
		sh := handlers.NewStruct(t)
		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			fh := h.Handler(f)
			sh.AddField(f.Name(), f.Type(), fh)
		}

		return sh

	case *types.Pointer:
		// TODO add *T support for supported T.

	}

	return nil
}

func (h *TypesHandlers) fs() *token.FileSet {
	return h.fset
}

func init() {
	builtinSupport = map[string]func() TypeHandler{
		"bool": func() TypeHandler {
			return handlers.Bool{}
		},
		"int": func() TypeHandler {
			return handlers.ArchInt()
		},
		"int8": func() TypeHandler {
			return handlers.Int8()
		},
		"int16": func() TypeHandler {
			return handlers.Int16()
		},
		"int32": func() TypeHandler {
			return handlers.Int32()
		},
		"int64": func() TypeHandler {
			return handlers.Int64()
		},
		"uint": func() TypeHandler {
			return handlers.ArchUint()
		},
		"uint8": func() TypeHandler {
			return handlers.Uint8()
		},
		"uint16": func() TypeHandler {
			return handlers.Uint16()
		},
		"uint32": func() TypeHandler {
			return handlers.Uint32()
		},
		"uint64": func() TypeHandler {
			return handlers.Uint64()
		},
		"float32": func() TypeHandler {
			return handlers.Float32()
		},
		"float64": func() TypeHandler {
			return handlers.Float64()
		},
		"[]byte": func() TypeHandler {
			return &handlers.Bytes{}
		},
		"string": func() TypeHandler {
			return &handlers.StringHandler{}
		},
		"[]int16": func() TypeHandler {
			return handlers.NewSliceInt16()
		},
		"[]int32": func() TypeHandler {
			return handlers.NewSliceInt32()
		},
		"[]int64": func() TypeHandler {
			return handlers.NewSliceInt64()
		},
		"[]uint16": func() TypeHandler {
			return handlers.NewSliceUint16()
		},
		"[]uint32": func() TypeHandler {
			return handlers.NewSliceUint32()
		},
		"[]uint64": func() TypeHandler {
			return handlers.NewSliceUint64()
		},
	}

	builtinHandlingSupport = append(builtinHandlingSupport,

		// varints
		func(p types.Type) TypeHandler {
			v, ok := p.(*types.Alias)
			if !ok {
				return nil
			}

			if v.Obj().Pkg().Path() != "github.com/sirkon/intypes" {
				return nil
			}

			switch v.Obj().Name() {
			case "VI":
				return handlers.ArchVarint()
			case "VI16":
				return handlers.Varint16()
			case "VI32":
				return handlers.Varint32()
			case "VI64":
				return handlers.Varint64()
			case "VU":
				return handlers.ArchUvarint()
			case "VU16":
				return handlers.Uvarint16()
			case "VU32":
				return handlers.Uvarint32()
			case "VU64":
				return handlers.Uvarint64()
			default:
				panic(fmt.Sprintf("type intypes.%s is not supported", v.Obj().Name()))
			}
		},

		// byte arrays
		func(p types.Type) TypeHandler {
			ord := bytesArrayMatch(p)
			if ord < 0 {
				return nil
			}

			return handlers.BytesArray(ord)
		},

		// auto-support
		func(p types.Type) TypeHandler {
			if !tdetect.TypeImplementsEncoder(p) || !tdetect.TypeImplementsDecoder(p) {
				return nil
			}

			return handlers.NewAuto(p)
		},
	)
}

// bytesArrayMatch return an array of bytes dimension if it is an [N]byte.
// Return -1 otherwise.
func bytesArrayMatch(typ types.Type) int {
	v, ok := typ.(*types.Array)
	if !ok {
		return -1
	}

	if !tdetect.IsBasic(v.Elem(), types.Byte) {
		return -1
	}

	return int(v.Len())
}
