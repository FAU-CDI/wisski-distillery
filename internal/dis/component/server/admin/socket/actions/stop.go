//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes pkglib errorsx
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/errorsx"
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

func (*Stop) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (a any, e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return nil, fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Down(ctx, out); err != nil {
		return nil, fmt.Errorf("failed to shutdown barrel: %w", err)
	}
	return nil, nil
}
