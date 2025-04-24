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
	if err := b.dependencies.Exporter.MakeExport(
		ctx,
		out,
		exporter.ExportTask{
			Dest:     "",
			Instance: nil,

			StagingOnly: false,
		},
	); err != nil {
		return nil, fmt.Errorf("failed to create export: %w", err)
	}
	return nil, nil
}
