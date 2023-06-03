package fenneg

import (
	"runtime/debug"

	"github.com/sirkon/fenneg/internal/app"
	"github.com/sirkon/message"
)

type commandVersion struct{}

// Run command.
func (commandVersion) Run(ctx *runContext) error {
	var version string
	info, ok := debug.ReadBuildInfo()
	if !ok {
		version = "(devel)"
	} else {
		version = info.Main.Version
	}

	message.Info(app.Name, "version", version)
	return nil
}
