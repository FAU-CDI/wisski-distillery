package cmd

//spellchecker:words github wisski distillery internal component auth goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/tkw1536/goprogram/exit"
)

// DisUser is the 'dis_user' command.
var DisUser wisski_distillery.Command = disUser{}

type disUser struct {
	CreateUser bool `description:"create a new user" long:"create" short:"c"`
	DeleteUser bool `description:"delete a user"     long:"delete" short:"d"`

	MakeAdmin   bool `description:"add admin permission to user"      long:"add-admin"    short:"a"`
	RemoveAdmin bool `description:"remove admin permission from user" long:"remove-admin" short:"A"`

	InfoUser  bool `description:"show information about a user" long:"info" short:"i"`
	ListUsers bool `description:"list all users"                long:"list" short:"l"`

	SetPassword   bool `description:"interactively set a user password"             long:"set-password"   short:"s"`
	UnsetPassword bool `description:"delete a users password and block the account" long:"unset-password" short:"u"`
	CheckPassword bool `description:"interactively check a user credential"         long:"check-password" short:"p"`

	EnableTOTP  bool `description:"interactively enroll a user in totp" long:"enable-totp"  short:"t"`
	DisableTOTP bool `description:"disable totp for a user"             long:"disable-totp" short:"v"`

	Positionals struct {
		User string `description:"username to manage. may be omitted for some actions" positional-arg-name:"USER"`
	} `positional-args:"true"`
}

func (disUser) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_user",
		Description: "manage distillery users",
	}
}

var errUserRequired = exit.NewErrorWithCode("`USER` argument is required", exit.ExitCommandArguments)

func (du disUser) AfterParse() error {
	var counter int
	for _, action := range []bool{
		du.CreateUser,
		du.InfoUser,
		du.DeleteUser,
		du.SetPassword,
		du.UnsetPassword,
		du.CheckPassword,
		du.ListUsers,
		du.DisableTOTP,
		du.EnableTOTP,
		du.MakeAdmin,
		du.RemoveAdmin,
	} {
		if action {
			counter++
		}
	}

	if counter != 1 {
		return errNoActionSelected
	}

	if !du.ListUsers && du.Positionals.User == "" {
		return errUserRequired
	}

	return nil
}

var errDisUserActionFailed = exit.NewErrorWithCode("action failed", exit.ExitGeneric)

func (du disUser) Run(context wisski_distillery.Context) (err error) {
	var userAction func(wisski_distillery.Context, *auth.AuthUser) error
	var genericAction func(wisski_distillery.Context) error

	switch {
	case du.ListUsers:
		genericAction = du.runListUsers
	case du.CreateUser:
		genericAction = du.runCreate

	case du.InfoUser:
		userAction = du.runInfo
	case du.DeleteUser:
		userAction = du.runDelete
	case du.SetPassword:
		userAction = du.runSetPassword
	case du.UnsetPassword:
		userAction = du.runUnsetPassword
	case du.CheckPassword:
		userAction = du.runCheckPassword

	case du.EnableTOTP:
		userAction = du.runEnableTOTP
	case du.DisableTOTP:
		userAction = du.runDisableTOTP
	case du.MakeAdmin:
		userAction = du.runMakeAdmin
	case du.RemoveAdmin:
		userAction = du.runRemoveAdmin
	}

	switch {
	case genericAction != nil:
		if err := genericAction(context); err != nil {
			return fmt.Errorf("%w: %w", errDisUserActionFailed, err)
		}
		return nil

	case userAction != nil:
		user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
		if err != nil {
			return fmt.Errorf("%w: failed to get user: %w", errDisUserActionFailed, err)
		}

		if err := userAction(context, user); err != nil {
			return fmt.Errorf("%w: %w", errDisUserActionFailed, err)
		}
		return nil
	}

	panic("never reached")
}

func (du disUser) runInfo(context wisski_distillery.Context, user *auth.AuthUser) error {
	_, _ = context.Println(user)
	return nil
}

func (du disUser) runCreate(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().CreateUser(context.Context, du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	_, _ = context.Println(user)
	return nil
}

func (du disUser) runDelete(context wisski_distillery.Context, user *auth.AuthUser) error {
	if err := user.Delete(context.Context); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

var errPasswordPolicy = exit.NewErrorWithCode("password policy failed: %s", exit.ExitGeneric)

func (du disUser) runSetPassword(context wisski_distillery.Context, user *auth.AuthUser) error {
	var passwd string
	{
		if _, err := context.Printf("Enter new password for user %s:", du.Positionals.User); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		passwd1, err := context.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if _, err := context.Println(); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}

		if _, err := context.Printf("Enter the same password again:"); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		passwd, err = context.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		_, _ = context.Println()

		if passwd != passwd1 {
			return errPasswordsNotIdentical
		}
		if err := user.CheckPasswordPolicy(passwd); err != nil {
			return fmt.Errorf("%w: %w", errPasswordPolicy, err)
		}
	}

	if err := user.SetPassword(context.Context, []byte(passwd)); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}

func (du disUser) runUnsetPassword(context wisski_distillery.Context, user *auth.AuthUser) error {
	if err := user.UnsetPassword(context.Context); err != nil {
		return fmt.Errorf("failed to unset password: %w", err)
	}
	return nil
}

func (du disUser) runCheckPassword(context wisski_distillery.Context, user *auth.AuthUser) error {
	if _, err := context.Printf("Enter password for %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	candidate, err := context.ReadPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	_, _ = context.Println()

	var passcode string
	if user.IsTOTPEnabled() {
		passcode, err = context.ReadPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		_, _ = context.Println()
	}

	if err := user.CheckCredentials(context.Context, []byte(candidate), passcode); err != nil {
		return fmt.Errorf("failed to check credentials: %w", err)
	}
	return nil
}

func (du disUser) runListUsers(context wisski_distillery.Context) error {
	users, err := context.Environment.Auth().Users(context.Context)
	if err != nil {
		return fmt.Errorf("failed to list all users: %w", err)
	}
	for _, user := range users {
		_, _ = context.Println(user)
	}
	return nil
}

func (du disUser) runEnableTOTP(context wisski_distillery.Context, user *auth.AuthUser) error {
	// get the secret
	key, err := user.NewTOTP(context.Context)
	if err != nil {
		return fmt.Errorf("failed to generate new totp: %w", err)
	}

	// print out the link
	url, err := auth.TOTPLink(key, 100, 100)
	if err != nil {
		return fmt.Errorf("failed to generate totp link: %w", err)
	}
	if _, err := context.Println(url); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	// request the passcode
	if _, err := context.Printf("Enter passcode for %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	passcode, err := context.ReadPassword()
	if err != nil {
		return fmt.Errorf("failed to read passcode: %w", err)
	}
	if _, err := context.Println(); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	// and enter it
	if err := user.EnableTOTP(context.Context, passcode); err != nil {
		return fmt.Errorf("failed to emable totp: %w", err)
	}
	return nil
}

func (du disUser) runDisableTOTP(context wisski_distillery.Context, user *auth.AuthUser) error {
	if err := user.DisableTOTP(context.Context); err != nil {
		return fmt.Errorf("failed to disable totp: %w", err)
	}
	return nil
}

func (du disUser) runMakeAdmin(context wisski_distillery.Context, user *auth.AuthUser) error {
	if err := user.MakeAdmin(context.Context); err != nil {
		return fmt.Errorf("failed to make admin: %w", err)
	}
	return nil
}

func (du disUser) runRemoveAdmin(context wisski_distillery.Context, user *auth.AuthUser) error {
	if err := user.MakeRegular(context.Context); err != nil {
		return fmt.Errorf("failed to make regular user: %w", err)
	}
	return nil
}
