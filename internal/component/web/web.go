package web

import (
	"embed"
	"fmt"
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

func (Web) Name() string {
	return "web"
}

func (web Web) Path() string {
	res := filepath.Join(web.Core.Config.DeployRoot, "core", web.Name())
	fmt.Println("debug====" + res)
	return res
}

func (Web) Context(parent component.InstallationContext) component.InstallationContext {
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
//go:embed web-https.env
var httpsResources embed.FS

func (web *Web) stackHTTPS(env environment.Environment) component.StackWithResources {
	return component.MakeStack(web, env, component.StackWithResources{
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

func (web *Web) stackHTTP(env environment.Environment) component.StackWithResources {
	return component.MakeStack(web, env, component.StackWithResources{
		Resources:   httpResources,
		ContextPath: "web-http",
		EnvPath:     "web-http.env",

		EnvContext: map[string]string{
			"DEFAULT_HOST": web.Config.DefaultDomain,
		},
	})
}
