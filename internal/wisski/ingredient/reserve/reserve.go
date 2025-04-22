//spellchecker:words reserve
package reserve

//spellchecker:words embed path filepath github wisski distillery internal component ingredient
import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Reserve implements reserving a WissKI Instance
// TODO: This should be integrated into the bookkeeping table.
type Reserve struct {
	ingredient.Base
}

//go:embed all:reserve
var reserveResources embed.FS

// Stack returns a stack representing the reserve instance.
func (reserve *Reserve) Stack() component.StackWithResources {
	liquid := ingredient.GetLiquid(reserve)
	config := ingredient.GetStill(reserve).Config
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: liquid.FilesystemBase,
		},

		Resources:   reserveResources,
		ContextPath: filepath.Join("reserve"),

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),

			"SLUG":          liquid.Slug,
			"HOST_RULE":     liquid.HostRule(),
			"HTTPS_ENABLED": config.HTTP.HTTPSEnabledEnv(),
		},
	}
}
