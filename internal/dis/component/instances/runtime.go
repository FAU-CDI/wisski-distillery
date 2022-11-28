package instances

import (
	"context"
	"embed"

	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errBootstrapFailedRuntime = exit.Error{
	Message:  "failed to update runtime",
	ExitCode: exit.ExitGeneric,
}

// Runtime contains runtime resources to be installed into any instance
//
//go:embed all:runtime
var runtimeResources embed.FS

// Update installs or updates runtime components needed by this component.
func (instances *Instances) Update(ctx context.Context, stream stream.IOStream) error {
	err := unpack.InstallDir(instances.Still.Environment, instances.Config.RuntimeDir(), "runtime", runtimeResources, func(dst, src string) {
		stream.Printf("[copy]  %s\n", dst)
	})
	if err != nil {
		return errBootstrapFailedRuntime.Wrap(err)
	}
	return nil
}
