package main

import (
	"go/types"

	"github.com/sirkon/errors"
	"github.com/sirkon/fenneg"
	"github.com/sirkon/message"
)

const examplePkg = "github.com/sirkon/fenneg/example/internal/example"

func main() {
	hnlrs, err := fenneg.NewTypesHandlers(
		// Add custom handler for the example.Index type
		fenneg.NewTypeHandler(
			func(ot types.Type) fenneg.TypeHandler {
				var t *types.Named

				if v, ok := ot.(*types.Pointer); ok {
					ot = v.Elem()
				}

				t, ok := ot.(*types.Named)
				if !ok {
					return nil
				}

				if t.Obj().Pkg().Path() != "github.com/sirkon/fenneg/example/internal/example" {
					return nil
				}

				if t.Obj().Name() != "Index" {
					return nil
				}

				return NewHandlerIndex(t)
			},
		),
	)
	if err != nil {
		message.Critical(errors.Wrap(err, "set up type handlers"))
	}

	r, err := fenneg.NewRunner("github.com/sirkon/errors", hnlrs)
	if err != nil {
		message.Fatal(errors.Wrap(err, "setup codegen runner"))
	}

	if err := r.OpLog().
		Source(examplePkg, "Source").
		Type(examplePkg, "TypeRecorder").
		LengthPrefix(true).
		Run(); err != nil {
		message.Critical(errors.Wrap(err, "process oplog"))
	}

	fenneg.Chill()

	if err := r.Struct(examplePkg, "Struct"); err != nil {
		message.Critical(errors.Wrap(err, "process struct"))
	}
}
