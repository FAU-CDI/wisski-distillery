package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Snapshot struct {
	component.Base
	dependencies struct {
		Exporter *exporter.Exporter
	}
}

var (
	_ WebsocketInstanceAction = (*Snapshot)(nil)
)

func (*Snapshot) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "snapshot",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (s *Snapshot) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	// TODO: return the path
	return nil, s.dependencies.Exporter.MakeExport(
		ctx,
		out,
		exporter.ExportTask{
			Dest:     "",
			Instance: instance,

			StagingOnly: false,
		},
	)
}
