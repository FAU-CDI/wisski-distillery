package cmd

//spellchecker:words github wisski distillery internal models goprogram exit pkglib status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/status"
)

// Cron is the 'cron' command.
var Rebuild wisski_distillery.Command = rebuild{}

type rebuild struct {
	Parallel int `default:"1" description:"run on (at most) this many instances in parallel. 0 for no limit." long:"parallel" short:"a"`

	System                bool   `description:"Update the system configuration according to other flags"                                                         long:"system-update"           short:"s"`
	PHPVersion            string `description:"update to specific php version to use for instance. See 'provision --list-php-versions' for available versions. " long:"php"                     short:"p"`
	IIPServer             bool   `description:"enable iip-server inside this instance"                                                                           long:"iip-server"              short:"i"`
	PHPDevelopment        bool   `description:"Include php development configuration"                                                                            long:"php-devel"               short:"d"`
	Flavor                string `description:"Use specific flavor. Use 'provision --list-flavors' to list flavors. "                                            long:"flavor"                  short:"f"`
	ContentSecurityPolicy string `description:"Setup ContentSecurityPolicy"                                                                                      long:"content-security-policy" short:"c"`

	Positionals struct {
		Slug []string `description:"slug of instance or instances to run rebuild" positional-arg-name:"SLUG" required:"0"`
	} `positional-args:"true"`
}

var errRebuildNoSystem = exit.NewErrorWithCode("flags for system reconfiguration have been set, but `--system' was not provided", exit.ExitCommandArguments)

func (rb rebuild) AfterParse() error {
	if rb.System {
		return nil
	}
	if rb.PHPVersion != "" || rb.PHPDevelopment || rb.ContentSecurityPolicy != "" {
		return errRebuildNoSystem
	}
	return nil
}

func (rebuild) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "rebuild",
		Description: "runs the rebuild script for several instances",
	}
}

var errRebuildFailed = exit.NewErrorWithCode("failed to run rebuild", exit.ExitGeneric)

func (rb rebuild) Run(context wisski_distillery.Context) (err error) {
	dis := context.Environment

	// find the instances
	wissKIs, err := dis.Instances().Load(context.Context, rb.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: failed to get instances: %w", errRebuildFailed, err)
	}

	// and do the actual rebuild
	if err := status.WriterGroup(context.Stderr, rb.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		sys := instance.System
		if rb.System {
			sys = models.System{
				PHP:                   rb.PHPVersion,
				IIPServer:             rb.IIPServer,
				PHPDevelopment:        rb.PHPDevelopment,
				ContentSecurityPolicy: rb.ContentSecurityPolicy,
			}
		}

		return instance.SystemManager().Apply(context.Context, writer, sys)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("rebuild %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to rebuild systems: %w", errRebuildFailed, err)
	}
	return nil
}
