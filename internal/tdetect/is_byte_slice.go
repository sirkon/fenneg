package tdetect

import "go/types"

// IsByteSlice checks if the type represents []byte
func IsByteSlice(v types.Type) bool {
	vv, ok := v.(*types.Slice)
	if !ok {
		return false
	}

	return IsBasic(vv.Elem(), types.Byte)
}
