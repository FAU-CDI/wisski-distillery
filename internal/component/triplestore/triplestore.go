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

//go:embed all:stack
var resources embed.FS

func (ts Triplestore) Stack(env environment.Environment) component.StackWithResources {
	return ts.ComponentBase.MakeStack(env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "stack",

		CopyContextFiles: []string{"graphdb.zip"}, // TODO: Move into constant?

		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
		},
	})
}
