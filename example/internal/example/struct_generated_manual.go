package example

import (
	"encoding/binary"
	"math"

	"github.com/sirkon/errors"
)

func StructDecodeManual(s *Struct, src []byte) error {
	pos := 0

	// Блок 1: ID (16) + ChangeID (16) + Repeat (4) + Theme (4) = 40 байт
	if pos+40 > len(src) {
		return errors.New("decode first block: buffer too small")
	}
	tmp := src[pos : pos+40]
	s.ID.Term = binary.LittleEndian.Uint64(tmp[0:8])
	s.ID.Index = binary.LittleEndian.Uint64(tmp[8:16])
	s.ChangeID.Term = binary.LittleEndian.Uint64(tmp[16:24])
	s.ChangeID.Index = binary.LittleEndian.Uint64(tmp[24:32])
	s.Repeat = binary.LittleEndian.Uint32(tmp[32:36])
	s.Theme = binary.LittleEndian.Uint32(tmp[36:40])
	pos += 40

	// Data []byte
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.Data length: buffer too small")
			}
			return errors.New("decode s.Data length: malformed uvarint")
		}
		pos += n
		if pos+int(size) > len(src) {
			return errors.New("decode s.Data content: buffer too small")
		}
		s.Data = src[pos : pos+int(size)]
		pos += int(size)
	}

	// Блок 2: Pi (4) + E (8) = 12 байт
	if pos+12 > len(src) {
		return errors.New("decode second block: buffer too small")
	}
	tmp = src[pos : pos+12]
	s.Pi = math.Float32frombits(binary.LittleEndian.Uint32(tmp[0:4]))
	s.E = math.Float64frombits(binary.LittleEndian.Uint64(tmp[4:12]))
	pos += 12

	// Field string — инлайн, без вызова функции
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.Field length: buffer too small")
			}
			return errors.New("decode s.Field length: malformed uvarint")
		}
		pos += n
		if pos+int(size) > len(src) {
			return errors.New("decode s.Field content: buffer too small")
		}
		s.Field = string(src[pos : pos+int(size)])
		pos += int(size)
	}

	// Блок 3: Int (8) + Uint (8) = 16 байт
	if pos+16 > len(src) {
		return errors.New("decode third block: buffer too small")
	}
	tmp = src[pos : pos+16]
	s.Int = int(binary.LittleEndian.Uint64(tmp[0:8]))
	s.Uint = uint(binary.LittleEndian.Uint64(tmp[8:16]))
	pos += 16

	// VarInt
	{
		val, n := binary.Varint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.VarInt: buffer too small")
			}
			return errors.New("decode s.VarInt: malformed varint")
		}
		s.VarInt = int(val)
		pos += n
	}

	// VarUint
	{
		val, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.VarUint: buffer too small")
			}
			return errors.New("decode s.VarUint: malformed uvarint")
		}
		s.VarUint = uint(val)
		pos += n
	}

	// BoolSlice — каждый bool занимает 1 байт
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.BoolSlice length: buffer too small")
			}
			return errors.New("decode s.BoolSlice length: malformed uvarint")
		}
		pos += n
		if pos+int(size) > len(src) {
			return errors.New("decode s.BoolSlice content: buffer too small")
		}
		s.BoolSlice = make([]bool, size)
		for i := 0; i < int(size); i++ {
			s.BoolSlice[i] = src[pos] != 0
			pos++
		}
	}

	// StringSlice — строки переменной длины, всё инлайн
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.StringSlice length: buffer too small")
			}
			return errors.New("decode s.StringSlice length: malformed uvarint")
		}
		pos += n
		s.StringSlice = make([]string, size)
		for i := 0; i < int(size); i++ {
			// Читаем длину очередной строки
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StringSlice[i] length: buffer too small")
				}
				return errors.New("decode s.StringSlice[i] length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.StringSlice[i] content: buffer too small")
			}
			s.StringSlice[i] = string(src[pos : pos+int(sz)])
			pos += int(sz)
		}
	}

	// BoolSliceSlice — двумерный слайс bool
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.BoolSliceSlice length: buffer too small")
			}
			return errors.New("decode s.BoolSliceSlice length: malformed uvarint")
		}
		pos += n
		s.BoolSliceSlice = make([][]bool, size)
		for i := 0; i < int(size); i++ {
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.BoolSliceSlice[i] length: buffer too small")
				}
				return errors.New("decode s.BoolSliceSlice[i] length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.BoolSliceSlice[i] content: buffer too small")
			}
			inner := make([]bool, sz)
			for j := 0; j < int(sz); j++ {
				inner[j] = src[pos] != 0
				pos++
			}
			s.BoolSliceSlice[i] = inner
		}
	}

	// StringSliceSlice — двумерный слайс строк
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.StringSliceSlice length: buffer too small")
			}
			return errors.New("decode s.StringSliceSlice length: malformed uvarint")
		}
		pos += n
		s.StringSliceSlice = make([][]string, size)
		for i := 0; i < int(size); i++ {
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StringSliceSlice[i] length: buffer too small")
				}
				return errors.New("decode s.StringSliceSlice[i] length: malformed uvarint")
			}
			pos += n2
			s.StringSliceSlice[i] = make([]string, sz)
			for j := 0; j < int(sz); j++ {
				sz2, n3 := binary.Uvarint(src[pos:])
				if n3 <= 0 {
					if n3 == 0 {
						return errors.New("decode s.StringSliceSlice[i][j] length: buffer too small")
					}
					return errors.New("decode s.StringSliceSlice[i][j] length: malformed uvarint")
				}
				pos += n3
				if pos+int(sz2) > len(src) {
					return errors.New("decode s.StringSliceSlice[i][j] content: buffer too small")
				}
				s.StringSliceSlice[i][j] = string(src[pos : pos+int(sz2)])
				pos += int(sz2)
			}
		}
	}

	// MapKFVF (map[uint32]uint32) — можно читать блоком пар
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.MapKFVF length: buffer too small")
			}
			return errors.New("decode s.MapKFVF length: malformed uvarint")
		}
		pos += n
		// Проверяем весь блок пар сразу
		if pos+8*int(size) > len(src) {
			return errors.New("decode s.MapKFVF content: buffer too small")
		}
		tmp = src[pos : pos+8*int(size)]
		s.MapKFVF = make(map[uint32]uint32, size)
		for i := 0; i < int(size); i++ {
			off := i * 8
			k := binary.LittleEndian.Uint32(tmp[off : off+4])
			v := binary.LittleEndian.Uint32(tmp[off+4 : off+8])
			s.MapKFVF[k] = v
		}
		pos += 8 * int(size)
	}

	// MapKFVV (map[uint32]string)
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.MapKFVV length: buffer too small")
			}
			return errors.New("decode s.MapKFVV length: malformed uvarint")
		}
		pos += n
		s.MapKFVV = make(map[uint32]string, size)
		for i := 0; i < int(size); i++ {
			// ключ uint32
			if pos+4 > len(src) {
				return errors.New("decode s.MapKFVV key: buffer too small")
			}
			k := binary.LittleEndian.Uint32(src[pos:])
			pos += 4
			// значение строка
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.MapKFVV value length: buffer too small")
				}
				return errors.New("decode s.MapKFVV value length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.MapKFVV value content: buffer too small")
			}
			v := string(src[pos : pos+int(sz)])
			pos += int(sz)
			s.MapKFVV[k] = v
		}
	}

	// MapKVVF (map[string]uint32)
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.MapKVVF length: buffer too small")
			}
			return errors.New("decode s.MapKVVF length: malformed uvarint")
		}
		pos += n
		s.MapKVVF = make(map[string]uint32, size)
		for i := 0; i < int(size); i++ {
			// ключ строка
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.MapKVVF key length: buffer too small")
				}
				return errors.New("decode s.MapKVVF key length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.MapKVVF key content: buffer too small")
			}
			k := string(src[pos : pos+int(sz)])
			pos += int(sz)
			// значение uint32
			if pos+4 > len(src) {
				return errors.New("decode s.MapKVVF value: buffer too small")
			}
			v := binary.LittleEndian.Uint32(src[pos:])
			pos += 4
			s.MapKVVF[k] = v
		}
	}

	// Struct (StructInternal)
	{
		if pos+8 > len(src) {
			return errors.New("decode s.Struct.A: buffer too small")
		}
		s.Struct.A = int(binary.LittleEndian.Uint64(src[pos:]))
		pos += 8
		// B строка
		sz, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.Struct.B length: buffer too small")
			}
			return errors.New("decode s.Struct.B length: malformed uvarint")
		}
		pos += n
		if pos+int(sz) > len(src) {
			return errors.New("decode s.Struct.B content: buffer too small")
		}
		s.Struct.B = string(src[pos : pos+int(sz)])
		pos += int(sz)
	}

	// StructSlice
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.StructSlice length: buffer too small")
			}
			return errors.New("decode s.StructSlice length: malformed uvarint")
		}
		pos += n
		s.StructSlice = make([]StructInternal, size)
		for i := 0; i < int(size); i++ {
			if pos+8 > len(src) {
				return errors.New("decode s.StructSlice[i].A: buffer too small")
			}
			s.StructSlice[i].A = int(binary.LittleEndian.Uint64(src[pos:]))
			pos += 8
			// B строка
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StructSlice[i].B length: buffer too small")
				}
				return errors.New("decode s.StructSlice[i].B length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.StructSlice[i].B content: buffer too small")
			}
			s.StructSlice[i].B = string(src[pos : pos+int(sz)])
			pos += int(sz)
		}
	}

	// StructMapInt
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.StructMapInt length: buffer too small")
			}
			return errors.New("decode s.StructMapInt length: malformed uvarint")
		}
		pos += n
		s.StructMapInt = make(map[int]StructInternal, size)
		for i := 0; i < int(size); i++ {
			if pos+8 > len(src) {
				return errors.New("decode s.StructMapInt key: buffer too small")
			}
			k := int(binary.LittleEndian.Uint64(src[pos:]))
			pos += 8
			// значение
			if pos+8 > len(src) {
				return errors.New("decode s.StructMapInt value.A: buffer too small")
			}
			v := StructInternal{
				A: int(binary.LittleEndian.Uint64(src[pos:])),
			}
			pos += 8
			// B строка
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StructMapInt value.B length: buffer too small")
				}
				return errors.New("decode s.StructMapInt value.B length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.StructMapInt value.B content: buffer too small")
			}
			v.B = string(src[pos : pos+int(sz)])
			pos += int(sz)
			s.StructMapInt[k] = v
		}
	}

	// StructMapStruct
	{
		size, n := binary.Uvarint(src[pos:])
		if n <= 0 {
			if n == 0 {
				return errors.New("decode s.StructMapStruct length: buffer too small")
			}
			return errors.New("decode s.StructMapStruct length: malformed uvarint")
		}
		pos += n
		s.StructMapStruct = make(map[StructInternal]StructInternal, size)
		for i := 0; i < int(size); i++ {
			// key
			if pos+8 > len(src) {
				return errors.New("decode s.StructMapStruct key.A: buffer too small")
			}
			k := StructInternal{
				A: int(binary.LittleEndian.Uint64(src[pos:])),
			}
			pos += 8
			sz, n2 := binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StructMapStruct key.B length: buffer too small")
				}
				return errors.New("decode s.StructMapStruct key.B length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.StructMapStruct key.B content: buffer too small")
			}
			k.B = string(src[pos : pos+int(sz)])
			pos += int(sz)

			// value
			if pos+8 > len(src) {
				return errors.New("decode s.StructMapStruct value.A: buffer too small")
			}
			v := StructInternal{
				A: int(binary.LittleEndian.Uint64(src[pos:])),
			}
			pos += 8
			sz, n2 = binary.Uvarint(src[pos:])
			if n2 <= 0 {
				if n2 == 0 {
					return errors.New("decode s.StructMapStruct value.B length: buffer too small")
				}
				return errors.New("decode s.StructMapStruct value.B length: malformed uvarint")
			}
			pos += n2
			if pos+int(sz) > len(src) {
				return errors.New("decode s.StructMapStruct value.B content: buffer too small")
			}
			v.B = string(src[pos : pos+int(sz)])
			pos += int(sz)

			s.StructMapStruct[k] = v
		}
	}

	if pos != len(src) {
		return errors.Newf("extra bytes left: %d", len(src)-pos)
	}
	return nil
}
