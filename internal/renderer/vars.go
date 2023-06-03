package renderer

import "sync"

var (
	// Path to a package with an API that is close to "github.com/sirkon/errors".
	// This package needs to implement New(f) and Wrap(f) functions and error
	// value they produce needs to have structured context methods like Int, Bool,
	// Str, etc.
	structuredErrorsPkgPath = "github.com/sirkon/errors"
)

// SetStructuredErrorsPkgPath set structured errors package path. Must be called
// before the first project initialization.
func SetStructuredErrorsPkgPath(p string) {
	once.Do(func() {
		structuredErrorsPkgPath = p
	})
}

var once sync.Once
