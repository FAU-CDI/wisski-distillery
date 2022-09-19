package sql

import (
	"context"
	"embed"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type SQL struct {
	component.ComponentBase

	ServerURL string // upstream server url

	PollContext  context.Context // context to abort polling with
	PollInterval time.Duration   // duration to wait for during wait

	sqlNetwork lazy.Lazy[string]
}

func (SQL) Name() string {
	return "sql"
}

//go:embed all:sql
var resources embed.FS

func (ssh *SQL) Stack(env environment.Environment) component.StackWithResources {
	return ssh.ComponentBase.MakeStack(env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "sql",

		MakeDirsPerm: environment.DefaultDirPerm,
		MakeDirs: []string{
			"data",
		},
	})
}
