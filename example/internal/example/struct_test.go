package example

import (
	"github.com/sirkon/deepequal"
	"github.com/sirkon/errors"
	"math"
	"testing"
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
