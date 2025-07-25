package cmd

//spellchecker:words github wisski distillery internal component instances models goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/goprogram/exit"
)

// DisGrant is the 'dis_grant' command.
var DisGrant wisski_distillery.Command = disGrant{}

type disGrant struct {
	AddAll     bool `description:"add grant to all WissKIs"               long:"add-all" short:"m"`
	AddUser    bool `description:"add or update a user to a given wisski" long:"add"     short:"a"`
	RemoveUser bool `description:"remove a user from a given wisski"      long:"remove"  short:"r"`

	DrupalAdmin bool `description:"grant user the admin role" long:"admin" short:"A"`

	Positionals struct {
		User       string `description:"distillery username" positional-arg-name:"USER"   required:"1-1"`
		Slug       string `description:"WissKI instance"     positional-arg-name:"SLUG"`
		DrupalUser string `description:"drupal username"     positional-arg-name:"DRUPAL"`
	} `positional-args:"true"`
}

func (disGrant) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_grant",
		Description: "grant distillery users access to specific WissKIs",
	}
}

var errNoSlugSelect = exit.NewErrorWithCode("slug not provided", exit.ExitCommandArguments)

func (dg disGrant) AfterParse() error {
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

var errFailedGrant = exit.NewErrorWithCode("unable to manage grants", exit.ExitGeneric)

func (dg disGrant) Run(context wisski_distillery.Context) (err error) {
	switch {
	case dg.AddUser:
		err = dg.runAddUser(context)
	case dg.AddAll:
		err = dg.runAddAll(context)
	case dg.RemoveUser:
		err = dg.runRemoveUser(context)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errFailedGrant, err)
	}
	return nil
}

func (dg disGrant) checkHasSlug(context wisski_distillery.Context) error {
	has, err := context.Environment.Instances().Has(context.Context, dg.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to check if instance exists: %w", err)
	}
	if !has {
		return instances.ErrWissKINotFound
	}
	return nil
}

func (dg disGrant) runAddUser(context wisski_distillery.Context) error {
	if err := dg.checkHasSlug(context); err != nil {
		return err
	}

	policy := context.Environment.Policy()
	if err := policy.Set(context.Context, models.Grant{
		User:            dg.Positionals.User,
		Slug:            dg.Positionals.Slug,
		DrupalUsername:  dg.Positionals.DrupalUser,
		DrupalAdminRole: dg.DrupalAdmin,
	}); err != nil {
		return fmt.Errorf("failed to set policy: %w", err)
	}
	return nil
}

func (dg disGrant) runRemoveUser(context wisski_distillery.Context) error {
	if err := dg.checkHasSlug(context); err != nil {
		return err
	}

	policy := context.Environment.Policy()
	if err := policy.Remove(context.Context, dg.Positionals.User, dg.Positionals.Slug); err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}
	return nil
}

func (dg disGrant) runAddAll(context wisski_distillery.Context) error {
	policy := context.Environment.Policy()

	instances, err := context.Environment.Instances().All(context.Context)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	for _, instance := range instances {
		if _, err := context.Printf("Adding grant for user %s to %s\n", dg.Positionals.User, instance.Slug); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		if err := policy.Set(context.Context, models.Grant{
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
