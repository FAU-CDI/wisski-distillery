package cmd

//spellchecker:words slog http sync time github wisski distillery internal wdlog cobra pkglib errorsx exit
import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
)

func NewServerCommand() *cobra.Command {
	impl := new(server)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "starts a server with information about this distillery",
		Args:  cobra.NoArgs,
		RunE:  impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Trigger, "trigger", false, "instead of running on the existing server, simply trigger a cron run")
	flags.StringVar(&impl.Bind, "bind", "127.0.0.1:8888", "address to listen on")
	flags.StringVar(&impl.InternalBind, "internal-bind", "127.0.0.1:9999", "address to listen on for internal server")

	return cmd
}

type server struct {
	Trigger      bool
	Bind         string
	InternalBind string
}

var errServerTrigger = exit.NewErrorWithCode("failed to trigger", cli.ExitGeneric)
var errServerGeneric = exit.NewErrorWithCode("unable to instantiate server", cli.ExitGeneric)

func (s *server) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get distillery: %w", err)
	}

	// if the caller requested a trigger, just trigger the cron tasks
	if s.Trigger {
		if err := dis.Control().Trigger(cmd.Context()); err != nil {
			return fmt.Errorf("%w: %w", errServerTrigger, err)
		}
	}

	{
		// create a channel for notifications
		notify, cancel := dis.Cron().Listen(cmd.Context())
		defer cancel()

		// start the cron tasks
		done := dis.Cron().Start(cmd.Context(), notify)
		defer func() {
			<-done
		}()
	}

	// and start the server
	public, internal, err := dis.Control().Server(cmd.Context(), cmd.ErrOrStderr())
	if err != nil {
		return fmt.Errorf("%w: %w", errServerGeneric, err)
	}

	// start the public listener
	publicS := http.Server{
		Handler:           public,
		ReadHeaderTimeout: 10 * time.Second,
	}
	publicC := make(chan error)
	var shutdownPublic func()
	{
		log := wdlog.Of(cmd.Context())
		log.Info(
			"listening public server",
			"bind", s.Bind,
		)

		publicL, err := net.Listen("tcp", s.Bind)
		if err != nil {
			return fmt.Errorf("%w: %w", errServerListen, err)
		}

		// shutdown the public server when done!
		shutdownPublic = sync.OnceFunc(func() {
			if err := publicS.Shutdown(cmd.Context()); err != nil {
				log.Error("failed to shutdown public server", slog.Any("error", err))
			}
		})
		defer shutdownPublic()

		go func() {
			publicC <- publicS.Serve(publicL)
		}()
	}

	// start the internal listener
	internalS := http.Server{
		Handler:           internal,
		ReadHeaderTimeout: 10 * time.Second,
	}
	internalC := make(chan error)

	var shutdownInternal func()
	{
		log := wdlog.Of(cmd.Context())
		log.Info(
			"listening internal server",
			"bind", s.InternalBind,
		)
		internalL, err := net.Listen("tcp", s.InternalBind)
		if err != nil {
			return fmt.Errorf("%w: %w", errServerListen, err)
		}

		// shutdown the internal server when done!
		shutdownInternal = sync.OnceFunc(func() {
			if err := internalS.Shutdown(cmd.Context()); err != nil {
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
		<-cmd.Context().Done()

		log := wdlog.Of(cmd.Context())
		log.Info("shutting down server")

		shutdownPublic()
		shutdownInternal()
	}()

	err = errorsx.Combine(<-internalC, <-publicC, err)
	if err != nil {
		return fmt.Errorf("%w: %w", errServerListen, err)
	}
	return nil
}
