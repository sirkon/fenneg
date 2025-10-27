package example

import (
	"encoding/binary"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirkon/errors"
	"github.com/sirkon/testlog"
)

func TestTypeRecorder(t *testing.T) {
	t.Run("fixeds", func(t *testing.T) {
		var tr TypeRecorder

		var (
			b = int32(12)
			c = uint16(13)
			d = float64(16)
			e = [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		)

		if err := tr.Fixeds(true, b, c, d, e); err != nil {
			testlog.Error(t, errors.Wrap(err, "encode"))
			return
		}

		ctrl := gomock.NewController(t)
		m := NewSourceMock(ctrl)
		m.EXPECT().Fixeds(true, b, c, d, e)

		buf, err := checkLength(tr.bufs[0])
		if err != nil {
			testlog.Error(t, errors.Wrap(err, "check buffer length"))
			return
		}

		if err := TypeRecorderDispatch(m, buf); err != nil {
			testlog.Error(t, errors.Wrap(err, "decode and dispatch"))
		}
	})

	t.Run("vars", func(t *testing.T) {
		var tr TypeRecorder

		var (
			a  = "Hello"
			b  = []byte("World!")
			c  = int16(1)
			d  = uint32(2)
			e  = EncDec(true)
			id = Index{
				Term:  1,
				Index: 2,
			}
		)

		if err := tr.Vars(a, b, c, d, e, id); err != nil {
			testlog.Error(t, errors.Wrap(err, "encode"))
		}

		ctrl := gomock.NewController(t)
		m := NewSourceMock(ctrl)
		m.EXPECT().Vars(a, b, c, d, e, id)

		buf, err := checkLength(tr.bufs[0])
		if err != nil {
			testlog.Error(t, errors.Wrap(err, "check buffer length"))
			return
		}

		if err := TypeRecorderDispatch(m, buf); err != nil {
			testlog.Error(t, errors.Wrap(err, "decode and dispatch"))
		}
	})
}

func checkLength(buf []byte) ([]byte, error) {
	l, off := binary.Uvarint(buf)
	if off <= 0 {
		if off == 0 {
			return nil, errors.New("buffer is too small")
		}

		return nil, errors.New("malformed uvarint data")
	}

	if l != uint64(len(buf)-off) {
		var e *errors.Error
		if l > uint64(len(buf)-off) {
			e = errors.New("record length encoded is larger than the length of the actual data")
		} else {
			e = errors.New("record length encoded is below the length of the actual data")
		}

		return nil, e.
			Uint64("record-length-encoded", l).
			Int("record-length-actual", len(buf)-off)
	}

	return buf[off:], nil
}
