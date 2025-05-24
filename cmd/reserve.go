package cmd

//spellchecker:words github wisski distillery internal component models logging goprogram exit pkglib errorsx
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/fsx"
)

// Reserve is the 'reserve' command.
var Reserve wisski_distillery.Command = reserve{}

type reserve struct {
	Positionals struct {
		Slug string `description:"name of instance to reserve" positional-arg-name:"slug" required:"1-1"`
	} `positional-args:"true"`
}

func (reserve) Description() wisski_distillery.Description {
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
)

func (r reserve) Run(context wisski_distillery.Context) (err error) {
	if err := r.run(context); err != nil {
		return fmt.Errorf("%w: %w", errReserveGeneric, err)
	}
	return nil
}

func (r reserve) run(context wisski_distillery.Context) (e error) {
	dis := context.Environment
	slug := r.Positionals.Slug

	// check that it doesn't already exist
	if _, err := logging.LogMessage(context.Stderr, "Reserving new WissKI instance %s", slug); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if exists, err := dis.Instances().Has(context.Context, slug); err != nil || exists {
		return fmt.Errorf("%q: %w: ", slug, errReserveAlreadyExists)
	}

	// make it in-memory
	instance, err := dis.Instances().Create(slug, models.System{})
	if err != nil {
		return fmt.Errorf("%w: %w", errProvisionGeneric, err)
	}

	// check that the base directory does not exist
	{
		if _, err := logging.LogMessage(context.Stderr, "Checking that base directory %s does not exist", instance.FilesystemBase); err != nil {
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
			return s.Install(context.Context, context.Stderr, component.InstallationContext{})
		}, context.Stderr, "Installing docker stack"); err != nil {
			return fmt.Errorf("failed to install docker stack: %w", err)
		}

		if err := logging.LogOperation(func() error {
			return s.Update(context.Context, context.Stderr, true)
		}, context.Stderr, "Updating docker stack"); err != nil {
			return fmt.Errorf("failed to update docker stack: %w", err)
		}
	}

	// and we're done!
	if _, err := logging.LogMessage(context.Stderr, "Instance has been reserved"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	_, _ = context.Printf("URL:      %s\n", instance.URL().String())

	return nil
}
