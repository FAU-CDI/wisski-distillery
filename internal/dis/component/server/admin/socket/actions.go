package socket

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/proto"
	"github.com/rs/zerolog"
)

func (sockets *Sockets) Actions(ctx context.Context) proto.ActionMap {
	logger := zerolog.Ctx(ctx)
	actions := make(proto.ActionMap, len(sockets.dependencies.Actions)+len(sockets.dependencies.IActions))

	// setup basic actions
	for _, a := range sockets.dependencies.Actions {
		name, action := sockets.regularAction(a)
		if _, ok := actions[name]; ok {
			logger.Warn().Str("name", name).Str("type", "regular").Msg("duplicate websocket action")
		}
		actions[name] = action

		logger.Info().
			Str("name", name).
			Str("type", "regular").
			Int("params", action.NumParams).
			Str("scope", string(action.Scope)).
			Str("scopeParam", action.ScopeParam).
			Msg("registering websocket action")
	}

	// setup instance actions
	for _, a := range sockets.dependencies.IActions {
		name, action := sockets.instanceAction(a)
		if _, ok := actions[name]; ok {
			zerolog.Ctx(ctx).Warn().Str("name", name).Str("type", "instance").Msg("duplicate websocket action")
		}
		actions[name] = action

		logger.Info().
			Str("name", name).
			Str("type", "instance").
			Int("params", action.NumParams-1).
			Str("scope", string(action.Scope)).
			Str("scopeParam", action.ScopeParam).
			Msg("registering websocket action")
	}

	return actions
}

func (sockets *Sockets) regularAction(a actions.WebsocketAction) (name string, action proto.Action) {
	meta := a.Action()

	return meta.Name, proto.Action{
		NumParams:  meta.NumParams + 1,
		Scope:      meta.Scope,
		ScopeParam: meta.ScopeParam,

		Handle: a.Act,
	}
}

func (sockets *Sockets) instanceAction(a actions.WebsocketInstanceAction) (name string, action proto.Action) {
	meta := a.Action()

	return meta.Name, proto.Action{
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
