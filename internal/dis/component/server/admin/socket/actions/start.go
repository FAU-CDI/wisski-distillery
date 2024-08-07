package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Start struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Start)(nil)
)

func (*Start) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "start",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*Start) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	return nil, instance.Barrel().Stack().Up(ctx, out)
}
