package fenneg

import (
	"github.com/sirkon/fenneg/internal/renderer"
	"github.com/sirkon/gogh"
)

type (
	// Project makes Project public.
	Project = gogh.Module[*renderer.Imports]

	// Package makes Package public.
	Package = gogh.Package[*renderer.Imports]

	// Go makes Go renderer public.
	Go = gogh.GoRenderer[*renderer.Imports]
)

// SetStructuredErrorsPkgPath sets custom errors package.
func SetStructuredErrorsPkgPath(p string) {
	renderer.SetStructuredErrorsPkgPath(p)
}
