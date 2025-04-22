//spellchecker:words binder
package binder

//spellchecker:words embed path filepath github wisski distillery internal component pkglib yamlx gopkg yaml
import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"
)

type Binder struct {
	component.Base
}

var (
	_ component.Installable = (*Binder)(nil)
)

func (binder *Binder) Path() string {
	return filepath.Join(component.GetStill(binder).Config.Paths.Root, "core", "binder")
}

func (binder *Binder) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed docker-compose.yml
var composeTemplate embed.FS

func (binder *Binder) Stack() component.StackWithResources {
	config := component.GetStill(binder).Config

	return component.MakeStack(binder, component.StackWithResources{
		ContextPath: ".",
		Resources:   composeTemplate,

		ComposerYML: func(root *yaml.Node) (*yaml.Node, error) {
			ports := config.Listen.ComposePorts("8000")
			if err := yamlx.ReplaceWith(root, ports, "services", "binder", "ports"); err != nil {
				return nil, err
			}

			command := config.HTTP.TCPMuxCommand("0.0.0.0:8000", "http:80", "http:443", "ssh:2222")
			if err := yamlx.ReplaceWith(root, command, "services", "binder", "command"); err != nil {
				return nil, err
			}

			return root, nil
		},

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),
		},
	})
}
