package sql

import (
	"context"
	"embed"
	"io/fs"
	"time"

	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
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

func (ssh SQL) Stack() stack.Installable {
	return ssh.ComponentBase.MakeStack(stack.Installable{
		Resources:   resources,
		ContextPath: "stack",

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			"data",
		},
	})
}
