package fenneg

import "github.com/sirkon/errors"

// commandOpLog handles operation log use case.
type commandOpLog struct {
	Source   sourcePoint `arg:"" help:"Source interface <path>:<interface> to generate code over it."`
	Type     sourcePoint `arg:"" help:"User type <path>:<type-name> to generate encoding methods for."`
	Dispatch identifier  `help:"Dispatching function name. Optional, <type-name>Dispatch will be used by default" short:"d"`

	Handler   sourcePoint `help:"Dispatcher argument type of a dispatching function. This type must implement source interface. Optional, source interface type will be used by default" short:"a"`
	LenPrefix bool        `help:"Write uvarint length prefix first." short:"l"`
}

// Run command.
func (c commandOpLog) Run(ctx *runContext) error {
	r, err := NewRunner(ctx.args.ErrorsPath.String(), ctx.handlers)
	if err != nil {
		return errors.Wrap(err, "setup codegen runner")
	}

	olr := r.OpLog().
		Source(c.Source.Path, c.Source.ID).
		Type(c.Type.Path, c.Type.ID)
	if c.Handler.IsValid() {
		olr = olr.DispatchHandler(c.Handler.Path, c.Handler.ID)
	}
	if c.Dispatch.String() != "" {
		olr = olr.DispatchFunc(c.Dispatch.String())
	}
	olr = olr.LengthPrefix(c.LenPrefix)

	if err := olr.Run(); err != nil {
		return errors.Wrap(err, "run oplog codegen")
	}

	return nil
}
