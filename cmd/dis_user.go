package cmd

//spellchecker:words github wisski distillery internal component auth goprogram exit
import (
	"bufio"
	"fmt"
	"strings"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewDisUserCommand() *cobra.Command {
	impl := new(disUser)

	cmd := &cobra.Command{
		Use:     "dis_user",
		Short:   "manage distillery users",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.CreateUser, "create", false, "create a new user")
	flags.BoolVar(&impl.DeleteUser, "delete", false, "delete a user")
	flags.BoolVar(&impl.MakeAdmin, "add-admin", false, "add admin permission to user")
	flags.BoolVar(&impl.RemoveAdmin, "remove-admin", false, "remove admin permission from user")
	flags.BoolVar(&impl.InfoUser, "info", false, "show information about a user")
	flags.BoolVar(&impl.ListUsers, "list", false, "list all users")
	flags.BoolVar(&impl.SetPassword, "set-password", false, "interactively set a user password")
	flags.BoolVar(&impl.UnsetPassword, "unset-password", false, "delete a users password and block the account")
	flags.BoolVar(&impl.CheckPassword, "check-password", false, "interactively check a user credential")
	flags.BoolVar(&impl.EnableTOTP, "enable-totp", false, "interactively enroll a user in totp")
	flags.BoolVar(&impl.DisableTOTP, "disable-totp", false, "disable totp for a user")

	return cmd
}

type disUser struct {
	CreateUser    bool
	DeleteUser    bool
	MakeAdmin     bool
	RemoveAdmin   bool
	InfoUser      bool
	ListUsers     bool
	SetPassword   bool
	UnsetPassword bool
	CheckPassword bool
	EnableTOTP    bool
	DisableTOTP   bool
	Positionals   struct {
		User string
	}
}

func (du *disUser) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		du.Positionals.User = args[0]
	}

	// Validate arguments
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

func (*disUser) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_user",
		Description: "manage distillery users",
	}
}

var errUserRequired = exit.NewErrorWithCode("`USER` argument is required", exit.ExitCommandArguments)
var errDisUserActionFailed = exit.NewErrorWithCode("action failed", exit.ExitGeneric)
var errPasswordsNotIdentical = exit.NewErrorWithCode("passwords not identical", exit.ExitGeneric)

type (
	_userAction    = func(*cobra.Command, *dis.Distillery, *auth.AuthUser) error
	_genericAction = func(*cobra.Command, *dis.Distillery) error
)

func (du *disUser) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errDisUserActionFailed, err)
	}

	var userAction _userAction
	var genericAction _genericAction

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
		if err := genericAction(cmd, dis); err != nil {
			return fmt.Errorf("%w: %w", errDisUserActionFailed, err)
		}
		return nil

	case userAction != nil:
		user, err := dis.Auth().User(cmd.Context(), du.Positionals.User)
		if err != nil {
			return fmt.Errorf("%w: failed to get user: %w", errDisUserActionFailed, err)
		}

		if err := userAction(cmd, dis, user); err != nil {
			return fmt.Errorf("%w: %w", errDisUserActionFailed, err)
		}
		return nil
	}

	panic("never reached")
}

func (du *disUser) runInfo(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), user)
	return nil
}

func (du *disUser) runCreate(cmd *cobra.Command, dis *dis.Distillery) error {
	user, err := dis.Auth().CreateUser(cmd.Context(), du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), user)
	return nil
}

func (du *disUser) runDelete(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if err := user.Delete(cmd.Context()); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

var errPasswordPolicy = exit.NewErrorWithCode("password policy failed: %s", exit.ExitGeneric)

func (du *disUser) runSetPassword(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	var passwd string
	{
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter new password for user %s:", du.Positionals.User); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		reader := bufio.NewReader(cmd.InOrStdin())
		passwd1, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		passwd1 = strings.TrimSpace(passwd1)
		if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter the same password again:"); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		passwd, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		passwd = strings.TrimSpace(passwd)
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		if passwd != passwd1 {
			return errPasswordsNotIdentical
		}
		if err := user.CheckPasswordPolicy(passwd); err != nil {
			return fmt.Errorf("%w: %w", errPasswordPolicy, err)
		}
	}

	if err := user.SetPassword(cmd.Context(), []byte(passwd)); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}

func (du *disUser) runUnsetPassword(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if err := user.UnsetPassword(cmd.Context()); err != nil {
		return fmt.Errorf("failed to unset password: %w", err)
	}
	return nil
}

func (du *disUser) runCheckPassword(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter password for %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	candidate, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	candidate = strings.TrimSpace(candidate)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	var passcode string
	if user.IsTOTPEnabled() {
		passcode, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		passcode = strings.TrimSpace(passcode)
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}

	if err := user.CheckCredentials(cmd.Context(), []byte(candidate), passcode); err != nil {
		return fmt.Errorf("failed to check credentials: %w", err)
	}
	return nil
}

func (du *disUser) runListUsers(cmd *cobra.Command, dis *dis.Distillery) error {
	users, err := dis.Auth().Users(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to list all users: %w", err)
	}
	for _, user := range users {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), user)
	}
	return nil
}

func (du *disUser) runEnableTOTP(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	// get the secret
	key, err := user.NewTOTP(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to generate new totp: %w", err)
	}

	// print out the link
	url, err := auth.TOTPLink(key, 100, 100)
	if err != nil {
		return fmt.Errorf("failed to generate totp link: %w", err)
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), url); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	// request the passcode
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter passcode for %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	reader := bufio.NewReader(cmd.InOrStdin())
	passcode, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read passcode: %w", err)
	}
	passcode = strings.TrimSpace(passcode)
	if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	// and enter it
	if err := user.EnableTOTP(cmd.Context(), passcode); err != nil {
		return fmt.Errorf("failed to emable totp: %w", err)
	}
	return nil
}

func (du *disUser) runDisableTOTP(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if err := user.DisableTOTP(cmd.Context()); err != nil {
		return fmt.Errorf("failed to disable totp: %w", err)
	}
	return nil
}

func (du *disUser) runMakeAdmin(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if err := user.MakeAdmin(cmd.Context()); err != nil {
		return fmt.Errorf("failed to make admin: %w", err)
	}
	return nil
}

func (du *disUser) runRemoveAdmin(cmd *cobra.Command, dis *dis.Distillery, user *auth.AuthUser) error {
	if err := user.MakeRegular(cmd.Context()); err != nil {
		return fmt.Errorf("failed to make regular user: %w", err)
	}
	return nil
}
