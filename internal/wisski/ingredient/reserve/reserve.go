//spellchecker:words reserve
package reserve

//spellchecker:words embed github wisski distillery internal component ingredient dockerx
import (
	"embed"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

// Reserve implements reserving a WissKI Instance
// TODO: This should be integrated into the bookkeeping table.
type Reserve struct {
	ingredient.Base
}

//go:embed all:reserve
var reserveResources embed.FS

// Stack returns a stack representing the reserve instance.
func (reserve *Reserve) OpenStack() (component.StackWithResources, error) {
	liquid := ingredient.GetLiquid(reserve)
	config := ingredient.GetStill(reserve).Config

	stack, err := dockerx.NewStack(liquid.Docker, liquid.FilesystemBase)
	if err != nil {
		return component.StackWithResources{}, fmt.Errorf("failed to create stack: %w", err)
	}

	return component.StackWithResources{
		Stack: stack,

		Resources:   reserveResources,
		ContextPath: "reserve",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),

			"SLUG":          liquid.Slug,
			"HOST_RULE":     liquid.HostRule(),
			"HTTPS_ENABLED": config.HTTP.HTTPSEnabledEnv(),
		},
	}, nil
}
