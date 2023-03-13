package triplestore

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

type Triplestore struct {
	component.Base

	BaseURL string // upstream server url

	PollInterval time.Duration // duration to wait for during wait
}

var (
	_ component.Backupable    = (*Triplestore)(nil)
	_ component.Snapshotable  = (*Triplestore)(nil)
	_ component.Installable   = (*Triplestore)(nil)
	_ component.Provisionable = (*Triplestore)(nil)
)

func (ts *Triplestore) Path() string {
	return filepath.Join(ts.Still.Config.Paths.Root, "core", "triplestore")
}

func (Triplestore) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:triplestore
//go:embed triplestore.env
var resources embed.FS

func (ts *Triplestore) Stack() component.StackWithResources {
	return component.MakeStack(ts, component.StackWithResources{
		Resources:   resources,
		ContextPath: "triplestore",

		CopyContextFiles: []string{"graphdb.zip"}, // TODO: Move into constant?

		EnvPath: "triplestore.env",
		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ts.Config.Docker.Network(),
		},

		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
			filepath.Join("data", "import"),
		},
	})
}
