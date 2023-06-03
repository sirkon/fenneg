package er

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirkon/fenneg/internal/renderer"
)

type (
	RType struct {
		m *strings.Builder
	}
	RAttr struct {
		m *strings.Builder
	}
)

func Return() *RType {
	var b strings.Builder
	return &RType{
		m: &b,
	}
}

func (e *RType) New(msg string) *RAttr {
	e.m.WriteString(`return $errors.New(`)
	e.m.WriteString(strconv.Quote(msg))
	e.m.WriteString(")")

	return &RAttr{
		m: e.m,
	}
}

func (e *RType) Newf(msg string, a ...string) *RAttr {
	e.m.WriteString("return $errors.Newf(")
	e.m.WriteString(strconv.Quote(msg))
	for _, w := range a {
		e.m.WriteString(", ")
		e.m.WriteString(w)
	}
	e.m.WriteString(")")

	return &RAttr{
		m: e.m,
	}
}

func (e *RType) Wrap(err string, msg string) *RAttr {
	e.m.WriteString(`return $errors.Wrap(`)
	e.m.WriteString(err)
	e.m.WriteString(`, `)
	e.m.WriteString(strconv.Quote(msg))
	e.m.WriteString(")")

	return &RAttr{
		m: e.m,
	}
}

func (e *RAttr) o(method string, k, v any) *RAttr {
	e.m.WriteByte('.')
	e.m.WriteString(method)
	e.m.WriteByte('(')

	var tmp string
	switch v := k.(type) {
	case string:
		tmp = Q(v).String()
	case fmt.Stringer:
		tmp = v.String()
	default:
		panic(fmt.Sprintf("unsupported key type %T", k))
	}
	e.m.WriteString(tmp)

	e.m.WriteString(", ")
	switch v := v.(type) {
	case string:
		tmp = v
	case fmt.Stringer:
		tmp = v.String()
	default:
		tmp = fmt.Sprint(v)
	}
	e.m.WriteString(tmp)
	e.m.WriteByte(')')

	return e
}

func (e *RAttr) Pfx(v any) *RAttr {
	e.m.WriteString(`.Pfx(`)
	var tmp string
	switch vv := v.(type) {
	case string:
		tmp = Q(vv).String()
	case fmt.Stringer:
		tmp = vv.String()
	default:
		panic(fmt.Sprintf("unsupporter prefix literal type %T", v))
	}
	e.m.WriteString(tmp)
	e.m.WriteByte(')')

	return e
}

func (e *RAttr) Bool(k, v any) *RAttr {
	return e.o("Bool", k, v)
}

func (e *RAttr) Int(k, v any) *RAttr {
	return e.o("Int", k, v)
}

func (e *RAttr) Int8(k, v any) *RAttr {
	return e.o("Int8", k, v)
}

func (e *RAttr) Int16(k, v any) *RAttr {
	return e.o("Int16", k, v)
}

func (e *RAttr) Int32(k, v any) *RAttr {
	return e.o("Int32", k, v)
}

func (e *RAttr) Int64(k, v any) *RAttr {
	return e.o("Int64", k, v)
}

func (e *RAttr) Uint(k, v any) *RAttr {
	return e.o("Uint", k, v)
}

func (e *RAttr) Uint8(k, v any) *RAttr {
	return e.o("Uint8", k, v)
}

func (e *RAttr) Uint16(k, v any) *RAttr {
	return e.o("Uint16", k, v)
}

func (e *RAttr) Uint32(k, v any) *RAttr {
	return e.o("Uint32", k, v)
}

func (e *RAttr) Uint64(k, v any) *RAttr {
	return e.o("Uint64", k, v)
}

func (e *RAttr) Float32(k, v any) *RAttr {
	return e.o("Float32", k, v)
}

func (e *RAttr) Float64(k, v any) *RAttr {
	return e.o("Float64", k, v)
}

func (e *RAttr) Str(k, v any) *RAttr {
	return e.o("Str", k, v)
}

func (e *RAttr) Stg(k, v any) *RAttr {
	return e.o("Stg", k, v)
}

func (e *RAttr) Strs(k, v any) *RAttr {
	return e.o("Strs", k, v)
}

func (e *RAttr) Type(k, v any) *RAttr {
	return e.o("Type", k, v)
}

func (e *RAttr) Any(k, v any) *RAttr {
	return e.o("Any", k, v)
}

func (e *RAttr) LenReq(v any) *RAttr {
	return e.o("Uint64", "length-required", fmt.Sprintf("uint64(%v)", v))
}

func (e *RAttr) LenSrc() *RAttr {
	return e.o("Int", "length-actual", "len($src)")
}

func (e *RAttr) Rend(r *renderer.Go, a ...any) {
	r = feedRenderer(r)
	r.L(e.m.String(), a...)
}
