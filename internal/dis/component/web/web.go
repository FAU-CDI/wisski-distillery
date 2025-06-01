package web

//spellchecker:words path filepath github wisski distillery internal component docker gopkg yaml embed
import (
	"fmt"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/docker"
	"gopkg.in/yaml.v3"

	_ "embed"
)

// Web implements the ingress gateway for the distillery.
//
// It consists of an nginx docker container and an optional letsencrypt container.
type Web struct {
	component.Base
	dependencies struct {
		Docker *docker.Docker
	}
}

var (
	_ component.Installable = (*Web)(nil)
)

func (web *Web) Path() string {
	return filepath.Join(component.GetStill(web).Config.Paths.Root, "core", "web")
}

func (*Web) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed docker-compose-http.yml
var dockerComposeHTTP []byte

//go:embed docker-compose-https.yml
var dockerComposeHTTPS []byte

func (web *Web) OpenStack() (component.StackWithResources, error) {
	var stack component.StackWithResources

	config := component.GetStill(web).Config
	stack.EnvContext = map[string]string{
		"DOCKER_NETWORK_NAME": config.Docker.Network(),
		"CERT_EMAIL":          config.HTTP.CertbotEmail,
	}

	if config.HTTP.HTTPSEnabled() {
		stack.ComposerYML = readYaml(dockerComposeHTTPS)
		stack.TouchFilesPerm = 0600
		stack.TouchFiles = []string{"acme.json"}
	} else {
		stack.ComposerYML = readYaml(dockerComposeHTTP)
	}

	//nolint:wrapcheck
	return component.OpenStack(web, web.dependencies.Docker, stack)
}

func readYaml(bytes []byte) func(*yaml.Node) (*yaml.Node, error) {
	return func(n *yaml.Node) (*yaml.Node, error) {
		var node yaml.Node
		err := yaml.Unmarshal(bytes, &node)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
		}
		return &node, nil
	}
}
