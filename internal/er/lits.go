package er

import "strconv"

type (
	L string
	Q string
)

func (l L) String() string {
	return string(l)
}

func (q Q) String() string {
	return strconv.Quote(string(q))
}
