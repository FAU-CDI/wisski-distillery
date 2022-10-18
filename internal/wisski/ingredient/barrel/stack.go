package barrel

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

//go:embed all:barrel barrel.env
var barrelResources embed.FS

// Barrel returns a stack representing the running WissKI Instance
func (barrel *Barrel) Stack() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: barrel.FilesystemBase,
			Env: barrel.Malt.Environment,
		},

		Resources:   barrelResources,
		ContextPath: filepath.Join("barrel"),
		EnvPath:     filepath.Join("barrel.env"),

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": barrel.Malt.Config.DockerNetworkName,

			"SLUG":          barrel.Slug,
			"VIRTUAL_HOST":  barrel.Domain(),
			"HTTPS_ENABLED": barrel.Malt.Config.HTTPSEnabledEnv(),

			"DATA_PATH":                   filepath.Join(barrel.FilesystemBase, "data"),
			"RUNTIME_DIR":                 barrel.Malt.Config.RuntimeDir(),
			"GLOBAL_AUTHORIZED_KEYS_FILE": barrel.Malt.Config.GlobalAuthorizedKeysFile,
		},

		MakeDirs: []string{"data", ".composer"},

		TouchFiles: []string{
			filepath.Join("data", "authorized_keys"),
		},
	}
}
