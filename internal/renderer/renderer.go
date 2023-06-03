package renderer

import "github.com/sirkon/gogh"

type (
	Project = gogh.Module[*Imports]
	Package = gogh.Package[*Imports]
	Go      = gogh.GoRenderer[*Imports]
)
