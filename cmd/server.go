package cmd

import (
	"net/http"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/tkw1536/goprogram/exit"
)

// Server is the 'server' command
var Server wisski_distillery.Command = server{}

type server struct {
	Prefix string `short:"p" long:"prefix" description:"prefix to listen under"`
	Bind   string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:8888"`
}

func (s server) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "server",
		Description: "Starts a server with information about this distillery",
	}
}

var errServerListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to listen",
}

func (s server) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	handler, err := dis.Control().Server(dis.Context(), context.IOStream)
	if err != nil {
		return err
	}

	context.Printf("Listening on %s\n", s.Bind)

	// make a new listener
	listener, err := dis.Core.Environment.Listen("tcp", s.Bind)
	if err != nil {
		return errServerListen.Wrap(err)
	}

	go func() {
		<-dis.Context().Done()
		listener.Close()
	}()

	// and serve that listener
	err = http.Serve(listener, http.StripPrefix(s.Prefix, handler))
	if err == nil {
		return nil
	}
	return errServerListen.Wrap(err)
}
