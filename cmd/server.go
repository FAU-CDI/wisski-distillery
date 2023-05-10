package cmd

import (
	"errors"
	"net"
	"net/http"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/rs/zerolog"
	"github.com/tkw1536/goprogram/exit"
)

// Server is the 'server' command
var Server wisski_distillery.Command = server{}

type server struct {
	Trigger      bool   `short:"t" long:"trigger" description:"instead of running on the existing server, simply trigger a cron run"`
	Bind         string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:8888"`
	InternalBind string `short:"i" long:"internal-bind" description:"address to listen on for internal server" default:"127.0.0.1:9999"`
}

func (s server) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "server",
		Description: "starts a server with information about this distillery",
	}
}

var errServerListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to listen",
}

var errServerTrigger = exit.Error{
	Message:  "failed to trigger",
	ExitCode: exit.ExitGeneric,
}

var errServerGeneric = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to instantiate server",
}

func (s server) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// if the caller requested a trigger, just trigger the cron tasks
	if s.Trigger {
		if err := dis.Control().Trigger(context.Context); err != nil {
			return errServerTrigger.Wrap(err)
		}
	}

	{
		// create a channel for notifications
		notify, cancel := dis.Cron().Listen(context.Context)
		defer cancel()

		// start the cron tasks
		done := dis.Cron().Start(context.Context, notify)
		defer func() {
			<-done
		}()
	}

	// and start the server
	public, internal, err := dis.Control().Server(context.Context, context.Stderr)
	if err != nil {
		return errServerGeneric.Wrap(err)
	}

	// start the public listener
	publicS := http.Server{Handler: public}
	publicC := make(chan error)
	{
		zerolog.Ctx(context.Context).Info().Str("bind", s.Bind).Msg("listening public server")
		publicL, err := net.Listen("tcp", s.Bind)
		if err != nil {
			return errServerListen.Wrap(err)
		}
		defer publicS.Shutdown(context.Context)
		go func() {
			publicC <- publicS.Serve(publicL)
		}()
	}

	// start the internal listener
	internalS := http.Server{Handler: internal}
	internalC := make(chan error)
	{
		zerolog.Ctx(context.Context).Info().Str("bind", s.InternalBind).Msg("listening internal server")
		internalL, err := net.Listen("tcp", s.InternalBind)
		if err != nil {
			return errServerListen.Wrap(err)
		}
		defer internalS.Shutdown(context.Context)
		go func() {
			internalC <- internalS.Serve(internalL)
		}()
	}

	go func() {
		<-context.Context.Done()

		zerolog.Ctx(context.Context).Info().Msg("shutting down server")
		publicS.Shutdown(context.Context)
		internalS.Shutdown(context.Context)
	}()

	return errServerListen.Wrap(errors.Join(<-internalC, <-publicC, err))
}
