package php

import (
	"context"
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
)

// NewServer returns a new server that can execute code within this distillery.
// When err == nil, the caller must call server.Close().
//
// See [PHPServer].
func (php *PHP) NewServer() *phpx.Server {
	return &phpx.Server{
		Context:  context.Background(),
		Executor: phpx.SpawnFunc(php.spawn),
	}
}

func (php *PHP) spawn(ctx context.Context, str stream.IOStream, code string) error {
	php.Dependencies.Barrel.Shell(ctx, str, "-c", shellescape.QuoteCommand([]string{"drush", "php:eval", code}))()
	return nil
}
