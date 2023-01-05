package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// DisGrant is the 'dis_grant' command
var DisGrant wisski_distillery.Command = disGrant{}

type disGrant struct {
	AddUser    bool `short:"a" long:"add" description:"add or update a user to a given wisski"`
	RemoveUser bool `short:"r" long:"remove" description:"remove a user from a given wisski"`

	DrupalAdmin bool `short:"A" long:"admin" description:"grant user the admin role"`

	Positionals struct {
		User       string `positional-arg-name:"USER" required:"1-1" description:"distillery username"`
		Slug       string `positional-arg-name:"SLUG" required:"1-1" description:"WissKI instance"`
		DrupalUser string `positional-arg-name:"DRUPAL" description:"drupal username"`
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

func (dg disGrant) AfterParse() error {
	var counter int
	for _, action := range []bool{
		dg.AddUser,
		dg.RemoveUser,
	} {
		if action {
			counter++
		}
	}

	if counter != 1 {
		return errNoActionSelected
	}

	return nil
}

func (dg disGrant) Run(context wisski_distillery.Context) error {
	switch {
	case dg.AddUser:
		return dg.runAddUser(context)
	case dg.RemoveUser:
		return dg.runRemoveUser(context)
	}
	panic("never reached")
}

func (dg disGrant) checkHasSlug(context wisski_distillery.Context) error {
	has, err := context.Environment.Instances().Has(context.Context, dg.Positionals.Slug)
	if err != nil {
		return err
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
	return policy.Set(context.Context, models.Grant{
		User:            dg.Positionals.User,
		Slug:            dg.Positionals.Slug,
		DrupalUsername:  dg.Positionals.DrupalUser,
		DrupalAdminRole: dg.DrupalAdmin,
	})
}

func (dg disGrant) runRemoveUser(context wisski_distillery.Context) error {
	if err := dg.checkHasSlug(context); err != nil {
		return err
	}

	policy := context.Environment.Policy()
	return policy.Remove(context.Context, dg.Positionals.User, dg.Positionals.Slug)
}
