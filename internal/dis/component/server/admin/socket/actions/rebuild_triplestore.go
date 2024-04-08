package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type RebuildTriplestore struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Snapshot)(nil)
)

func (wsa *RebuildTriplestore) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "rebuild_triplestore",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (wsa *RebuildTriplestore) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.TRB().RebuildTriplestore(ctx, out, false)
}
