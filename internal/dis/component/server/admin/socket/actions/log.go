package actions

import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
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

func (*InstanceLog) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	err := instance.Barrel().Stack().Attach(ctx, stream.NonInteractive(out), false, "barrel")
	if err != nil {
		return nil, fmt.Errorf("failed to redirect log output: %w", err)
	}
	return nil, nil
}
