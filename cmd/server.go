package cmd

import (
	"net/http"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/cancel"
	"github.com/rs/zerolog"
	"github.com/tkw1536/goprogram/exit"
)

// Server is the 'server' command
var Server wisski_distillery.Command = server{}

type server struct {
	Trigger bool   `short:"t" long:"trigger" description:"instead of running on the existing server, simply trigger a cron run"`
	Prefix  string `short:"p" long:"prefix" description:"prefix to listen under"`
	Bind    string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:8888"`
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

func (s server) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// if the caller requested a trigger, just trigger the cron tasks
	if s.Trigger {
		return dis.Control().Trigger(context.Context, context.Environment.Environment)
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
	handler, err := dis.Control().Server(context.Context, context.Stderr)
	if err != nil {
		return err
	}

	zerolog.Ctx(context.Context).Info().Str("bind", s.Bind).Msg("listening")

	// create a new listener
	listener, err := dis.Still.Environment.Listen("tcp", s.Bind)
	if err != nil {
		return errServerListen.Wrap(err)
	}

	go func() {
		<-context.Context.Done()
		listener.Close()
	}()

	server := http.Server{
		Handler: http.StripPrefix(s.Prefix, handler),
	}

	err, _ = cancel.WithContext(context.Context, func(start func()) error {
		start()
		return server.Serve(listener)
	}, func() {
		zerolog.Ctx(context.Context).Info().Msg("shutting down server")
		server.Shutdown(context.Context)
	})

	return errServerListen.Wrap(err)
}
