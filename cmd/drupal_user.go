package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
)

// DrupalUser is the 'drupal_user' setting
var DrupalUser wisski_distillery.Command = duser{}

type duser struct {
	Passwd      bool `short:"p" long:"password" description:"reset password for user"`
	Login       bool `short:"l" long:"login" description:"print url to login as"`
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to manage"`
		User string `positional-arg-name:"USER" description:"username to manage"`
	} `positional-args:"true"`
}

func (duser) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "drupal_user",
		Description: "set a password for a specific user",
	}
}

var errPasswordsNotIdentical = exit.Error{
	Message:  "Passwords are not identical",
	ExitCode: exit.ExitGeneric,
}

func (du duser) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(du.Positionals.Slug)
	if err != nil {
		return err
	}

	if du.Passwd {
		return du.resetPassword(context, instance)
	}
	return du.login(context, instance)
}

func (du duser) login(context wisski_distillery.Context, instance *wisski.WissKI) error {
	link, err := instance.Drush().Login(context.IOStream, du.Positionals.User)
	if err != nil {
		return err
	}
	context.Println(link)
	return nil
}

func (du duser) resetPassword(context wisski_distillery.Context, instance *wisski.WissKI) error {
	context.Printf("Enter new password for user %s:", du.Positionals.User)
	passwd1, err := context.IOStream.ReadPassword()
	if err != nil {
		return err
	}

	context.Printf("Enter the same password again:")
	passwd2, err := context.IOStream.ReadPassword()
	if err != nil {
		return err
	}

	if passwd1 != passwd2 {
		return errPasswordsNotIdentical
	}

	return instance.Drush().ResetPassword(context.IOStream, du.Positionals.User, passwd1)
}
