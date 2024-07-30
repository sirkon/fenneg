// Code generated by fenneg version (devel). DO NOT EDIT.

package example

import (
	"encoding/binary"
	"math"

	"github.com/sirkon/errors"
	"github.com/sirkon/varsize"
)

const (
	sourceCodeFixeds = 1
	sourceCodeVars   = 2
)

// Fixeds encodes arguments tuple of this method.
func (t *TypeRecorder) Fixeds(a bool, b int32, c uint16, d float64, e [12]byte) error {
	var buf []byte
	{
		bufSize := 4 + 1 + 4 + 2 + 8 + 12
		buf = t.allocateBuffer(varsize.Uint(uint64(bufSize)) + bufSize)

		// Encode record length.
		buf = binary.AppendUvarint(buf, uint64(bufSize))
	}

	// Encode branch (method) code.
	buf = binary.LittleEndian.AppendUint32(buf, uint32(sourceCodeFixeds))

	// Encode a(bool).
	if a {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}

	// Encode b(int32).
	buf = binary.LittleEndian.AppendUint32(buf, uint32(b))

	// Encode c(uint16).
	buf = binary.LittleEndian.AppendUint16(buf, c)

	// Encode d(float64).
	buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(d))

	// Encode e([12]byte).
	buf = append(buf, e[:]...)

	return t.writeBuffer(buf)
}

// Vars encodes arguments tuple of this method.
func (t *TypeRecorder) Vars(a string, b []byte, c int16, d uint32, e EncDec, id Index) error {
	lenA := varsize.Uint(uint(len(a))) + len(a)
	lenB := varsize.Len(b) + len(b)
	lenC := varsize.Int(c)
	lenD := varsize.Uint(d)
	lenE := e.Len()
	var buf []byte
	{
		bufSize := 4 + lenA + lenB + lenC + lenD + lenE + 16
		buf = t.allocateBuffer(varsize.Uint(uint64(bufSize)) + bufSize)

		// Encode record length.
		buf = binary.AppendUvarint(buf, uint64(bufSize))
	}

	// Encode branch (method) code.
	buf = binary.LittleEndian.AppendUint32(buf, uint32(sourceCodeVars))

	// Encode a(string).
	buf = binary.AppendUvarint(buf, uint64(len(a)))
	buf = append(buf, a...)

	// Encode b([]byte).
	buf = binary.AppendUvarint(buf, uint64(len(b)))
	buf = append(buf, b...)

	// Encode c(int16).
	buf = binary.AppendVarint(buf, int64(c))

	// Encode d(uint32).
	buf = binary.AppendUvarint(buf, uint64(d))

	// Encode e(EncDec).
	buf = e.Encode(buf)

	// Encode id(Index).
	buf = binary.LittleEndian.AppendUint64(buf, id.Term)
	buf = binary.LittleEndian.AppendUint64(buf, id.Index)

	return t.writeBuffer(buf)
}

// TypeRecorderDispatch dispatches encoded data made with TypeRecorder
func TypeRecorderDispatch(disp Source, rec []byte) error {
	if len(rec) < 4 {
		return errors.New("decode branch code: record buffer is too small").Uint64("length-required", uint64(4)).Int("length-actual", len(rec))
	}

	branch := binary.LittleEndian.Uint32(rec[:4])
	rec = rec[4:]

	switch branch {
	case sourceCodeFixeds:
		// Decode a(bool).
		var a bool
		if len(rec) < 1 {
			return errors.New("decode Fixeds.a(bool): record buffer is too small").Uint64("length-required", uint64(1)).Int("length-actual", len(rec))
		}
		if rec[0] != 0 {
			a = true
		} else {
			a = false
		}
		rec = rec[1:]

		// Decode b(int32).
		var b int32
		if len(rec) < 4 {
			return errors.New("decode Fixeds.b(int32): record buffer is too small").Uint64("length-required", uint64(4)).Int("length-actual", len(rec))
		}
		b = int32(binary.LittleEndian.Uint32(rec))
		rec = rec[4:]

		// Decode c(uint16).
		var c uint16
		if len(rec) < 2 {
			return errors.New("decode Fixeds.c(uint16): record buffer is too small").Uint64("length-required", uint64(2)).Int("length-actual", len(rec))
		}
		c = binary.LittleEndian.Uint16(rec)
		rec = rec[2:]

		// Decode d(float64).
		var d float64
		if len(rec) >= 8 {
			keyRec := binary.LittleEndian.Uint64(rec)
			d = math.Float64frombits(keyRec)
		} else {
			return errors.New("decode Fixeds.d(float64): record buffer is too small").Uint64("length-required", uint64(8)).Int("length-actual", len(rec))
		}
		rec = rec[8:]

		// Decode e([12]byte).
		var e [12]byte
		if len(rec) < 12 {
			return errors.New("decode Fixeds.e([12]byte): record buffer is too small").Uint64("length-required", uint64(12)).Int("length-actual", len(rec))
		}
		copy(e[:12], rec)
		rec = rec[12:]

		if len(rec) > 0 {
			return errors.New("decode Fixeds: the record was not emptied after the last argument decoded").Int("record-bytes-left", len(rec))
		}

		if err := disp.Fixeds(a, b, c, d, e); err != nil {
			return errors.Wrap(err, "call Fixeds")
		}

		return nil

	case sourceCodeVars:
		// Decode a(string).
		var a string
		{
			size, off := binary.Uvarint(rec)
			if off <= 0 {
				if off == 0 {
					return errors.New("decode Vars.a(string) length: record buffer is too small")
				}
				return errors.New("decode Vars.a(string) length: malformed uvarint sequence")
			}
			rec = rec[off:]
			if int(size) > len(rec) {
				return errors.New("decode Vars.a(string) content: record buffer is too small").Uint64("length-required", uint64(int(size))).Int("length-actual", len(rec))
			}
			a = string(rec[:size])
			rec = rec[size:]
		}

		// Decode b([]byte).
		var b []byte
		{
			size, off := binary.Uvarint(rec)
			if off <= 0 {
				if off == 0 {
					return errors.New("decode Vars.b([]byte) length: record buffer is too small")
				}
				return errors.New("decode Vars.b([]byte) length: malformed uvarint sequence")
			}
			rec = rec[off:]
			if uint64(len(rec)) < size {
				return errors.New("decode Vars.b([]byte) content: record buffer is too small").Uint64("length-required", uint64(size)).Int("length-actual", len(rec))
			}
			b = rec[:size]
			rec = rec[size:]
		}

		// Decode c(int16).
		var c int16
		{
			val, off := binary.Varint(rec)
			if off <= 0 {
				if off == 0 {
					return errors.New("decode Vars.c(int16): record buffer is too small")
				}
				return errors.New("decode Vars.c(int16): malformed varint sequence")
			}
			c = int16(val)
			rec = rec[off:]
		}

		// Decode d(uint32).
		var d uint32
		{
			val, off := binary.Uvarint(rec)
			if off <= 0 {
				if off == 0 {
					return errors.New("decode Vars.d(uint32): record buffer is too small")
				}
				return errors.New("decode Vars.d(uint32): malformed uvarint sequence")
			}
			d = uint32(val)
			rec = rec[off:]
		}

		// Decode e(EncDec).
		var e EncDec
		if recRest, err := e.Decode(rec); err != nil {
			return errors.Wrap(err, "decode Vars.e(EncDec)")
		} else {
			rec = recRest
		}

		// Decode id(Index).
		var id Index
		if len(rec) < 16 {
			return errors.New("decode Vars.id(Index): record buffer is too small").Uint64("length-required", uint64(16)).Int("length-actual", len(rec))
		}
		id.Term = binary.LittleEndian.Uint64(rec)
		id.Index = binary.LittleEndian.Uint64(rec[8:])
		rec = rec[16:]

		if len(rec) > 0 {
			return errors.New("decode Vars: the record was not emptied after the last argument decoded").Int("record-bytes-left", len(rec))
		}

		if err := disp.Vars(a, b, c, d, e, id); err != nil {
			return errors.Wrap(err, "call Vars")
		}

		return nil

	default:
		return errors.Newf("invalid branch code %d", branch).Uint32("invalid-branch-code", branch)
	}

	return nil
}
