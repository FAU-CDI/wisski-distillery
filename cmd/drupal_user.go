package cmd

//spellchecker:words github wisski distillery internal status wstatus goprogram exit pkglib
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	wstatus "github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/status"
)

// DrupalUser is the 'drupal_user' setting.
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

var (
	errNoActionSelected = exit.NewErrorWithCode("exactly one action must be selected", exit.ExitGeneric)
	errUserParameter    = exit.NewErrorWithCode("incorrect username parameter", exit.ExitGeneric)
)

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

var (
	errPasswordsNotIdentical  = exit.NewErrorWithCode("passwords are not identical", exit.ExitGeneric)
	errDrupalUserActionFailed = exit.NewErrorWithCode("action failed", exit.ExitGeneric)
)

func (du drupalUser) Run(context wisski_distillery.Context) (err error) {
	instance, err := context.Environment.Instances().WissKI(context.Context, du.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: failed to get WissKI: %w", errDrupalUserActionFailed, err)
	}

	switch {
	case du.CheckCommonPasswords:
		err = du.checkCommonPassword(context, instance)
	case du.CheckPasswdInteractive:
		err = du.checkPasswordInteractive(context, instance)
	case du.ResetPasswd:
		err = du.resetPassword(context, instance)
	case du.Login:
		err = du.login(context, instance)
	default:
		panic("never reached")
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errDrupalUserActionFailed, err)
	}
	return nil
}

func (du drupalUser) login(context wisski_distillery.Context, instance *wisski.WissKI) error {
	link, err := instance.Users().Login(context.Context, nil, du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to login user: %w", err)
	}
	_, _ = context.Println(link)
	return nil
}

func (du drupalUser) checkCommonPassword(context wisski_distillery.Context, instance *wisski.WissKI) error {
	users := instance.Users()

	entities, err := users.All(context.Context, nil)
	if err != nil {
		return fmt.Errorf("failed to list all users: %w", err)
	}

	if err := status.RunErrorGroup(context.Stderr, status.Group[wstatus.DrupalUser, error]{
		PrefixString: func(item wstatus.DrupalUser, index int) string {
			return fmt.Sprintf("User[%q]: ", item.Name)
		},
		PrefixAlign: true,
		Handler: func(user wstatus.DrupalUser, index int, writer io.Writer) (e error) {
			pv, err := users.GetPasswordValidator(context.Context, string(user.Name))
			if err != nil {
				return fmt.Errorf("failed to get password validator: %w", err)
			}
			defer errorsx.Close(pv, &e, "password validator")

			return pv.CheckDictionary(context.Context, writer)
		},
	}, entities); err != nil {
		return fmt.Errorf("failed to get check for common passwords: %w", err)
	}
	return nil
}

func (du drupalUser) checkPasswordInteractive(context wisski_distillery.Context, instance *wisski.WissKI) (e error) {
	validator, err := instance.Users().GetPasswordValidator(context.Context, du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to get password validator: %w", err)
	}
	defer errorsx.Close(validator, &e, "validator")

	for {
		if _, err := context.Printf("Enter a password to check:"); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		candidate, err := context.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if _, err := context.Println(); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}

		if candidate == "" {
			break
		}

		if validator.Check(context.Context, candidate) {
			_, _ = context.Println("check passed")
		} else {
			_, _ = context.Println("check did not pass")
		}
	}

	return nil
}

func (du drupalUser) resetPassword(context wisski_distillery.Context, instance *wisski.WissKI) error {
	if _, err := context.Printf("Enter new password for user %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	passwd1, err := context.ReadPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	_, _ = context.Println()

	if _, err := context.Printf("Enter the same password again:"); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	passwd2, err := context.ReadPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	_, _ = context.Println()

	if passwd1 != passwd2 {
		return errPasswordsNotIdentical
	}

	if err := instance.Users().SetPassword(context.Context, nil, du.Positionals.User, passwd1); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}
