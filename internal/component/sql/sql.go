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

//go:embed all:stack
var resources embed.FS

func (ssh SQL) Stack() component.Installable {
	return ssh.ComponentBase.MakeStack(component.Installable{
		Resources:   resources,
		ContextPath: "stack",

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			"data",
		},
	})
}
