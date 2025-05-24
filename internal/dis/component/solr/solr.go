//spellchecker:words solr
package solr

//spellchecker:words embed path filepath time github wisski distillery internal component docker
import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/docker"
)

type Solr struct {
	component.Base

	BaseURL string // upstream solr url

	PollInterval time.Duration // duration to wait for during wait

	dependencies struct {
		Docker *docker.Docker
	}
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

func (solr *Solr) OpenStack() (component.StackWithResources, error) {
	return component.OpenStack(solr, solr.dependencies.Docker, component.StackWithResources{
		Resources:   resources,
		ContextPath: "solr",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": component.GetStill(solr).Config.Docker.Network(),
		},

		MakeDirs: []string{
			"data",
		},
	})
}
