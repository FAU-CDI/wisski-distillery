package solr

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

type Solr struct {
	component.Base

	BaseURL string // upstream solr url

	PollInterval time.Duration // duration to wait for during wait
}

var (
	_ component.Installable = (*Solr)(nil)
)

func (s *Solr) Path() string {
	return filepath.Join(component.GetStill(s).Config.Paths.Root, "core", "solr")
}

func (*Solr) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:solr
var resources embed.FS

func (solr *Solr) Stack() component.StackWithResources {
	return component.MakeStack(solr, component.StackWithResources{
		Resources:   resources,
		ContextPath: "solr",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": component.GetStill(solr).Config.Docker.Network(),
		},

		MakeDirs: []string{
			filepath.Join("data"),
		},
	})
}
