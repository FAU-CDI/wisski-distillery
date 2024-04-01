package barrel

import (
	"embed"
	"path/filepath"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

//go:embed all:barrel
var barrelResources embed.FS

const localSettingsName = "settings.local.php"

//go:embed local.settings.php
var localSettingsTemplate string

// Barrel returns a stack representing the running WissKI Instance
func (barrel *Barrel) Stack() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: barrel.FilesystemBase,
		},

		Resources:   barrelResources,
		ContextPath: filepath.Join("barrel"),

		CreateFiles: map[string]string{
			localSettingsName: localSettingsTemplate,
		},

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": barrel.Malt.Config.Docker.Network(),

			"SLUG":            barrel.Slug,
			"HOST_RULE":       barrel.HostRule(),
			"WISSKI_HOSTNAME": barrel.Hostname(),
			"HTTPS_ENABLED":   barrel.Malt.Config.HTTP.HTTPSEnabledEnv(),

			"DATA_PATH":   filepath.Join(barrel.FilesystemBase, "data"),
			"RUNTIME_DIR": barrel.Malt.Config.Paths.RuntimeDir(),

			"LOCAL_SETTINGS_PATH":  filepath.Join(barrel.FilesystemBase, localSettingsName),
			"LOCAL_SETTINGS_MOUNT": LocalSettingsPath,

			"BARREL_BASE_IMAGE":       barrel.GetDockerBaseImage(),
			"IIP_SERVER_ENABLED":      barrel.GetIIPServerEnabled(),
			"OPCACHE_MODE":            barrel.OpCacheMode(),
			"CONTENT_SECURITY_POLICY": barrel.ContentSecurityPolicy,
		},

		MakeDirs: []string{"data", ".composer"},
	}
}
