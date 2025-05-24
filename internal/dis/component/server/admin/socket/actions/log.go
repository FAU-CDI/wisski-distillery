//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes pkglib errorsx stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/stream"
)

type InstanceLog struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Snapshot)(nil)
)

func (*InstanceLog) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "instance_log",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstanceLog) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (a any, e error) {
	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return nil, fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Attach(ctx, stream.NonInteractive(out), false); err != nil {
		return nil, fmt.Errorf("failed to attach to stack: %w", err)
	}
	return nil, nil
}
