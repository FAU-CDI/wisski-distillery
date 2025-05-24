//spellchecker:words instances
package instances

//spellchecker:words context embed errors github wisski distillery internal component unpack
import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
)

var errBootstrapFailedRuntime = errors.New("failed to update runtime")

// Runtime contains runtime resources to be installed into any instance
//
//go:embed all:runtime
var runtimeResources embed.FS

// Update installs or updates runtime components needed by this component.
func (instances *Instances) Update(ctx context.Context, progress io.Writer) error {
	err := unpack.InstallDir(component.GetStill(instances).Config.Paths.RuntimeDir(), "runtime", runtimeResources, func(dst, src string) {
		// no sensible way to report errors
		_, _ = fmt.Fprintf(progress, "[copy]  %s\n", dst)
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errBootstrapFailedRuntime, err)
	}
	return nil
}
