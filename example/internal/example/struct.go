package example

import "github.com/sirkon/intypes"

type Struct struct {
	ID          Index
	ChangeID    Index
	Repeat      uint32
	Theme       uint32
	Data        []byte
	Pi          float32
	E           float64
	Field       string
	Int         int
	Uint        uint
	VarInt      intypes.VI
	VarUint     intypes.VU
	BoolSlice   []bool
	StringSlice []string
}
