package socket

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/FAU-CDI/process_over_websocket/proto"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/tkw1536/pkglib/contextx"

	"github.com/rs/zerolog"
)

func (sockets *Sockets) Actions(ctx context.Context) proto.Handler {
	logger := zerolog.Ctx(ctx)

	actions := make(map[string]*actionable, len(sockets.dependencies.Actions)+len(sockets.dependencies.IActions))
	// setup basic actions
	for _, a := range sockets.dependencies.Actions {
		action, exec := sockets.regularAction(a)
		if _, ok := actions[action.Name]; ok {
			logger.Warn().Str("name", action.Name).Str("type", "regular").Msg("duplicate websocket action")
			continue
		}
		actions[action.Name] = exec

		logger.Info().
			Str("name", action.Name).
			Str("type", "regular").
			Int("params", action.NumParams).
			Str("scope", string(action.Scope)).
			Str("scopeParam", action.ScopeParam).
			Msg("registering websocket action")
	}

	// setup instance actions
	for _, a := range sockets.dependencies.IActions {
		action, exec := sockets.instanceAction(a)
		if _, ok := actions[action.Name]; ok {
			logger.Warn().Str("name", action.Name).Str("type", "instance").Msg("duplicate websocket action")
		}
		actions[action.Name] = exec

		logger.Info().
			Str("name", action.Name).
			Str("type", "instance").
			Int("params", action.NumParams).
			Str("scope", string(action.Scope)).
			Str("scopeParam", action.ScopeParam).
			Msg("registering websocket action")
	}

	return proto.HandlerFunc(func(r *http.Request, name string, args ...string) (p proto.Process, err error) {
		action, ok := actions[name]
		if !ok {
			return nil, proto.ErrHandlerUnknownProcess
		}

		if err := action.Validate(r, args...); err != nil {
			return nil, err
		}

		return proto.ProcessFunc(func(ictx context.Context, input io.Reader, output io.Writer, args ...string) (res any, err error) {
			// defer func() {
			//	logger.Err(err).Str("action", name).Msg("finished pow action")
			// }()
			return action.Run(contextx.WithValuesOf(ictx, ctx), input, output, args...)
		}), nil
	})
}

func (sockets *Sockets) regularAction(a actions.WebsocketAction) (actions.Action, *actionable) {
	meta := a.Action()
	return meta, &actionable{
		Validate: func(r *http.Request, args ...string) error {
			if err := sockets.dependencies.Auth.CheckScope(meta.ScopeParam, meta.Scope, r); err != nil {
				return errors.Join(err, proto.ErrHandlerAuthorizationDenied)
			}

			if len(args) != meta.NumParams {
				return proto.ErrHandlerInvalidArgs
			}
			return nil
		},
		Run: func(ctx context.Context, input io.Reader, output io.Writer, args ...string) (res any, err error) {
			err = a.Act(ctx, input, output, args...)
			return err == nil, err
		},
	}
}

func (sockets *Sockets) instanceAction(a actions.WebsocketInstanceAction) (actions.InstanceAction, *actionable) {
	meta := a.Action()
	return meta, &actionable{
		Validate: func(r *http.Request, args ...string) error {
			if err := sockets.dependencies.Auth.CheckScope(meta.ScopeParam, meta.Scope, r); err != nil {
				return errors.Join(err, proto.ErrHandlerAuthorizationDenied)
			}

			if len(args) != meta.NumParams+1 {
				return proto.ErrHandlerInvalidArgs
			}
			return nil
		},
		Run: func(ctx context.Context, input io.Reader, output io.Writer, args ...string) (res any, err error) {
			instance, err := sockets.dependencies.Instances.WissKI(ctx, args[0])
			if err != nil {
				return nil, err
			}

			{
				err := a.Act(ctx, instance, input, output, args[1:]...)
				return err == nil, err
			}
		},
	}
}

type actionable struct {
	Validate func(*http.Request, ...string) error
	Run      func(ctx context.Context, input io.Reader, output io.Writer, args ...string) (any, error)
}
