// Code generated by fenneg version (devel). DO NOT EDIT.

package example

import (
	"encoding/binary"

	"github.com/sirkon/errors"
	"github.com/sirkon/varsize"
)

func StructLen(s *Struct) int {
	if s == nil {
		return 0
	}

	// len ID(Index).
	// len ChangeID(Index).
	// len Repeat(uint32).
	// len Theme(uint32).
	// len Data([]byte).
	// len Field(string).
	// len Int(int).
	// len Uint(uint).
	// len VarInt(int).
	// len VarUint(uint).

	lenData := varsize.Len(s.Data) + len(s.Data)
	lenField := varsize.Uint(uint(len(s.Field))) + len(s.Field)
	lenVarInt := varsize.Int(s.VarInt)

	return 16 + 16 + 4 + 4 + lenData + lenField + 64 + 8 + lenVarInt + 8
}
func StructEncode(dst []byte, s *Struct) []byte {
	if s == nil {
		return dst
	}
	if dst == nil {
		dst = make([]byte, 0, StructLen(s))
	}

	// Encode ID(Index).
	dst = binary.LittleEndian.AppendUint64(dst, s.ID.Term)
	dst = binary.LittleEndian.AppendUint64(dst, s.ID.Index)

	// Encode ChangeID(Index).
	dst = binary.LittleEndian.AppendUint64(dst, s.ChangeID.Term)
	dst = binary.LittleEndian.AppendUint64(dst, s.ChangeID.Index)

	// Encode Repeat(uint32).
	dst = binary.LittleEndian.AppendUint32(dst, s.Repeat)

	// Encode Theme(uint32).
	dst = binary.LittleEndian.AppendUint32(dst, s.Theme)

	// Encode Data([]byte).
	dst = binary.AppendUvarint(dst, uint64(len(s.Data)))
	dst = append(dst, s.Data...)

	// Encode Field(string).
	dst = binary.AppendUvarint(dst, uint64(len(s.Field)))
	dst = append(dst, s.Field...)

	// Encode Int(int).
	dst = binary.LittleEndian.AppendUint64(dst, uint64(s.Int))

	// Encode Uint(uint).
	dst = binary.LittleEndian.AppendUint64(dst, uint64(s.Uint))

	// Encode VarInt(int).
	dst = binary.AppendVarint(dst, int64(s.VarInt))

	// Encode VarUint(uint).
	dst = binary.LittleEndian.AppendUint64(dst, uint64(s.VarUint))

	return dst
}

// StructEncode decodes content of Struct.
func StructDecode(s *Struct, src []byte) (err error) {
	// Decode ID(Index).
	if len(src) < 16 {
		return errors.New("decode s.ID(Index): record buffer is too small").Uint64("length-required", uint64(16)).Int("length-actual", len(src))
	}
	s.ID.Term = binary.LittleEndian.Uint64(src)
	s.ID.Index = binary.LittleEndian.Uint64(src[8:])
	src = src[16:]

	// Decode ChangeID(Index).
	if len(src) < 16 {
		return errors.New("decode s.ChangeID(Index): record buffer is too small").Uint64("length-required", uint64(16)).Int("length-actual", len(src))
	}
	s.ChangeID.Term = binary.LittleEndian.Uint64(src)
	s.ChangeID.Index = binary.LittleEndian.Uint64(src[8:])
	src = src[16:]

	// Decode Repeat(uint32).
	if len(src) < 4 {
		return errors.New("decode s.Repeat(uint32): record buffer is too small").Uint64("length-required", uint64(4)).Int("length-actual", len(src))
	}
	s.Repeat = binary.LittleEndian.Uint32(src)
	src = src[4:]

	// Decode Theme(uint32).
	if len(src) < 4 {
		return errors.New("decode s.Theme(uint32): record buffer is too small").Uint64("length-required", uint64(4)).Int("length-actual", len(src))
	}
	s.Theme = binary.LittleEndian.Uint32(src)
	src = src[4:]

	// Decode Data([]byte).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.Data([]byte) length: record buffer is too small")
			}
			return errors.New("decode s.Data([]byte) length: malformed uvarint sequence")
		}
		src = src[off:]
		if uint64(len(src)) < size {
			return errors.New("decode s.Data([]byte) content: record buffer is too small").Uint64("length-required", uint64(size)).Int("length-actual", len(src))
		}
		s.Data = src[:size]
		src = src[size:]
	}

	// Decode Field(string).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.Field(string) length: record buffer is too small")
			}
			return errors.New("decode s.Field(string) length: malformed uvarint sequence")
		}
		src = src[off:]
		if int(size) > len(src) {
			return errors.New("decode s.Field(string) content: record buffer is too small").Uint64("length-required", uint64(int(size))).Int("length-actual", len(src))
		}
		s.Field = string(src[:size])
		src = src[size:]
	}

	// Decode Int(int).
	if len(src) < 64 {
		return errors.New("decode s.Int(int): record buffer is too small").Uint64("length-required", uint64(64)).Int("length-actual", len(src))
	}
	s.Int = int(binary.LittleEndian.Uint64(src))
	src = src[64:]

	// Decode Uint(uint).
	if len(src) < 8 {
		return errors.New("decode s.Uint(uint): record buffer is too small").Uint64("length-required", uint64(8)).Int("length-actual", len(src))
	}
	s.Uint = uint(binary.LittleEndian.Uint64(src))
	src = src[8:]

	// Decode VarInt(int).
	{
		val, off := binary.Varint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.VarInt(int): record buffer is too small")
			}
			return errors.New("decode s.VarInt(int): malformed varint sequence")
		}
		s.VarInt = int(val)
		src = src[off:]
	}

	// Decode VarUint(uint).
	if len(src) < 8 {
		return errors.New("decode s.VarUint(uint): record buffer is too small").Uint64("length-required", uint64(8)).Int("length-actual", len(src))
	}
	s.VarUint = uint(binary.LittleEndian.Uint64(src))
	src = src[8:]

	if len(src) > 0 {
		return errors.New("the record was not emptied after the last argument decoded").Int("record-bytes-left", len(src))
	}

	return nil
}
