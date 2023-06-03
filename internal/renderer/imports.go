package renderer

import "github.com/sirkon/gogh"

// NewImports для конструктора рендерера.
func NewImports(imp *gogh.Imports) *Imports {
	return &Imports{
		imp: imp,
	}
}

// Imports is an extension of gogh.Imports to provide some frequently
// used imports.
type Imports struct {
	imp *gogh.Imports
}

// Add to implement gogh.Importer
func (i *Imports) Add(pkgpath string) *gogh.ImportAliasControl {
	return i.imp.Add(pkgpath)
}

// Module to implement gogh.Importer
func (i *Imports) Module(relpath string) *gogh.ImportAliasControl {
	return i.imp.Module(relpath)
}

// Imports дto implement gogh.Importer
func (i *Imports) Imports() *gogh.Imports {
	return i.imp
}

// Errors structured errors package. It is github.com/sirkon/errors
// by default. The value can be overrided
func (i *Imports) Errors() *gogh.ImportAliasControl {
	return i.imp.Add(structuredErrorsPkgPath)
}

// Binary импорт пакета binary из стандартной библиотеки.
func (i *Imports) Binary() *gogh.ImportAliasControl {
	return i.imp.Add("encoding/binary")
}

// Unsafe импорт пакета unsafe из стандартной библиотеки.
func (i *Imports) Unsafe() *gogh.ImportAliasControl {
	return i.imp.Add("unsafe")
}

// Varsize импорт пакета github.com/sirkon/varsize.
func (i *Imports) Varsize() *gogh.ImportAliasControl {
	return i.imp.Add("github.com/sirkon/varsize")
}

var _ gogh.Importer = &Imports{}
