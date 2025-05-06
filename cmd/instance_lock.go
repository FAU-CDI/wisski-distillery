package cmd

//spellchecker:words github wisski distillery internal ingredient locker goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/tkw1536/goprogram/exit"
)

// InstanceLock is then 'instance_lock' command.
var InstanceLock wisski_distillery.Command = instanceLock{}

type instanceLock struct {
	Lock        bool `description:"lock the provided instance"   long:"lock"   short:"l"`
	Unlock      bool `description:"unlock the provided instance" long:"unlock" short:"u"`
	Positionals struct {
		Slug string `description:"slug of instance to lock or unlock" positional-arg-name:"SLUG" required:"1-1"`
	} `positional-args:"true"`
}

func (instanceLock) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "instance_lock",
		Description: "locks or unlocks an instance",
	}
}

func (l instanceLock) AfterParse() error {
	if l.Lock == l.Unlock {
		return errLockUnlockExcluded
	}
	return nil
}

var (
	errLockUnlockExcluded = exit.NewErrorWithCode("exactly one of `--lock` and `--unlock` must be provied", exit.ExitCommandArguments)
	errNotUnlock          = exit.NewErrorWithCode("unable to unlock instance: not locked", exit.ExitCommandArguments)
	errInstanceLockWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
)

func (l instanceLock) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, l.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errInstanceLockWissKI, err)
	}

	if l.Unlock {
		if !instance.Locker().TryUnlock(context.Context) {
			return errNotUnlock
		}
		_, _ = context.Println("unlocked")
		return nil
	}

	if !instance.Locker().TryLock(context.Context) {
		return locker.ErrLocked
	}

	_, _ = context.Println("locked")
	return nil
}
