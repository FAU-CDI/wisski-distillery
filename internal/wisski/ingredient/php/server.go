package php

import (
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/phpserver"
	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
)

type Server = phpserver.Server

// NewServer returns a new server that can execute code within this distillery.
// When err == nil, the caller must call server.Close().
//
// See [PHPServer].
func (php *PHP) NewServer() (*Server, error) {
	return phpserver.New(func(str stream.IOStream, script string) {
		php.Barrel.Shell(str, "-c", shellescape.QuoteCommand([]string{"drush", "php:eval", script}))
	})
}
