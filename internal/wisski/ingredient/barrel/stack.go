//spellchecker:words barrel
package barrel

//spellchecker:words embed path filepath github wisski distillery internal component ingredient dockerx
import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/yamlx"
	"gopkg.in/yaml.v3"
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

		ComposerYML: func(root *yaml.Node) (*yaml.Node, error) {
			labelsNode, err := yamlx.Find(root, "services", "barrel", "labels")
			if err != nil {
				return nil, fmt.Errorf("failed to find labels: %w", err)
			}
			var labels []string
			if err := labelsNode.Decode(&labels); err != nil {
				return nil, fmt.Errorf("failed to decode labels: %w", err)
			}

			// add the middleware labels
			middleswares := barrel.makeMidlewares()
			labels = append(labels, makeMiddlewareLabels(middleswares...)...)

			if err := yamlx.ReplaceWith(root, labels, "services", "barrel", "labels"); err != nil {
				return nil, fmt.Errorf("failed to replace labels: %w", err)
			}

			return root, nil
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

func (barrel *Barrel) makeMidlewares() []map[string]string {

	middleswares := []map[string]string{
		map[string]string{
			"headers.customresponseheaders.x-drupal-cache":         "",
			"headers.customresponseheaders.x-drupal-dynamic-cache": "",
			"headers.customresponseheaders.x-generator":            "",
			"headers.customresponseheaders.x-powered-by":           "",
			"headers.customresponseheaders.Server":                 "",
		},
	}

	liquid := ingredient.GetLiquid(barrel)
	if len(liquid.IPAllowlist) != 0 {
		middleswares = append(middleswares, map[string]string{
			"ipallowlist.sourcerange": liquid.IPAllowlist,
		})
	}

	return middleswares
}

func makeMiddlewareLabels(middleswares ...map[string]string) (labels []string) {
	var counter int
	var names []string
	for _, middleware := range middleswares {
		if len(middleware) == 0 {
			continue
		}
		counter++
		name := fmt.Sprintf("wisski_%d_${SLUG}", counter)
		names = append(names, name+"@docker")
		for key, value := range middleware {
			labels = append(labels, fmt.Sprintf("traefik.http.middlewares.%s.%s=%s", name, key, value))
		}
	}

	if len(names) > 0 {
		labels = append(labels, "traefik.http.routers.wisski_${SLUG}.middlewares="+strings.Join(names, ","))
	}
	return labels

}
