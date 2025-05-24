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

func (*Start) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (a any, e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return nil, fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Up(ctx, out); err != nil {
		return nil, fmt.Errorf("failed to start barrel: %w", err)
	}
	return nil, nil
}
