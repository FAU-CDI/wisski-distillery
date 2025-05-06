package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"
	"net"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// SSH is the 'ssh' command.
var SSH wisski_distillery.Command = ssh{}

type ssh struct {
	Bind           string `default:"127.0.0.1:2223"                         description:"address to listen on" long:"bind"  short:"b"`
	PrivateKeyPath string `description:"path to store private host keys in" long:"private-key-path"            required:"1" short:"p"`
}

func (s ssh) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "ssh",
		Description: "starts the ssh server to allow clients to connect to this distillery",
	}
}

var (
	errSSHServer = exit.NewErrorWithCode("unable to listen server", exit.ExitGeneric)
	errSSHListen = exit.NewErrorWithCode("unable to listen", exit.ExitGeneric)
)

func (s ssh) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	server, err := dis.SSH().Server(context.Context, s.PrivateKeyPath, context.Stderr)
	if err != nil {
		return fmt.Errorf("%w: %w", errSSHServer, err)
	}

	_, _ = context.Printf("Listening on %s\n", s.Bind)

	// make a new listener
	listener, err := net.Listen("tcp", s.Bind)
	if err != nil {
		return fmt.Errorf("%w: %w", errSSHListen, err)
	}

	go func() {
		<-context.Context.Done()
		_ = listener.Close() // it is either closed or it isn't
	}()

	// and serve that listener
	err = server.Serve(listener)
	if err != nil {
		return fmt.Errorf("%w: %w", errServerListen, err)
	}
	return nil
}
