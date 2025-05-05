//spellchecker:words socket
package socket

//spellchecker:words context errors http github process over websocket proto wisski distillery internal component server admin socket actions wdlog pkglib contextx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/process_over_websocket/proto"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/errorsx"
)

func (sockets *Sockets) Actions(ctx context.Context) proto.Handler {
	logger := wdlog.Of(ctx)

	actions := make(map[string]*actionable, len(sockets.dependencies.Actions)+len(sockets.dependencies.IActions))
	// setup basic actions
	for _, a := range sockets.dependencies.Actions {
		action, exec := sockets.regularAction(a)
		if _, ok := actions[action.Name]; ok {
			logger.Warn(
				"duplicate websocket action",
				"name", action.Name,
				"type", "regular",
			)
			continue
		}
		actions[action.Name] = exec

		logger.Info(
			"registering websocket action",

			"name", action.Name,
			"type", "regular",
			"params", action.NumParams,
			"scope", string(action.Scope),
			"scopeParam", action.ScopeParam,
		)
	}

	// setup instance actions
	for _, a := range sockets.dependencies.IActions {
		action, exec := sockets.instanceAction(a)
		if _, ok := actions[action.Name]; ok {
			logger.Warn(
				"duplicate websocket action",
				"name", action.Name,
				"type", "instance",
			)
		}
		actions[action.Name] = exec

		logger.Info(
			"registering websocket action",

			"name", action.Name,
			"type", "instance",
			"params", action.NumParams,
			"scope", string(action.Scope),
			"scopeParam", action.ScopeParam,
		)
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
				return errorsx.Combine(err, proto.ErrHandlerAuthorizationDenied)
			}

			if len(args) != meta.NumParams {
				return proto.ErrHandlerInvalidArgs
			}
			return nil
		},
		Run: func(ctx context.Context, input io.Reader, output io.Writer, args ...string) (res any, err error) {
			return a.Act(ctx, input, output, args...)
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
				return nil, fmt.Errorf("cannot find instance %q: %w", args[0], err)
			}

			return a.Act(ctx, instance, input, output, args[1:]...)
		},
	}
}

type actionable struct {
	Validate func(*http.Request, ...string) error
	Run      func(ctx context.Context, input io.Reader, output io.Writer, args ...string) (any, error)
}
