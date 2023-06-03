package tdetect

import "go/types"

// TypeImplementsEncoder checks if type implements encoder interface:
//
//     type Encoder interface{
//         Len() int
//         Encode([]byte) error
//     }
func TypeImplementsEncoder(t types.Type) bool {
	return typeImplementsEncoder(t, 0)
}

// TypeImplementsDecoder checks if type implements decoder interface:
//
//     type Decoder interface{
//         Decode([]byte) ([]byte, error)
//     }
func TypeImplementsDecoder(t types.Type) bool {
	return typeImplementsDecoder(t, 0)
}

func typeImplementsEncoder(t types.Type, depth int) bool {
	switch v := t.(type) {
	case *types.Named:
		if v.Obj().IsAlias() {
			return typeImplementsEncoder(v.Underlying(), depth)
		}

		var hasLen bool
		var hasEncode bool
		for i := 0; i < v.NumMethods(); i++ {
			switch m := v.Method(i); m.Name() {
			case "Len":
				s := m.Type().(*types.Signature)
				if s.Params().Len() != 0 {
					return false
				}
				if s.Results().Len() != 1 {
					return false
				}
				if !IsBasic(s.Results().At(0).Type(), types.Int) {
					return false
				}

				hasLen = true

			case "Encode":
				s := m.Type().(*types.Signature)

				if s.Params().Len() != 1 {
					return false
				}
				if s.Results().Len() != 1 {
					return false
				}
				if !isBytes(s.Params().At(0).Type()) {
					return false
				}
				if !isBytes(s.Results().At(0).Type()) {
					return false
				}

				hasEncode = true
			}
		}

		return hasLen && hasEncode

	case *types.Pointer:
		return depth < 1 && typeImplementsEncoder(v.Elem(), depth+1)

	default:
		return false
	}
}

func typeImplementsDecoder(t types.Type, depth int) bool {
	switch v := t.(type) {
	case *types.Named:
		if v.Obj().IsAlias() {
			return typeImplementsEncoder(v.Underlying(), depth)
		}

		for i := 0; i < v.NumMethods(); i++ {
			m := v.Method(i)
			if m.Name() != "Decode" {
				continue
			}

			s := m.Type().(*types.Signature)

			// Same checks as for Encode, but we won't do a func
			// here because it is just a coincidence, not a rule.
			if s.Params().Len() != 1 {
				return false
			}
			if s.Results().Len() != 2 {
				return false
			}
			if !isBytes(s.Params().At(0).Type()) {
				return false
			}
			if !isBytes(s.Results().At(0).Type()) {
				return false
			}
			if !IsErrorType(s.Results().At(1).Type()) {
				return false
			}

			return true
		}

		return false

	case *types.Pointer:
		return depth < 2 && typeImplementsDecoder(v.Elem(), depth+1)

	default:
		return false
	}
}

func isBytes(t types.Type) bool {
	v, ok := t.(*types.Slice)
	if !ok {
		return false
	}

	return IsBasic(v.Elem(), types.Byte)
}
