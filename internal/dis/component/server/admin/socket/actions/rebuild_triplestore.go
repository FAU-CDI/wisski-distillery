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

type RebuildTriplestore struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*RebuildTriplestore)(nil)
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

func (wsa *RebuildTriplestore) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	size, err := instance.TRB().RebuildTriplestore(ctx, out, false)
	if err != nil {
		return 0, fmt.Errorf("failed to rebuild triplestore: %w", err)
	}
	return size, nil
}
