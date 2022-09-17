package web

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

// Web implements the ingress gateway for the distillery.
//
// It consists of an nginx docker container and an optional letsencrypt container.
type Web struct {
	component.ComponentBase
}

func (Web) Name() string {
	return "web"
}

func (web Web) Stack() component.Installable {
	if web.Config.HTTPSEnabled() {
		return web.stackHTTPS()
	} else {
		return web.stackHTTP()
	}
}

//go:embed all:web-https
//go:embed web-https.env
var httpsResources embed.FS

func (web Web) stackHTTPS() component.Installable {
	return web.MakeStack(component.Installable{
		Resources:   httpsResources,
		ContextPath: "web-https",
		EnvPath:     "web-https.env",

		EnvContext: map[string]string{
			"DEFAULT_HOST": web.Config.DefaultDomain,
		},
	})
}

//go:embed all:web-http
//go:embed web-http.env
var httpResources embed.FS

func (web Web) stackHTTP() component.Installable {
	return web.MakeStack(component.Installable{
		Resources:   httpResources,
		ContextPath: "web-http",
		EnvPath:     "web-http.env",

		EnvContext: map[string]string{
			"DEFAULT_HOST": web.Config.DefaultDomain,
		},
	})
}
