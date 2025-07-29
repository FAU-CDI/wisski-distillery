//spellchecker:words main
package main

//spellchecker:words context runtime debug ggman internal pkglib exit
import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/FAU-CDI/wisski-distillery/cmd"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words workdir
func main() {
	// build the parameters
	params, err := cli.ParamsFromEnv()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// and run the command
	cmd := cmd.NewCommand(ctx, params)
	if err := cmd.Execute(); err != nil {
		code, _ := exit.CodeFromError(err)
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
		code.Return()
	}
}
