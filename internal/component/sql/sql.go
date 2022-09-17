package sql

import (
	"context"
	"embed"
	"io/fs"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

type SQL struct {
	component.ComponentBase

	ServerURL string // upstream server url

	PollContext  context.Context // context to abort polling with
	PollInterval time.Duration   // duration to wait for during wait
}

func (SQL) Name() string {
	return "sql"
}

//go:embed all:sql
var resources embed.FS

func (ssh SQL) Stack() component.StackWithResources {
	return ssh.ComponentBase.MakeStack(component.StackWithResources{
		Resources:   resources,
		ContextPath: "sql",

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			"data",
		},
	})
}
