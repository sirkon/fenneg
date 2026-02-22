package example

import (
	"encoding/binary"

	"github.com/sirkon/errors"
)

func PerfcheckDecodeManual(p *Perfcheck, src []byte) (err error) {
	// Проверяем, что буфер содержит ровно 80 байт (или хотя бы 80?)
	if len(src) != 80 {
		if len(src) < 80 {
			return errors.New("decode Perfcheck: record buffer is too small").
				Uint64("length-required", uint64(80)).
				Int("length-actual", len(src))
		}
		return errors.Newf("the buffer still has %d bytes left after decode", len(src)-80).
			Int("record-bytes-left", len(src)-80)
	}

	// Читаем все поля напрямую по фиксированным смещениям
	p.F0 = int64(binary.LittleEndian.Uint64(src[0:]))
	p.F1 = int64(binary.LittleEndian.Uint64(src[8:]))
	p.F2 = int64(binary.LittleEndian.Uint64(src[16:]))
	p.F3 = int64(binary.LittleEndian.Uint64(src[24:]))
	p.F4 = int64(binary.LittleEndian.Uint64(src[32:]))
	p.F5 = int64(binary.LittleEndian.Uint64(src[40:]))
	p.F6 = int64(binary.LittleEndian.Uint64(src[48:]))
	p.F7 = int64(binary.LittleEndian.Uint64(src[56:]))
	p.F8 = int64(binary.LittleEndian.Uint64(src[64:]))
	p.F9 = int64(binary.LittleEndian.Uint64(src[72:]))

	return nil
}
