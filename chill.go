package fenneg

import "github.com/sirkon/fenneg/internal/handlers"

// Chill should be called to support int and uint types as int64 and uint64 respectively.
func Chill() {
	handlers.Chill()
}
