package cmd

//spellchecker:words github wisski distillery internal ingredient locker goprogram exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewInstanceLockCommand() *cobra.Command {
	impl := new(instanceLock)

	cmd := &cobra.Command{
		Use:     "instance_lock SLUG",
		Short:   "locks or unlocks an instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Lock, "lock", false, "lock the provided instance")
	flags.BoolVar(&impl.Unlock, "unlock", false, "unlock the provided instance")

	return cmd
}

type instanceLock struct {
	Lock        bool
	Unlock      bool
	Positionals struct {
		Slug string
	}
}

func (l *instanceLock) ParseArgs(cmd *cobra.Command, args []string) error {
	l.Positionals.Slug = args[0]

	if l.Lock == l.Unlock {
		return exit.NewErrorWithCode("exactly one of `--lock` and `--unlock` must be provied", exit.ExitCommandArguments)
	}
	return nil
}

var (
	errLockNoInstance = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
	errLockFailed     = exit.NewErrorWithCode("failed to update instance lock", exit.ExitGeneric)
)

func (l *instanceLock) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w %q: %w", errLockNoInstance, l.Positionals.Slug, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), l.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w %q: %w", errLockNoInstance, l.Positionals.Slug, err)
	}

	if l.Unlock {
		if err := instance.Locker().TryUnlock(cmd.Context()); err != nil {
			return fmt.Errorf("%w: %w", errLockFailed, err)
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "unlocked")
		return nil
	}

	if err := instance.Locker().TryLock(cmd.Context()); err != nil {
		return fmt.Errorf("%w: %w", errLockFailed, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "locked")
	return nil
}
