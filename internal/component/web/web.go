package web

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

// Web implements the web component
type Web struct {
	component.ComponentBase
}

func (Web) Name() string {
	return "web"
}

//go:embed all:stack
//go:embed web.env
var resources embed.FS

func (web Web) Stack() component.Installable {
	HTTPS_METHOD := "nohttp"
	if web.Config.HTTPSEnabled() {
		HTTPS_METHOD = "redirect"
	}

	return web.MakeStack(component.Installable{
		Resources:   resources,
		ContextPath: "stack",
		EnvPath:     "web.env",

		EnvContext: map[string]string{
			"DEFAULT_HOST": web.Config.DefaultDomain,
			"HTTPS_METHOD": HTTPS_METHOD,
		},
	})
}
