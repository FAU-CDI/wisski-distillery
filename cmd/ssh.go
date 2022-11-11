package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// SSH is the 'ssh' command
var SSH wisski_distillery.Command = ssh{}

type ssh struct {
	Bind string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:2223"`
}

func (s ssh) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "ssh",
		Description: "Starts the ssh server to allow clients to connect to this distillery",
	}
}

var errSSHListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to listen",
}

func (s ssh) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	server, err := dis.SSH().Server(dis.Context(), context.IOStream)
	if err != nil {
		return err
	}

	context.Printf("Listening on %s\n", s.Bind)

	// make a new listener
	listener, err := dis.Still.Environment.Listen("tcp", s.Bind)
	if err != nil {
		return errServerListen.Wrap(err)
	}

	go func() {
		<-dis.Context().Done()
		listener.Close()
	}()

	// and serve that listener
	err = server.Serve(listener)
	return errServerListen.Wrap(err)
}
