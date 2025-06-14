package handlers

func isVariadic(t Type) bool {
	return t.Len() <= 0
}

func isFixed(t Type) bool {
	return t.Len() > 0
}
