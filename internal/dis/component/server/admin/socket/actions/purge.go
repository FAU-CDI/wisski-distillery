//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes instances purger
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Purge struct {
	component.Base
	dependencies struct {
		Purger *purger.Purger
	}
}

var (
	_ WebsocketInstanceAction = (*Purge)(nil)
)

func (*Purge) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "purge",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (p *Purge) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := p.dependencies.Purger.Purge(ctx, out, instance.Slug); err != nil {
		return nil, fmt.Errorf("failed to purge system: %w", err)
	}
	return nil, nil
}
