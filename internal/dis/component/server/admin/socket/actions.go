package socket

import (
	"context"
	"io"
)

func (sockets *Sockets) Actions() ActionMap {
	actions := make(ActionMap, len(sockets.dependencies.Actions)+len(sockets.dependencies.IActions))

	// setup basic actions
	for _, a := range sockets.dependencies.Actions {
		a := a
		meta := a.Action()
		actions[meta.Name] = Action{
			NumParams:  meta.NumParams,
			Scope:      meta.Scope,
			ScopeParam: meta.ScopeParam,

			Handle: a.Act,
		}
	}

	// setup instance actions
	for _, a := range sockets.dependencies.IActions {
		a := a
		meta := a.Action()
		actions[meta.Name] = Action{
			NumParams:  meta.NumParams + 1,
			Scope:      meta.Scope,
			ScopeParam: meta.ScopeParam,

			Handle: func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
				instance, err := sockets.dependencies.Instances.WissKI(ctx, params[0])
				if err != nil {
					return err
				}
				return a.Act(ctx, instance, in, out, params[1:]...)
			},
		}
	}

	return actions
}
