package cmd

//spellchecker:words errors slog http sync github wisski distillery internal wdlog goprogram exit
import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
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
			return errServerTrigger.WrapError(err)
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
		return errServerGeneric.WrapError(err)
	}

	// start the public listener
	publicS := http.Server{Handler: public}
	publicC := make(chan error)
	var shutdownPublic func()
	{
		log := wdlog.Of(context.Context)
		log.Info(
			"listening public server",
			"bind", s.Bind,
		)

		publicL, err := net.Listen("tcp", s.Bind)
		if err != nil {
			return errServerListen.WrapError(err)
		}

		// shutdown the public server when done!
		shutdownPublic = sync.OnceFunc(func() {
			if err := publicS.Shutdown(context.Context); err != nil {
				log.Error("failed to shutdown public server", slog.Any("error", err))
			}
		})
		defer shutdownPublic()

		go func() {
			publicC <- publicS.Serve(publicL)
		}()
	}

	// start the internal listener
	internalS := http.Server{Handler: internal}
	internalC := make(chan error)

	var shutdownInternal func()
	{
		log := wdlog.Of(context.Context)
		log.Info(
			"listening internal server",
			"bind", s.InternalBind,
		)
		internalL, err := net.Listen("tcp", s.InternalBind)
		if err != nil {
			return errServerListen.WrapError(err)
		}

		// shutdown the internal server when done!
		shutdownInternal = sync.OnceFunc(func() {
			if err := internalS.Shutdown(context.Context); err != nil {
				log.Error("failed to shutdown internal server", slog.Any("error", err))
			}
		})
		defer shutdownInternal()

		go func() {
			internalC <- internalS.Serve(internalL)
		}()
	}

	// shutdown everything when the context closes
	go func() {
		<-context.Context.Done()

		log := wdlog.Of(context.Context)
		log.Info("shutting down server")

		shutdownPublic()
		shutdownInternal()
	}()

	return errServerListen.WrapError(errors.Join(<-internalC, <-publicC, err))
}
