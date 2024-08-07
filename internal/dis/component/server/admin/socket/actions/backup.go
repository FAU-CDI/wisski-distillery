package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
)

type Backup struct {
	component.Base
	dependencies struct {
		Exporter *exporter.Exporter
	}
}

var (
	_ WebsocketAction = (*Backup)(nil)
)

func (*Backup) Action() Action {
	return Action{
		Name:      "backup",
		Scope:     scopes.ScopeUserAdmin,
		NumParams: 0,
	}
}

func (b *Backup) Act(ctx context.Context, in io.Reader, out io.Writer, params ...string) (any, error) {
	return nil, b.dependencies.Exporter.MakeExport(
		ctx,
		out,
		exporter.ExportTask{
			Dest:     "",
			Instance: nil,

			StagingOnly: false,
		},
	)
}
