//spellchecker:words instances
package instances

//spellchecker:words context embed github wisski distillery internal component unpack goprogram exit
import (
	"context"
	"embed"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/tkw1536/goprogram/exit"
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
func (instances *Instances) Update(ctx context.Context, progress io.Writer) error {
	err := unpack.InstallDir(component.GetStill(instances).Config.Paths.RuntimeDir(), "runtime", runtimeResources, func(dst, src string) {
		fmt.Fprintf(progress, "[copy]  %s\n", dst)
	})
	if err != nil {
		return errBootstrapFailedRuntime.WrapError(err)
	}
	return nil
}
