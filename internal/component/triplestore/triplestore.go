package triplestore

import (
	"context"
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

type Triplestore struct {
	component.ComponentBase

	BaseURL string // upstream server url

	PollContext  context.Context // context to abort polling with
	PollInterval time.Duration   // duration to wait for during wait
}

func (Triplestore) Name() string {
	return "triplestore"
}

func (ts Triplestore) Path() string {
	return filepath.Join(ts.Core.Config.DeployRoot, "core", ts.Name())
}

func (Triplestore) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:triplestore
var resources embed.FS

func (ts *Triplestore) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(ts, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "triplestore",

		CopyContextFiles: []string{"graphdb.zip"}, // TODO: Move into constant?

		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
		},
	})
}
