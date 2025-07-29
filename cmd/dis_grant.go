package cmd

//spellchecker:words github wisski distillery internal component instances models goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewDisGrantCommand() *cobra.Command {
	impl := new(disGrant)

	cmd := &cobra.Command{
		Use:     "dis_grant USER [SLUG] [DRUPAL_USER]",
		Short:   "grant distillery users access to specific WissKIs",
		Args:    cobra.RangeArgs(1, 3),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.AddAll, "add-all", false, "add grant to all WissKIs")
	flags.BoolVar(&impl.AddUser, "add", false, "add or update a user to a given wisski")
	flags.BoolVar(&impl.RemoveUser, "remove", false, "remove a user from a given wisski")
	flags.BoolVar(&impl.DrupalAdmin, "admin", false, "grant user the admin role")

	return cmd
}

type disGrant struct {
	AddAll      bool
	AddUser     bool
	RemoveUser  bool
	DrupalAdmin bool
	Positionals struct {
		User       string
		Slug       string
		DrupalUser string
	}
}

func (dg *disGrant) ParseArgs(cmd *cobra.Command, args []string) error {
	dg.Positionals.User = args[0]
	if len(args) >= 2 {
		dg.Positionals.Slug = args[1]
	}
	if len(args) >= 3 {
		dg.Positionals.DrupalUser = args[2]
	}

	// Validate arguments
	var counter int
	for _, action := range []bool{
		dg.AddUser,
		dg.RemoveUser,
		dg.AddAll,
	} {
		if action {
			counter++
		}
	}

	if counter != 1 {
		return errNoActionSelected
	}

	if !dg.AddAll && dg.Positionals.Slug == "" {
		return errNoSlugSelect
	}

	return nil
}

func (*disGrant) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_grant",
		Description: "grant distillery users access to specific WissKIs",
	}
}

var errNoSlugSelect = exit.NewErrorWithCode("slug not provided", exit.ExitCommandArguments)
var errNoActionSelected = exit.NewErrorWithCode("no action selected", exit.ExitCommandArguments)
var errFailedGrant = exit.NewErrorWithCode("unable to manage grants", exit.ExitGeneric)

func (dg *disGrant) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errFailedGrant, err)
	}

	switch {
	case dg.AddUser:
		err = dg.runAddUser(cmd, dis)
	case dg.AddAll:
		err = dg.runAddAll(cmd, dis)
	case dg.RemoveUser:
		err = dg.runRemoveUser(cmd, dis)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errFailedGrant, err)
	}
	return nil
}

func (dg *disGrant) checkHasSlug(cmd *cobra.Command, dis *dis.Distillery) error {
	has, err := dis.Instances().Has(cmd.Context(), dg.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to check if instance exists: %w", err)
	}
	if !has {
		return instances.ErrWissKINotFound
	}
	return nil
}

func (dg *disGrant) runAddUser(cmd *cobra.Command, dis *dis.Distillery) error {
	if err := dg.checkHasSlug(cmd, dis); err != nil {
		return err
	}

	policy := dis.Policy()
	if err := policy.Set(cmd.Context(), models.Grant{
		User:            dg.Positionals.User,
		Slug:            dg.Positionals.Slug,
		DrupalUsername:  dg.Positionals.DrupalUser,
		DrupalAdminRole: dg.DrupalAdmin,
	}); err != nil {
		return fmt.Errorf("failed to set policy: %w", err)
	}
	return nil
}

func (dg *disGrant) runRemoveUser(cmd *cobra.Command, dis *dis.Distillery) error {
	if err := dg.checkHasSlug(cmd, dis); err != nil {
		return err
	}

	policy := dis.Policy()
	if err := policy.Remove(cmd.Context(), dg.Positionals.User, dg.Positionals.Slug); err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}
	return nil
}

func (dg *disGrant) runAddAll(cmd *cobra.Command, dis *dis.Distillery) error {
	policy := dis.Policy()

	instances, err := dis.Instances().All(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	for _, instance := range instances {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Adding grant for user %s to %s\n", dg.Positionals.User, instance.Slug); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		if err := policy.Set(cmd.Context(), models.Grant{
			User:            dg.Positionals.User,
			Slug:            instance.Slug,
			DrupalUsername:  dg.Positionals.User,
			DrupalAdminRole: dg.DrupalAdmin,
		}); err != nil {
			return fmt.Errorf("failed to add grant for instance %q to user: %w", instance.Slug, err)
		}
	}

	return nil
}
