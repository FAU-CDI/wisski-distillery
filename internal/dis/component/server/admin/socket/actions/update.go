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

type Update struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Update)(nil)
)

func (*Update) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "update",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (u *Update) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Composer().Update(ctx, out); err != nil {
		return nil, fmt.Errorf("failed to update composer: %w", err)
	}
	return nil, nil
}
