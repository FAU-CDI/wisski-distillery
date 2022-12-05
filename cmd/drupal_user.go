package cmd

import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	wstatus "github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/status"
)

// DrupalUser is the 'drupal_user' setting
var DrupalUser wisski_distillery.Command = drupalUser{}

type drupalUser struct {
	CheckCommonPasswords   bool `short:"d" long:"check-common-passwords" description:"check for most common passwords. operates on all users concurrently."`
	CheckPasswdInteractive bool `short:"c" long:"check-password" description:"interactively check user password"`
	ResetPasswd            bool `short:"r" long:"reset-password" description:"reset password for user"`
	Login                  bool `short:"l" long:"login" description:"print url to login as"`
	Positionals            struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to manage"`
		User string `positional-arg-name:"USER" description:"username to manage. may be omitted for some actions"`
	} `positional-args:"true"`
}

func (drupalUser) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "drupal_user",
		Description: "set a password for a specific user",
	}
}

var errNoActionSelected = exit.Error{
	Message:  "exactly one action must be selected",
	ExitCode: exit.ExitGeneric,
}

var errUserParameter = exit.Error{
	Message:  "incorrect username parameter",
	ExitCode: exit.ExitGeneric,
}

func (du drupalUser) AfterParse() error {
	var count int
	for _, s := range []bool{
		du.CheckCommonPasswords,
		du.CheckPasswdInteractive,
		du.ResetPasswd,
		du.Login,
	} {
		if s {
			count++
		}
	}
	if count != 1 {
		return errNoActionSelected
	}

	if du.CheckCommonPasswords != (du.Positionals.User == "") {
		return errUserParameter
	}

	return nil
}

var errPasswordsNotIdentical = exit.Error{
	Message:  "passwords are not identical",
	ExitCode: exit.ExitGeneric,
}

func (du drupalUser) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, du.Positionals.Slug)
	if err != nil {
		return err
	}

	switch {
	case du.CheckCommonPasswords:
		return du.checkCommonPassword(context, instance)
	case du.CheckPasswdInteractive:
		return du.checkPasswordInteractive(context, instance)
	case du.ResetPasswd:
		return du.resetPassword(context, instance)
	case du.Login:
		return du.login(context, instance)
	}
	panic("never reached")
}

func (du drupalUser) login(context wisski_distillery.Context, instance *wisski.WissKI) error {
	link, err := instance.Users().Login(context.Context, nil, du.Positionals.User)
	if err != nil {
		return err
	}
	context.Println(link)
	return nil
}

var errPasswordFound = exit.Error{
	Message:  "user had a dictionary password",
	ExitCode: 5,
}

func (du drupalUser) checkCommonPassword(context wisski_distillery.Context, instance *wisski.WissKI) error {
	users := instance.Users()

	entities, err := users.All(context.Context, nil)
	if err != nil {
		return err
	}

	return status.RunErrorGroup(context.Stderr, status.Group[wstatus.User, error]{
		PrefixString: func(item wstatus.User, index int) string {
			return fmt.Sprintf("User[%q]: ", item.Name)
		},
		PrefixAlign: true,
		Handler: func(user wstatus.User, index int, writer io.Writer) error {
			pv, err := users.GetPasswordValidator(context.Context, string(user.Name))
			if err != nil {
				return err
			}
			defer pv.Close()

			return pv.CheckDictionary(context.Context, writer)
		},
	}, entities)
}

func (du drupalUser) checkPasswordInteractive(context wisski_distillery.Context, instance *wisski.WissKI) error {
	validator, err := instance.Users().GetPasswordValidator(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}
	defer validator.Close()

	for {
		context.Printf("Enter a password to check:")
		candidate, err := context.IOStream.ReadPassword()
		if err != nil {
			return err
		}
		context.Println()

		if candidate == "" {
			break
		}

		if validator.Check(context.Context, candidate) {
			context.Println("check passed")
		} else {
			context.Println("check did not pass")
		}
	}

	return nil
}

func (du drupalUser) resetPassword(context wisski_distillery.Context, instance *wisski.WissKI) error {
	context.Printf("Enter new password for user %s:", du.Positionals.User)
	passwd1, err := context.IOStream.ReadPassword()
	if err != nil {
		return err
	}
	context.Println()

	context.Printf("Enter the same password again:")
	passwd2, err := context.IOStream.ReadPassword()
	if err != nil {
		return err
	}
	context.Println()

	if passwd1 != passwd2 {
		return errPasswordsNotIdentical
	}

	return instance.Users().SetPassword(context.Context, nil, du.Positionals.User, passwd1)
}
