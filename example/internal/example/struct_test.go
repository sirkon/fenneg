package example

import (
	"math"
	"testing"

	"github.com/sirkon/deepequal"
	"github.com/sirkon/errors"
)

func TestStruct(t *testing.T) {
	sample := Struct{
		ID: Index{
			Term:  1,
			Index: 2,
		},
		ChangeID: Index{
			Term:  3,
			Index: 4,
		},
		Repeat:      1,
		Theme:       2,
		Data:        []byte{1, 2, 3},
		Pi:          math.Pi,
		E:           math.E,
		Field:       "field",
		Int:         1,
		Uint:        2,
		VarInt:      3,
		VarUint:     4,
		BoolSlice:   []bool{true, false, true},
		StringSlice: []string{"a", "b", "c"},
		BoolSliceSlice: [][]bool{
			{
				true,
				false,
			},
			{
				false,
				true,
			},
		},
		StringSliceSlice: [][]string{
			{"a", "b", "c"},
			{"c", "b", "a"},
		},
		MapKFVF: map[uint32]uint32{1: 2, 2: 3},
		MapKFVV: map[uint32]string{1: "1", 2: "12"},
		MapKVVF: map[string]uint32{"1": 1, "12": 12},
		Struct: StructInternal{
			A: 123,
			B: "Hello!",
		},
		StructSlice: []StructInternal{
			{
				A: 12,
				B: "Hello",
			},
			{
				A: 78,
				B: "World",
			},
		},
		StructMapInt: map[int]StructInternal{
			1: {
				A: 21,
				B: "olleH",
			},
			300: {
				A: 87,
				B: "dlorW",
			},
		},
		StructMapStruct: map[StructInternal]StructInternal{
			{
				A: 12,
				B: "Hello",
			}: {
				A: 21,
				B: "olleH",
			},
			{
				A: 78,
				B: "World",
			}: {
				A: 87,
				B: "dlroW",
			},
		},
	}

	var buf []byte
	buf = StructEncode(buf, &sample)

	var res Struct
	if err := StructDecode(&res, buf); err != nil {
		t.Fatal(errors.Wrap(err, "decode struct"))
	}

	if !deepequal.Equal(res, sample) {
		deepequal.SideBySide(t, "struct", sample, res)
	} else {
		deepequal.SideBySide(t, "struct", sample, res)
	}
}
