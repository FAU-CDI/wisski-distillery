package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Stop struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Stop)(nil)
)

func (*Stop) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "stop",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*Stop) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	return nil, instance.Barrel().Stack().Down(ctx, out)
}
