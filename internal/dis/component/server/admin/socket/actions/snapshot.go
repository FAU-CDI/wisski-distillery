//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes exporter
import (
	"context"
	"fmt"
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
	if err := s.dependencies.Exporter.MakeExport(
		ctx,
		out,
		exporter.ExportTask{
			Dest:     "",
			Instance: instance,

			StagingOnly: false,
		},
	); err != nil {
		return nil, fmt.Errorf("failed to make export: %w", err)
	}
	return nil, nil
}
