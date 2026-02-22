package example

import (
	"encoding/binary"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/sirkon/errors"
)

var (
	perfSample = &Perfcheck{
		F0: 0,
		F1: 1,
		F2: 2,
		F3: 3,
		F4: 4,
		F5: 5,
		F6: 6,
		F7: 7,
		F8: 8,
		F9: 9,
	}

	perfSampleData = PerfcheckEncode(nil, perfSample)
)

func BenchmarkPerfcheckAppend(b *testing.B) {
	b.Run("encode", func(b *testing.B) {
		src := perfSample
		data := make([]byte, 0, PerfcheckLen(src))
		for b.Loop() {
			data = data[:0]
			data = PerfcheckEncode(data, src)
		}

		var got Perfcheck
		if err := PerfcheckDecode(&got, data); err != nil {
			b.Fatal(errors.Wrap(err, "append decode encoded data"))
		}

		assert.Equal(b, perfSample, &got, "check if encoding is right")
	})

	b.Run("decode", func(b *testing.B) {
		var got Perfcheck

		for b.Loop() {
			if err := PerfcheckDecode(&got, perfSampleData); err != nil {
				b.Fatal(errors.Wrap(err, "append decode perfSample data"))
			}
		}

		assert.Equal(b, perfSample, &got, "check if decoding is right")
	})
}

func BenchmarkPerfcheckCopy(b *testing.B) {
	b.Run("encode", func(b *testing.B) {
		src := perfSample
		data := make([]byte, 0, PerfcheckLen(src))
		for b.Loop() {
			data = data[:0]
			data = PerfcheckCopyEncode(data, src)
		}

		var got Perfcheck
		if err := PerfcheckDecode(&got, data); err != nil {
			b.Fatal(errors.Wrap(err, "copy decode encoded data"))
		}

		assert.Equal(b, perfSample, &got, "check if encoding is right")
	})

	b.Run("decode", func(b *testing.B) {
		var got Perfcheck

		for b.Loop() {
			if err := PerfcheckCopyDecode(&got, perfSampleData); err != nil {
				b.Fatal(errors.Wrap(err, "copy decode perfSample data"))
			}
		}

		assert.Equal(b, perfSample, &got, "check if decoding is right")
	})
}

// PerfcheckEncode encodes Perfcheck into dst using a single length guard
// and direct-copy semantics. The resulting layout is identical to the
// append-based version, but avoids repeated slice growth and bounds checks.
func PerfcheckCopyEncode(dst []byte, p *Perfcheck) []byte {
	if p == nil {
		return dst
	}

	const size = 8 * 10 // total bytes required

	// Ensure capacity; grow if necessary.
	if cap(dst)-len(dst) < size {
		newBuf := make([]byte, len(dst), len(dst)+size)
		copy(newBuf, dst)
		dst = newBuf
	}

	// Extend slice to fit encoded data.
	start := len(dst)
	dst = dst[:start+size]
	i := start

	// Sequential fixed-size writes.
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F0))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F1))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F2))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F3))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F4))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F5))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F6))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F7))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F8))
	i += 8
	binary.LittleEndian.PutUint64(dst[i:], uint64(p.F9))
	i += 8

	return dst
}

// PerfcheckDecode decodes content of Perfcheck using length-guard + copying semantics.
func PerfcheckCopyDecode(p *Perfcheck, src []byte) (err error) {
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
