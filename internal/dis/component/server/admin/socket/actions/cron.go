package actions

import (
	"context"
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

func (c *Cron) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.Drush().Cron(ctx, out)
}
