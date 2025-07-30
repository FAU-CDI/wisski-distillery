package cmd

//spellchecker:words github wisski distillery internal status wstatus cobra pkglib errorsx exit nobufio
import (
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	wstatus "github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/nobufio"
	"go.tkw01536.de/pkglib/status"
)

func NewDrupalUserCommand() *cobra.Command {
	impl := new(drupalUser)

	cmd := &cobra.Command{
		Use:     "drupal_user SLUG [USER]",
		Short:   "set a password for a specific user",
		Args:    cobra.RangeArgs(1, 2),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.CheckCommonPasswords, "check-common-passwords", false, "check for most common passwords. operates on all users concurrently.")
	flags.BoolVar(&impl.CheckPasswdInteractive, "check-password", false, "interactively check user password")
	flags.BoolVar(&impl.ResetPasswd, "reset-password", false, "reset password for user")
	flags.BoolVar(&impl.Login, "login", false, "print url to login as")

	return cmd
}

type drupalUser struct {
	CheckCommonPasswords   bool
	CheckPasswdInteractive bool
	ResetPasswd            bool
	Login                  bool
	Positionals            struct {
		Slug string
		User string
	}
}

func (du *drupalUser) ParseArgs(cmd *cobra.Command, args []string) error {
	du.Positionals.Slug = args[0]
	if len(args) >= 2 {
		du.Positionals.User = args[1]
	}

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
	errUserParameter = exit.NewErrorWithCode("incorrect username parameter", cli.ExitGeneric)

	errDrupalUserActionFailed = exit.NewErrorWithCode("action failed", cli.ExitGeneric)
)

func (du *drupalUser) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: failed to get WissKI: %w", errDrupalUserActionFailed, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), du.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: failed to get WissKI: %w", errDrupalUserActionFailed, err)
	}

	switch {
	case du.CheckCommonPasswords:
		err = du.checkCommonPassword(cmd, instance)
	case du.CheckPasswdInteractive:
		err = du.checkPasswordInteractive(cmd, instance)
	case du.ResetPasswd:
		err = du.resetPassword(cmd, instance)
	case du.Login:
		err = du.login(cmd, instance)
	default:
		panic("never reached")
	}

	if err != nil {
		return fmt.Errorf("%w: %w", errDrupalUserActionFailed, err)
	}
	return nil
}

func (du *drupalUser) login(cmd *cobra.Command, instance *wisski.WissKI) error {
	link, err := instance.Users().Login(cmd.Context(), nil, du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to login user: %w", err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), link)
	return nil
}

func (du *drupalUser) checkCommonPassword(cmd *cobra.Command, instance *wisski.WissKI) error {
	users := instance.Users()

	entities, err := users.All(cmd.Context(), nil)
	if err != nil {
		return fmt.Errorf("failed to list all users: %w", err)
	}

	if err := status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[wstatus.DrupalUser, error]{
		PrefixString: func(item wstatus.DrupalUser, index int) string {
			return fmt.Sprintf("User[%q]: ", item.Name)
		},
		PrefixAlign: true,
		Handler: func(user wstatus.DrupalUser, index int, writer io.Writer) (e error) {
			pv, err := users.GetPasswordValidator(cmd.Context(), string(user.Name))
			if err != nil {
				return fmt.Errorf("failed to get password validator: %w", err)
			}
			defer errorsx.Close(pv, &e, "password validator")

			return pv.CheckDictionary(cmd.Context(), writer)
		},
	}, entities); err != nil {
		return fmt.Errorf("failed to get check for common passwords: %w", err)
	}
	return nil
}

func (du *drupalUser) checkPasswordInteractive(cmd *cobra.Command, instance *wisski.WissKI) (e error) {
	validator, err := instance.Users().GetPasswordValidator(cmd.Context(), du.Positionals.User)
	if err != nil {
		return fmt.Errorf("failed to get password validator: %w", err)
	}
	defer errorsx.Close(validator, &e, "validator")

	for {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter a password to check:"); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}
		candidate, err := nobufio.ReadPassword(cmd.InOrStdin())
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("failed to write text: %w", err)
		}

		if candidate == "" {
			break
		}

		if validator.Check(cmd.Context(), candidate) {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "check passed")
		} else {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "check did not pass")
		}
	}

	return nil
}

func (du *drupalUser) resetPassword(cmd *cobra.Command, instance *wisski.WissKI) error {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter new password for user %s:", du.Positionals.User); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	passwd1, err := nobufio.ReadPassword(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter the same password again:"); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}
	passwd2, err := nobufio.ReadPassword(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	if passwd1 != passwd2 {
		return errPasswordsNotIdentical
	}

	if err := instance.Users().SetPassword(cmd.Context(), nil, du.Positionals.User, passwd1); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}
