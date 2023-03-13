package web

import (
	"bytes"
	"embed"
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// Web implements the ingress gateway for the distillery.
//
// It consists of an nginx docker container and an optional letsencrypt container.
type Web struct {
	component.Base
}

var (
	_ component.Installable = (*Web)(nil)
)

func (web *Web) Path() string {
	return filepath.Join(web.Still.Config.Paths.Root, "core", "web")
}

func (*Web) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed web.env
var webEnv embed.FS

//go:embed docker-compose-http.yml
var dockerComposeHTTP []byte

//go:embed docker-compose-https.yml
var dockerComposeHTTPS []byte

func (web *Web) Stack() component.StackWithResources {
	var stack component.StackWithResources
	stack.Resources = webEnv
	stack.EnvPath = "web.env"

	stack.EnvContext = map[string]string{
		"DOCKER_NETWORK_NAME": web.Config.Docker.Network(),
		"CERT_EMAIL":          web.Config.HTTP.CertbotEmail,
	}

	if web.Config.HTTP.HTTPSEnabled() {
		stack.ReadComposeFile = func() (io.Reader, error) {
			return bytes.NewReader(dockerComposeHTTPS), nil
		}
		stack.TouchFilesPerm = 0600
		stack.TouchFiles = []string{"acme.json"}
	} else {
		stack.ReadComposeFile = func() (io.Reader, error) {
			return bytes.NewReader(dockerComposeHTTP), nil
		}
	}

	return component.MakeStack(web, stack)
}
