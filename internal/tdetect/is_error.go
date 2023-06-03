package tdetect

import "go/types"

// IsErrorType checks if type is the error.
func IsErrorType(t types.Type) bool {
	n := t.(*types.Named)
	for {
		if n.Underlying() == nil {
			break
		}

		if !n.Obj().IsAlias() {
			break
		}

		u, ok := n.Underlying().(*types.Named)
		if !ok {
			break
		}

		n = u
	}

	ifc, ok := n.Underlying().(*types.Interface)
	if !ok {
		return false
	}

	if ifc.NumMethods() != 1 {
		return false
	}

	m := ifc.Method(0)
	if m.Name() != "Error" {
		return false
	}

	s := m.Type().(*types.Signature)
	if s.Params().Len() != 0 {
		return false
	}

	if s.Results().Len() != 1 {
		return false
	}

	v, ok := s.Results().At(0).Type().(*types.Basic)
	if !ok {
		return false
	}

	return v.Kind() == types.String
}
