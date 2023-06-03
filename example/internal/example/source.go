package example

import (
	"github.com/sirkon/errors"
	"github.com/sirkon/intypes"
)

// Source interface to generate testing stuff.
type Source interface {
	Fixeds(a bool, b int32, c uint16, d float64, e [12]byte) error
	Vars(a string, b []byte, c intypes.VI16, d intypes.VU32, e EncDec, id Index) error
}

// EncDec auto-apply type.
type EncDec bool

// Len to auto-apply encoding.
func (x EncDec) Len() int {
	return 1
}

// Encode to auto-apply encoding.
func (x EncDec) Encode(buf []byte) []byte {
	if x {
		return append(buf, 1)
	}

	return append(buf, 0)
}

// Decode to auto-apply decoding.
func (x *EncDec) Decode(rec []byte) ([]byte, error) {
	if len(rec) < 1 {
		return nil, errors.New("the record buffer is too small").
			Int("length-required", 1).
			Int("length-actual", 0)
	}

	*x = rec[0] != 0
	return rec[1:], nil
}

// Index to be handled by custom type handler.
type Index struct {
	Term  uint64
	Index uint64
}
