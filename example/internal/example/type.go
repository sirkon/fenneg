package example

// TypeRecorder implement
type TypeRecorder struct {
	bufs [][]byte
}

func (t *TypeRecorder) allocateBuffer(n int) []byte {
	return make([]byte, 0, n)
}

func (t *TypeRecorder) writeBuffer(buf []byte) error {
	t.bufs = append(t.bufs, buf)
	return nil
}
