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

type Prune struct {
	component.Base
	dependencies struct {
		exporter *exporter.Exporter
	}
}

var _ WebsocketAction = (*Prune)(nil)

func (*Prune) Action() Action {
	return Action{
		Name:      "prune",
		Scope:     scopes.ScopeUserAdmin,
		NumParams: 0,
	}
}

func (pa *Prune) Act(ctx context.Context, in io.Reader, out io.Writer, params ...string) (any, error) {
	err := pa.dependencies.exporter.PruneExports(ctx, out)
	if err != nil {
		return nil, fmt.Errorf("failed to prune exports: %w", err)
	}
	return nil, nil
}
