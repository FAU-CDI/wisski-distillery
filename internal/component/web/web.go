package web

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Web implements the ingress gateway for the distillery.
//
// It consists of an nginx docker container and an optional letsencrypt container.
type Web struct {
	component.ComponentBase
}

func (web *Web) Path() string {
	return filepath.Join(web.Still.Config.DeployRoot, "core", "web")
}

func (*Web) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

func (web Web) Stack(env environment.Environment) component.StackWithResources {
	if web.Config.HTTPSEnabled() {
		return web.stackHTTPS(env)
	} else {
		return web.stackHTTP(env)
	}
}

//go:embed all:web-https
//go:embed web.env
var httpsResources embed.FS

func (web *Web) stackHTTPS(env environment.Environment) component.StackWithResources {
	return component.MakeStack(web, env, component.StackWithResources{
		Resources:   httpsResources,
		ContextPath: "web-https",
		EnvPath:     "web.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": web.Config.DockerNetworkName,
			"CERT_EMAIL":          web.Config.CertbotEmail,
		},
		TouchFilesPerm: 0600,
		TouchFiles:     []string{"acme.json"},
	})
}

//go:embed all:web-http
//go:embed web.env
var httpResources embed.FS

func (web *Web) stackHTTP(env environment.Environment) component.StackWithResources {
	return component.MakeStack(web, env, component.StackWithResources{
		Resources:   httpResources,
		ContextPath: "web-http",
		EnvPath:     "web.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": web.Config.DockerNetworkName,
			"CERT_EMAIL":          web.Config.CertbotEmail,
		},
	})
}
