package tdetect

import "go/types"

// IsBasic checks if the type is the basic one with
// the kind given.
func IsBasic(t types.Type, kind types.BasicKind) bool {
	v, ok := t.(*types.Basic)
	if !ok {
		return false
	}

	return v.Kind() == kind
}
