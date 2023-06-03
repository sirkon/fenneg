package main

import (
	"go/types"
	"os"

	"github.com/sirkon/errors"
	"github.com/sirkon/fenneg"
	"github.com/sirkon/message"
)

func main() {
	const appName = "test"

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

	// This is just for example of course, do not overwrite like this
	// in your code for a general purpose utility.
	// It is OK though to write like this for something really specific
	// what is not needed anywhere else.
	//
	// Here we take Source interface and render methods for TypeRecorder.
	os.Args = []string{
		appName,
		"op-log",
		"-l", // Need to generate uvarint length prefix first.
		"github.com/sirkon/fenneg/example/internal/example:Source",
		"github.com/sirkon/fenneg/example/internal/example:TypeRecorder",
	}

	if err := fenneg.Run(appName, hnlrs); err != nil {
		message.Critical(errors.Wrap(err, "generate oplog thing"))
	}

	os.Args = []string{
		appName,
		"struct",
		"github.com/sirkon/fenneg/example/internal/example:Struct",
	}
	if err := fenneg.Run(appName, hnlrs); err != nil {
		message.Critical(errors.Wrap(err, "generate struct thing"))
	}
}
