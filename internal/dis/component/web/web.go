package web

import (
	"embed"
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

func (web Web) Stack() component.StackWithResources {
	if web.Config.HTTP.HTTPSEnabled() {
		return web.stackHTTPS()
	} else {
		return web.stackHTTP()
	}
}

//go:embed all:web-https
//go:embed web.env
var httpsResources embed.FS

func (web *Web) stackHTTPS() component.StackWithResources {
	return component.MakeStack(web, component.StackWithResources{
		Resources:   httpsResources,
		ContextPath: "web-https",
		EnvPath:     "web.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": web.Config.Docker.Network,
			"CERT_EMAIL":          web.Config.HTTP.CertbotEmail,
		},
		TouchFilesPerm: 0600,
		TouchFiles:     []string{"acme.json"},
	})
}

//go:embed all:web-http
//go:embed web.env
var httpResources embed.FS

func (web *Web) stackHTTP() component.StackWithResources {
	return component.MakeStack(web, component.StackWithResources{
		Resources:   httpResources,
		ContextPath: "web-http",
		EnvPath:     "web.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": web.Config.Docker.Network,
			"CERT_EMAIL":          web.Config.HTTP.CertbotEmail,
		},
	})
}
