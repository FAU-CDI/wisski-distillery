package cmd

//spellchecker:words github wisski distillery internal models cobra pkglib exit status
import (
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/status"
)

func NewRebuildCommand() *cobra.Command {
	impl := new(rebuild)

	cmd := &cobra.Command{
		Use:     "rebuild SLUG...",
		Short:   "runs the rebuild script for several instances",
		Args:    cobra.ArbitraryArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVar(&impl.Parallel, "parallel", 1, "run on (at most) this many instances in parallel. 0 for no limit.")
	flags.BoolVar(&impl.System, "system-update", false, "Update the system configuration according to other flags")
	flags.StringVar(&impl.PHPVersion, "php", "", "update to specific php version to use for instance. See 'provision --list-php-versions' for available versions.")
	flags.BoolVar(&impl.IIPServer, "iip-server", false, "enable iip-server inside this instance")
	flags.BoolVar(&impl.PHPDevelopment, "php-devel", false, "Include php development configuration")
	flags.StringVar(&impl.Flavor, "flavor", "", "Use specific flavor. Use 'provision --list-flavors' to list flavors.")
	flags.StringVar(&impl.ContentSecurityPolicy, "content-security-policy", "", "Setup ContentSecurityPolicy")

	return cmd
}

type rebuild struct {
	Parallel int

	System                bool
	PHPVersion            string
	IIPServer             bool
	PHPDevelopment        bool
	Flavor                string
	ContentSecurityPolicy string

	Positionals struct {
		Slug []string
	}
}

func (rb *rebuild) ParseArgs(cmd *cobra.Command, args []string) error {
	rb.Positionals.Slug = args

	if rb.System {
		return nil
	}
	if rb.PHPVersion != "" || rb.PHPDevelopment || rb.ContentSecurityPolicy != "" {
		return errRebuildNoSystem
	}
	return nil
}

var errRebuildNoSystem = exit.NewErrorWithCode("flags for system reconfiguration have been set, but `--system' was not provided", cli.ExitCommandArguments)
var errRebuildFailed = exit.NewErrorWithCode("failed to run rebuild", cli.ExitGeneric)

func (rb *rebuild) Exec(cmd *cobra.Command, args []string) (err error) {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: failed to get instances: %w", errRebuildFailed, err)
	}

	// find the instances
	wissKIs, err := dis.Instances().Load(cmd.Context(), rb.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: failed to get instances: %w", errRebuildFailed, err)
	}

	// and do the actual rebuild
	if err := status.WriterGroup(cmd.ErrOrStderr(), rb.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		sys := instance.System
		if rb.System {
			sys = models.System{
				PHP:                   rb.PHPVersion,
				IIPServer:             rb.IIPServer,
				PHPDevelopment:        rb.PHPDevelopment,
				ContentSecurityPolicy: rb.ContentSecurityPolicy,
			}
		}

		return instance.SystemManager().Apply(cmd.Context(), writer, sys)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("rebuild %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to rebuild systems: %w", errRebuildFailed, err)
	}
	return nil
}
