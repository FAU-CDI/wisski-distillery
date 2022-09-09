package cmd

import (
	"net/http"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/tkw1536/goprogram/exit"
)

// ResolverServer is the 'resolver_server' command
var ResolverServer wisski_distillery.Command = server{
	Desc: wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "resolver_server",
		Description: "Starts a global resolver server",
	},
	Server: func(context wisski_distillery.Context) (http.Handler, error) {
		return context.Environment.Resolver().Server(context.IOStream)
	},
}

// DisServer is the 'dis_server' command
var DisServer wisski_distillery.Command = server{
	Desc: wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_server",
		Description: "Starts a server with information about this distillery",
	},
	Server: func(context wisski_distillery.Context) (http.Handler, error) {
		return context.Environment.Server(), nil
	},
}

type server struct {
	Prefix string `short:"p" long:"prefix" description:"prefix to listen under"`
	Bind   string `short:"b" long:"bind" description:"address to listen on" default:"127.0.0.1:8888"`

	Desc   wisski_distillery.Description
	Server func(context wisski_distillery.Context) (http.Handler, error)
}

func (s server) Description() wisski_distillery.Description {
	return s.Desc
}

var errServerListen = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to listen",
}

func (s server) Run(context wisski_distillery.Context) error {
	handler, err := s.Server(context)
	if err != nil {
		return err
	}

	context.Printf("Listening on %s\n", s.Bind)
	err = http.ListenAndServe(s.Bind, http.StripPrefix(s.Prefix, handler))
	if err == nil {
		return nil
	}
	return errServerListen.Wrap(err)
}
