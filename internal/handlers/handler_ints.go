package handlers

import (
	"strconv"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// ArchInt handles int as int64.
func ArchInt() Int {
	chillCheck()

	return Int(0)
}

// Int8 handles int8.
func Int8() Int {
	return Int(8)
}

// Int16 handles int16.
func Int16() Int {
	return Int(16)
}

// Int32 handles int32.
func Int32() Int {
	return Int(32)
}

// Int64 handles int64.
func Int64() Int {
	return Int(64)
}

// Int handles intXXX.
type Int int

// Name to implement TypeHandler.
func (i Int) Name(*renderer.Go) string {
	if i == 0 {
		return "int"
	}

	return "int" + strconv.Itoa(int(i))
}

// Pre to implement TypeHandler.
func (i Int) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (i Int) Len() int {
	return i.bytes()
}

// LenExpr to implement TypeHandler.
func (i Int) LenExpr(r *renderer.Go, src string) string {
	return strconv.Itoa(i.bytes())
}

// Encoding to implement TypeHandler.
func (i Int) Encoding(r *renderer.Go, dst, src string) {
	switch i {
	case 0:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = $bin.LittleEndian.AppendUint64($dst, uint64($src))`)
	case 8:
		r.L(`$dst = append($dst, uint8($src))`)
	default:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = $bin.LittleEndian.AppendUint$0($dst, $1($src))`, i, i.Name(r))
	}
}

// Decoding to implement TypeHandler.
func (i Int) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")

	r.L(`if len($src) < $0 {`, i.bytes())
	er.Return().New("$decode: $recordTooSmall").LenReq(i.bytes()).LenSrc().Rend(r)
	r.L(`}`)

	switch i {
	case 0:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = int($bin.LittleEndian.Uint64($src))`, i)
	case 8:
		r.L(`$dst = $0($src[0])`, i.Name(r))
	default:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = int$0($bin.LittleEndian.Uint$0($src))`, i)
	}

	return false
}

func (i Int) bytes() int {
	if i == 0 {
		return 64
	}

	return int(i >> 3)
}
