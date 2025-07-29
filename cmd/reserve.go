package cmd

//spellchecker:words github wisski distillery internal component models logging goprogram exit pkglib errorsx
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

func NewReserveCommand() *cobra.Command {
	impl := new(reserve)

	cmd := &cobra.Command{
		Use:     "reserve",
		Short:   "reserves a new instance",
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type reserve struct {
	Positionals struct {
		Slug string
	}
}

func (r *reserve) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		r.Positionals.Slug = args[0]
	}
	return nil
}

func (*reserve) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "reserve",
		Description: "reserves a new instance",
	}
}

// TODO: AfterParse to check instance!

var (
	errReserveAlreadyExists = exit.NewErrorWithCode("instance already exists", exit.ExitGeneric)
	errReserveGeneric       = exit.NewErrorWithCode("unable to provision instance", exit.ExitGeneric)
	errReserveStack         = exit.NewErrorWithCode("failed to open stack", exit.ExitGeneric)
	errProvisionGeneric     = exit.NewErrorWithCode("unable to provision instance", exit.ExitGeneric)
)

func (r *reserve) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errReserveGeneric, err)
	}

	if err := r.run(cmd, dis); err != nil {
		return fmt.Errorf("%w: %w", errReserveGeneric, err)
	}
	return nil
}

func (r *reserve) run(cmd *cobra.Command, dis *dis.Distillery) (e error) {
	slug := r.Positionals.Slug

	// check that it doesn't already exist
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Reserving new WissKI instance %s", slug); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if exists, err := dis.Instances().Has(cmd.Context(), slug); err != nil || exists {
		return fmt.Errorf("%q: %w: ", slug, errReserveAlreadyExists)
	}

	// make it in-memory
	instance, err := dis.Instances().Create(slug, models.System{})
	if err != nil {
		return fmt.Errorf("%w: %w", errProvisionGeneric, err)
	}

	// check that the base directory does not exist
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Checking that base directory %s does not exist", instance.FilesystemBase); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		exists, err := fsx.Exists(instance.FilesystemBase)
		if err != nil {
			return fmt.Errorf("%w: %w", errProvisionGeneric, err)
		}
		if exists {
			return fmt.Errorf("%q: %w", slug, errReserveAlreadyExists)
		}
	}

	// get the stack
	s, err := instance.Reserve().OpenStack()
	if err != nil {
		return fmt.Errorf("%w: %w", errReserveStack, err)
	}
	defer errorsx.Close(s, &e, "stack")

	{
		if err := logging.LogOperation(func() error {
			return s.Install(cmd.Context(), cmd.ErrOrStderr(), component.InstallationContext{})
		}, cmd.ErrOrStderr(), "Installing docker stack"); err != nil {
			return fmt.Errorf("failed to install docker stack: %w", err)
		}

		if err := logging.LogOperation(func() error {
			return s.Update(cmd.Context(), cmd.ErrOrStderr(), true)
		}, cmd.ErrOrStderr(), "Updating docker stack"); err != nil {
			return fmt.Errorf("failed to update docker stack: %w", err)
		}
	}

	// and we're done!
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Instance has been reserved"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL:      %s\n", instance.URL().String())

	return nil
}
