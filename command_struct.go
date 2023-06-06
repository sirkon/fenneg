package fenneg

import "github.com/sirkon/errors"

const errorCannotContinue errors.Const = "cannot continue"

// commandStruct handles structure use case.
type commandStruct struct {
	Source sourcePoint `arg:"" help:"Source structure <path>:<struct> to generate code for its encoding/decoding."`
}

// Run command.
func (c *commandStruct) Run(ctx *runContext) error {
	r, err := NewRunner(ctx.args.ErrorsPath.String(), ctx.handlers)
	if err != nil {
		return errors.Wrap(err, "setup codegen runner")
	}

	if err := r.Struct(c.Source.Path, c.Source.ID); err != nil {
		return errors.Wrap(err, "run codegen for the given struct")
	}

	return nil
}
