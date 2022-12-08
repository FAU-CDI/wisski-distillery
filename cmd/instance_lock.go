package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/tkw1536/goprogram/exit"
)

// InstanceLock is then 'instance_lock' command
var InstanceLock wisski_distillery.Command = instanceLock{}

type instanceLock struct {
	Lock        bool `short:"l" long:"lock" description:"lock the provided instance"`
	Unlock      bool `short:"u" long:"unlock" description:"unlock the provided instance"`
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to lock or unlock"`
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

var errLockUnlockExcluded = exit.Error{
	Message:  "exactly one of `--lock` and `--unlock` must be provied",
	ExitCode: exit.ExitCommandArguments,
}

func (l instanceLock) AfterParse() error {
	if l.Lock == l.Unlock {
		return errLockUnlockExcluded
	}
	return nil
}

var errNotUnlock = exit.Error{
	Message:  "unable to unlock instance: not locked",
	ExitCode: exit.ExitCommandArguments,
}

func (l instanceLock) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, l.Positionals.Slug)
	if err != nil {
		return err
	}

	if l.Unlock {
		if !instance.Locker().TryUnlock(context.Context) {
			return errNotUnlock
		}
		context.Println("unlocked")
		return nil
	}

	if !instance.Locker().TryLock(context.Context) {
		return locker.Locked
	}

	context.Println("locked")
	return nil
}
