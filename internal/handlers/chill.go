package handlers

var keepItChill bool

// Chill call this if you want to support int and uint types treated as int64 and uint64 respectively.
func Chill() {
	keepItChill = true
}

func chillCheck() {
	if !keepItChill {
		panic("make fenneg.Chill() call to support int and uint types")
	}
}
