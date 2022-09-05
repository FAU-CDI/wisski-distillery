package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

// WebComponent represents the 'web' layer belonging to a distillery
type WebComponent struct {
	dis *Distillery
}

// Web returns the WebComponent belonging to this distillery
func (dis *Distillery) Web() WebComponent {
	return WebComponent{dis: dis}
}

func (WebComponent) Name() string {
	return "web"
}

func (web WebComponent) Stack() stack.Installable {
	return web.dis.makeComponentStack(web, stack.Installable{
		EnvFileContext: map[string]string{
			"DEFAULT_HOST": web.dis.Config.DefaultDomain,
		},
	})
}

func (web WebComponent) Path() string {
	return web.Stack().Dir
}
