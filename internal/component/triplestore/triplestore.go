package triplestore

import (
	"context"
	"embed"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
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

func (ts Triplestore) Stack() component.StackWithResources {
	return ts.ComponentBase.MakeStack(component.StackWithResources{
		Resources:   resources,
		ContextPath: "stack",

		CopyContextFiles: []string{"graphdb.zip"}, // TODO: Move into constant?

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
		},
	})
}
