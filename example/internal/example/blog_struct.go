package example

import (
	"time"

	"github.com/sirkon/intypes"
)

type Location struct {
	Valid bool
	File  string
	Line  intypes.VI
}

type Blog struct {
	Header   uint8
	CRC32    uint32
	Time     time.Time
	Level    uint8
	Location Location
	Payload  []byte
}
