package handlers

import (
	"strconv"

	"github.com/sirkon/fenneg/internal/er"
	"github.com/sirkon/fenneg/internal/renderer"
)

// ArchUint handles uint as uint64.
func ArchUint() Uint {
	chillCheck()

	return Uint(0)
}

// Uint8 handles uint8.
func Uint8() Uint {
	return Uint(8)
}

// Uint16 handles uint16.
func Uint16() Uint {
	return Uint(16)
}

// Uint32 handles uint32.
func Uint32() Uint {
	return Uint(32)
}

// Uint64 handles uint64.
func Uint64() Uint {
	return Uint(64)
}

// Uint handles uintXXX.
type Uint int

// Name to implement TypeHandler.
func (i Uint) Name(*renderer.Go) string {
	if i == 0 {
		return "uint"
	}

	return "uint" + strconv.Itoa(int(i))
}

// Pre to implement TypeHandler.
func (i Uint) Pre(r *renderer.Go, src string) {}

// Len to implement TypeHandler.
func (i Uint) Len() int {
	return i.bytes()
}

// LenExpr to implement TypeHandler.
func (i Uint) LenExpr(r *renderer.Go, src string) string {
	return strconv.Itoa(i.bytes())
}

// Encoding to implement TypeHandler.
func (i Uint) Encoding(r *renderer.Go, dst, src string) {
	switch i {
	case 0:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = $bin.LittleEndian.AppendUint64($dst, uint64($src))`)
	case 8:
		r.L(`$dst = append($dst, $src)`)
	default:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = $bin.LittleEndian.AppendUint$0($dst, $src)`, i)
	}
}

// Decoding to implement TypeHandler.
func (i Uint) Decoding(r *renderer.Go, dst, src string) bool {
	r.Imports().Errors().Ref("errors")

	r.L(`if len($src) < $0 {`, i.bytes())
	er.Return().New("$decode: $recordTooSmall").LenReq(i.bytes()).LenSrc().Rend(r)
	r.L(`}`)

	switch i {
	case 0:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = uint($bin.LittleEndian.Uint64($src))`, i)
	case 8:
		r.L(`$dst = $src[0]`)
	default:
		r.Imports().Binary().Ref("bin")
		r.L(`$dst = $bin.LittleEndian.Uint$0($src)`, i)
	}

	return false
}

func (i Uint) bytes() int {
	if i == 0 {
		return 8
	}

	return int(i >> 3)
}
