package php

//spellchecker:words context embed github wisski distillery internal phpx pkglib stream
import (
	"context"
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/tkw1536/pkglib/stream"
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
	_ = php.dependencies.Barrel.ShellScript(ctx, str, "drush", "php:eval", code)
	return nil
}
