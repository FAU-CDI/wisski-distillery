package cmd

import (
	"net"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// SSH is the 'ssh' command
var SSH wisski_distillery.Command = ssh{}

type ssh struct {
	Bind           string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:2223"`
	PrivateKeyPath string `short:"p" long:"private-key-path" description:"path to store private host keys in" required:"1"`
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

var errSSHServer = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to listen server",
}

var errSSHListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to listen",
}

func (s ssh) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	server, err := dis.SSH().Server(context.Context, s.PrivateKeyPath, context.Stderr)
	if err != nil {
		return errSSHServer.WrapError(err)
	}

	context.Printf("Listening on %s\n", s.Bind)

	// make a new listener
	listener, err := net.Listen("tcp", s.Bind)
	if err != nil {
		return errSSHListen.WrapError(err)
	}

	go func() {
		<-context.Context.Done()
		listener.Close()
	}()

	// and serve that listener
	err = server.Serve(listener)
	return errServerListen.WrapError(err)
}
