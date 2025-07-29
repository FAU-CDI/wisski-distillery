package cmd

//spellchecker:words encoding json github wisski distillery internal component provision models ingredient barrel manager logging goprogram exit
import (
	"encoding/json"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/manager"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewProvisionCommand() *cobra.Command {
	impl := new(pv)

	cmd := &cobra.Command{
		Use:     "provision SLUG",
		Short:   "creates a new instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.StringVar(&impl.PHPVersion, "php", "", "specific php version to use for instance. See 'provision --list-php-versions' for available versions.")
	flags.BoolVar(&impl.ListPHPVersions, "list-php-versions", false, "List available php versions")
	flags.BoolVar(&impl.IIPServer, "iip-server", false, "enable iip-server inside this instance")
	flags.BoolVar(&impl.PHPDevelopment, "php-devel", false, "Include php development configuration")
	flags.StringVar(&impl.Flavor, "flavor", "", "Use specific flavor. Use '--list-flavors' to list flavors.")
	flags.BoolVar(&impl.ListFlavors, "list-flavors", false, "List all known flavors")
	flags.StringVar(&impl.ContentSecurityPolicy, "content-security-policy", "", "Setup ContentSecurityPolicy")

	return cmd
}

type pv struct {
	PHPVersion            string
	ListPHPVersions       bool
	IIPServer             bool
	PHPDevelopment        bool
	Flavor                string
	ListFlavors           bool
	ContentSecurityPolicy string
	Positionals           struct {
		Slug string
	}
}

func (p *pv) ParseArgs(cmd *cobra.Command, args []string) error {
	p.Positionals.Slug = args[0]

	if !p.ListFlavors && !p.ListPHPVersions && p.Positionals.Slug == "" {
		return errProvisionMissingSlug
	}
	return nil
}

func (*pv) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "provision",
		Description: "creates a new instance",
	}
}

var errProvisionMissingSlug = exit.NewErrorWithCode("must provide a slug", exit.ExitCommandArguments)

// TODO: AfterParse to check instance!

func (p *pv) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get distillery: %w", err)
	}

	if p.ListFlavors {
		return p.listFlavors(cmd)
	}
	if p.ListPHPVersions {
		return p.listPHPVersions(cmd)
	}

	instance, err := dis.Provision().Provision(cmd.ErrOrStderr(), cmd.Context(), provision.Flags{
		Slug:   p.Positionals.Slug,
		Flavor: p.Flavor,
		System: models.System{
			PHP:                   p.PHPVersion,
			IIPServer:             p.IIPServer,
			PHPDevelopment:        p.PHPDevelopment,
			ContentSecurityPolicy: p.ContentSecurityPolicy,
		},
	})
	if err != nil {
		return fmt.Errorf("%q: %w: %w", p.Positionals.Slug, errProvisionGeneric, err)
	}

	// and we're done!
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Instance has been provisioned"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL:      %s\n", instance.URL().String())
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Username: %s\n", instance.DrupalUsername)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Password: %s\n", instance.DrupalPassword)

	return nil
}

func (*pv) listFlavors(cmd *cobra.Command) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manager.Profiles()); err != nil {
		return fmt.Errorf("failed to encode flavors: %w", err)
	}
	return nil
}

func (*pv) listPHPVersions(cmd *cobra.Command) error {
	for _, v := range models.KnownPHPVersions() {
		if v == models.DefaultPHPVersion {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s (default)\n", v); err != nil {
				return fmt.Errorf("failed to print message: %w", err)
			}
		} else {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), v); err != nil {
				return fmt.Errorf("failed to print message: %w", err)
			}
		}
	}
	return nil
}
