// Code generated by fenneg version (devel). DO NOT EDIT.

package example

import (
	"encoding/binary"
	"math"

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
	// len Pi(float32).
	// len E(float64).
	// len Field(string).
	// len Int(int).
	// len Uint(uint).
	// len VarInt(int).
	// len VarUint(uint).
	// len BoolSlice([]bool).
	// len StringSlice([]string).
	// len BoolSliceSlice([][]bool).
	// len StringSliceSlice([][]string).

	lenData := varsize.Len(s.Data) + len(s.Data)
	lenField := varsize.Uint(uint(len(s.Field))) + len(s.Field)
	lenVarInt := varsize.Int(s.VarInt)
	lenVarUint := varsize.Uint(s.VarUint)
	lenBoolSlice := varsize.Len(s.BoolSlice) + len(s.BoolSlice)*1
	lenStringSlice := varsize.Len(s.StringSlice)
	for _, item := range s.StringSlice {
		lenItem := varsize.Uint(uint(len(item))) + len(item)
		lenStringSlice += lenItem
	}
	lenBoolSliceSlice := varsize.Len(s.BoolSliceSlice)
	for _, item := range s.BoolSliceSlice {
		lenItem := varsize.Len(item) + len(item)*1
		lenBoolSliceSlice += lenItem
	}
	lenStringSliceSlice := varsize.Len(s.StringSliceSlice)
	for _, item := range s.StringSliceSlice {
		lenItem := varsize.Len(item)
		for _, item2 := range item {
			lenItem2 := varsize.Uint(uint(len(item2))) + len(item2)
			lenItem += lenItem2
		}
		lenStringSliceSlice += lenItem
	}

	return 16 + 16 + 4 + 4 + lenData + 4 + 8 + lenField + 8 + 8 + lenVarInt + lenVarUint + lenBoolSlice + lenStringSlice + lenBoolSliceSlice + lenStringSliceSlice
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

	// Encode Pi(float32).
	dst = binary.LittleEndian.AppendUint32(dst, math.Float32bits(s.Pi))

	// Encode E(float64).
	dst = binary.LittleEndian.AppendUint64(dst, math.Float64bits(s.E))

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
	dst = binary.AppendUvarint(dst, uint64(s.VarUint))

	// Encode BoolSlice([]bool).
	dst = binary.AppendUvarint(dst, uint64(len(s.BoolSlice)))
	for _, v := range s.BoolSlice {
		if v {
			dst = append(dst, 1)
		} else {
			dst = append(dst, 0)
		}
	}

	// Encode StringSlice([]string).
	dst = binary.AppendUvarint(dst, uint64(len(s.StringSlice)))
	for _, v := range s.StringSlice {
		dst = binary.AppendUvarint(dst, uint64(len(v)))
		dst = append(dst, v...)
	}

	// Encode BoolSliceSlice([][]bool).
	dst = binary.AppendUvarint(dst, uint64(len(s.BoolSliceSlice)))
	for _, v := range s.BoolSliceSlice {
		dst = binary.AppendUvarint(dst, uint64(len(v)))
		for _, vV := range v {
			if vV {
				dst = append(dst, 1)
			} else {
				dst = append(dst, 0)
			}
		}
	}

	// Encode StringSliceSlice([][]string).
	dst = binary.AppendUvarint(dst, uint64(len(s.StringSliceSlice)))
	for _, v := range s.StringSliceSlice {
		dst = binary.AppendUvarint(dst, uint64(len(v)))
		for _, vV := range v {
			dst = binary.AppendUvarint(dst, uint64(len(vV)))
			dst = append(dst, vV...)
		}
	}

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

	// Decode Pi(float32).
	if len(src) >= 4 {
		keySrc := binary.LittleEndian.Uint32(src)
		s.Pi = math.Float32frombits(keySrc)
	} else {
		return errors.New("decode s.Pi(float32): record buffer is too small").Uint64("length-required", uint64(4)).Int("length-actual", len(src))
	}
	src = src[4:]

	// Decode E(float64).
	if len(src) >= 8 {
		keySrc := binary.LittleEndian.Uint64(src)
		s.E = math.Float64frombits(keySrc)
	} else {
		return errors.New("decode s.E(float64): record buffer is too small").Uint64("length-required", uint64(8)).Int("length-actual", len(src))
	}
	src = src[8:]

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
	if len(src) < 8 {
		return errors.New("decode s.Int(int): record buffer is too small").Uint64("length-required", uint64(8)).Int("length-actual", len(src))
	}
	s.Int = int(binary.LittleEndian.Uint64(src))
	src = src[8:]

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
	{
		val, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.VarUint(uint): record buffer is too small")
			}
			return errors.New("decode s.VarUint(uint): malformed uvarint sequence")
		}
		s.VarUint = uint(val)
		src = src[off:]
	}

	// Decode BoolSlice([]bool).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.BoolSlice([]bool) length: record buffer is too small")
			}
			return errors.New("decode s.BoolSlice([]bool) length: malformed uvarint sequence")
		}
		src = src[off:]
		s.BoolSlice = make([]bool, size)
		for i := 0; i < int(size); i++ {
			if len(src) < 1 {
				return errors.New("decode s.BoolSlice[i]([]bool): record buffer is too small").Uint64("length-required", uint64(1)).Int("length-actual", len(src))
			}
			if src[0] != 0 {
				s.BoolSlice[i] = true
			}
			src = src[1:]
		}
	}

	// Decode StringSlice([]string).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.StringSlice([]string) length: record buffer is too small")
			}
			return errors.New("decode s.StringSlice([]string) length: malformed uvarint sequence")
		}
		src = src[off:]
		s.StringSlice = make([]string, size)
		for i := 0; i < int(size); i++ {
			{
				size2, off2 := binary.Uvarint(src)
				if off2 <= 0 {
					if off2 == 0 {
						return errors.New("decode s.StringSlice[i]([]string) length: record buffer is too small")
					}
					return errors.New("decode s.StringSlice[i]([]string) length: malformed uvarint sequence")
				}
				src = src[off2:]
				if int(size2) > len(src) {
					return errors.New("decode s.StringSlice[i]([]string) content: record buffer is too small").Uint64("length-required", uint64(int(size2))).Int("length-actual", len(src))
				}
				s.StringSlice[i] = string(src[:size2])
				src = src[size2:]
			}
		}
	}

	// Decode BoolSliceSlice([][]bool).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.BoolSliceSlice([][]bool) length: record buffer is too small")
			}
			return errors.New("decode s.BoolSliceSlice([][]bool) length: malformed uvarint sequence")
		}
		src = src[off:]
		s.BoolSliceSlice = make([][]bool, size)
		for i := 0; i < int(size); i++ {
			{
				size2, off2 := binary.Uvarint(src)
				if off2 <= 0 {
					if off2 == 0 {
						return errors.New("decode s.BoolSliceSlice[i]([][]bool) length: record buffer is too small")
					}
					return errors.New("decode s.BoolSliceSlice[i]([][]bool) length: malformed uvarint sequence")
				}
				src = src[off2:]
				s.BoolSliceSlice[i] = make([]bool, size2)
				for i2 := 0; i2 < int(size2); i2++ {
					if len(src) < 1 {
						return errors.New("decode s.BoolSliceSlice[i][i2]([][]bool): record buffer is too small").Uint64("length-required", uint64(1)).Int("length-actual", len(src))
					}
					if src[0] != 0 {
						s.BoolSliceSlice[i][i2] = true
					}
					src = src[1:]
				}
			}
		}
	}

	// Decode StringSliceSlice([][]string).
	{
		size, off := binary.Uvarint(src)
		if off <= 0 {
			if off == 0 {
				return errors.New("decode s.StringSliceSlice([][]string) length: record buffer is too small")
			}
			return errors.New("decode s.StringSliceSlice([][]string) length: malformed uvarint sequence")
		}
		src = src[off:]
		s.StringSliceSlice = make([][]string, size)
		for i := 0; i < int(size); i++ {
			{
				size2, off2 := binary.Uvarint(src)
				if off2 <= 0 {
					if off2 == 0 {
						return errors.New("decode s.StringSliceSlice[i]([][]string) length: record buffer is too small")
					}
					return errors.New("decode s.StringSliceSlice[i]([][]string) length: malformed uvarint sequence")
				}
				src = src[off2:]
				s.StringSliceSlice[i] = make([]string, size2)
				for i2 := 0; i2 < int(size2); i2++ {
					{
						size3, off3 := binary.Uvarint(src)
						if off3 <= 0 {
							if off3 == 0 {
								return errors.New("decode s.StringSliceSlice[i][i2]([][]string) length: record buffer is too small")
							}
							return errors.New("decode s.StringSliceSlice[i][i2]([][]string) length: malformed uvarint sequence")
						}
						src = src[off3:]
						if int(size3) > len(src) {
							return errors.New("decode s.StringSliceSlice[i][i2]([][]string) content: record buffer is too small").Uint64("length-required", uint64(int(size3))).Int("length-actual", len(src))
						}
						s.StringSliceSlice[i][i2] = string(src[:size3])
						src = src[size3:]
					}
				}
			}
		}
	}

	if len(src) > 0 {
		return errors.Newf("the buffer still has %d bytes left after the last argument decoded", len(src)).Int("record-bytes-left", len(src))
	}

	return nil
}
