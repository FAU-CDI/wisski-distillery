package cmd

import (
	"net/http"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/tkw1536/goprogram/exit"
)

// DisServer is the 'dis_server' command
var DisServer wisski_distillery.Command = disServer{}

type disServer struct {
	Prefix string `short:"p" long:"prefix" description:"prefix to listen under"`
	Bind   string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:8888"`
}

func (disServer) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_server",
		Description: "Starts a server with information about this distillery",
	}
}

var errServerListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to listen",
}

func (s disServer) Run(context wisski_distillery.Context) error {
	server := context.Environment.Server()

	context.Printf("Listening on %s\n", s.Bind)
	err := http.ListenAndServe(s.Bind, http.StripPrefix(s.Prefix, server))
	if err == nil {
		return nil
	}
	return errServerListen.Wrap(err)
}
