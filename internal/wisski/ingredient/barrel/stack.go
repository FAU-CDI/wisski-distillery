//spellchecker:words barrel
package barrel

//spellchecker:words embed path filepath github wisski distillery internal component ingredient dockerx
import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

//go:embed all:barrel
var barrelResources embed.FS

const localSettingsName = "settings.local.php"

//go:embed local.settings.php
var localSettingsTemplate string

const phpIniName = "custom.ini"

//go:embed custom.ini
var phpIniTemplate string

// Barrel returns a stack representing the running WissKI Instance.
func (barrel *Barrel) OpenStack() (component.StackWithResources, error) {
	liquid := ingredient.GetLiquid(barrel)
	config := ingredient.GetStill(barrel).Config

	stack, err := dockerx.NewStack(liquid.Docker, liquid.FilesystemBase)
	if err != nil {
		return component.StackWithResources{}, fmt.Errorf("failed to get docker client: %w", err)
	}

	return component.StackWithResources{
		Stack: stack,

		Resources:   barrelResources,
		ContextPath: "barrel",

		CreateFiles: map[string]string{
			localSettingsName: localSettingsTemplate,
			phpIniName:        phpIniTemplate,
		},

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),

			"SLUG":            liquid.Slug,
			"HOST_RULE":       liquid.HostRule(),
			"WISSKI_HOSTNAME": liquid.Hostname(),
			"HTTPS_ENABLED":   config.HTTP.HTTPSEnabledEnv(),

			"DATA_PATH":   filepath.Join(liquid.FilesystemBase, "data"),
			"RUNTIME_DIR": config.Paths.RuntimeDir(),

			"LOCAL_SETTINGS_PATH":  filepath.Join(liquid.FilesystemBase, localSettingsName),
			"LOCAL_SETTINGS_MOUNT": LocalSettingsPath,

			"PHP_INI_PATH":  filepath.Join(liquid.FilesystemBase, phpIniName),
			"PHP_INI_MOUNT": PHPIniPath,

			"BARREL_BASE_IMAGE":       liquid.GetDockerBaseImage(),
			"IIP_SERVER_ENABLED":      liquid.GetIIPServerEnabled(),
			"PHP_CONFIG_MODE":         liquid.PHPDevelopmentMode(),
			"CONTENT_SECURITY_POLICY": liquid.ContentSecurityPolicy,
		},

		MakeDirs: []string{"data", ".composer"},
	}, nil
}
