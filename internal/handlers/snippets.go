package handlers

import "strings"

func isVariadic(t Type) bool {
	return t.Len() <= 0
}

func isFixed(t Type) bool {
	return t.Len() > 0
}

func dotIsSep(head string, a ...string) string {
	var buf strings.Builder
	buf.WriteString(strings.ReplaceAll(head, ".", "_"))
	for _, p := range a {
		buf.WriteByte('_')
		buf.WriteString(strings.ReplaceAll(p, ".", "_"))
	}

	return buf.String()
}
