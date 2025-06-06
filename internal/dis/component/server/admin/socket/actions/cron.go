//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Cron struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Cron)(nil)
)

func (*Cron) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "cron",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (c *Cron) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Drush().Cron(ctx, out); err != nil {
		return nil, fmt.Errorf("failed to run cron: %w", err)
	}
	return nil, nil
}
