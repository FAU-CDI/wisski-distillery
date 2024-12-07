package cmd

//spellchecker:words github wisski distillery internal models goprogram exit pkglib status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/status"
)

// Cron is the 'cron' command
var Rebuild wisski_distillery.Command = rebuild{}

type rebuild struct {
	Parallel int `short:"a" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`

	System                bool   `short:"s" long:"system-update" description:"Update the system configuration according to other flags"`
	PHPVersion            string `short:"p" long:"php" description:"update to specific php version to use for instance. See 'provision --list-php-versions' for available versions. "`
	IIPServer             bool   `short:"i" long:"iip-server" description:"enable iip-server inside this instance"`
	OPCacheDevelopment    bool   `short:"o" long:"opcache-devel" description:"Include opcache development configuration"`
	Flavor                string `short:"f" long:"flavor" description:"Use specific flavor. Use 'provision --list-flavors' to list flavors. "`
	ContentSecurityPolicy string `short:"c" long:"content-security-policy" description:"Setup ContentSecurityPolicy"`

	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance or instances to run rebuild"`
	} `positional-args:"true"`
}

var errRebuildNoSystem = exit.Error{
	Message:  "flags for system reconfiguration have been set, but `--system' was not provided",
	ExitCode: exit.ExitCommandArguments,
}

func (rb rebuild) AfterParse() error {
	if rb.System {
		return nil
	}
	if rb.PHPVersion != "" || rb.OPCacheDevelopment || rb.ContentSecurityPolicy != "" {
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

var errRebuildFailed = exit.Error{
	Message:  "failed to run rebuild",
	ExitCode: exit.ExitGeneric,
}

func (rb rebuild) Run(context wisski_distillery.Context) (err error) {
	defer errRebuildFailed.DeferWrap(&err)

	dis := context.Environment

	// find the instances
	wissKIs, err := dis.Instances().Load(context.Context, rb.Positionals.Slug...)
	if err != nil {
		return err
	}

	// and do the actual rebuild
	return status.WriterGroup(context.Stderr, rb.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		sys := instance.System
		if rb.System {
			sys = models.System{
				PHP:                   rb.PHPVersion,
				IIPServer:             rb.IIPServer,
				OpCacheDevelopment:    rb.OPCacheDevelopment,
				ContentSecurityPolicy: rb.ContentSecurityPolicy,
			}
		}

		return instance.SystemManager().Apply(context.Context, writer, sys)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("rebuild %q", item.Slug)
	}))
}
