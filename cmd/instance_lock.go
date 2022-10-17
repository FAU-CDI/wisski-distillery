package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// InstanceLock is then 'instance_lock' command
var InstanceLock wisski_distillery.Command = instanceLock{}

type instanceLock struct {
	Lock        bool `short:"l" long:"lock" description:"Lock the provided WissKI instance"`
	Unlock      bool `short:"u" long:"unlock" description:"Unlock the provided WissKI instance"`
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
		Description: "Locks or unlocks a WissKI instance",
	}
}

var errLockUnlockExcluded = exit.Error{
	Message:  "Exactly one of `--lock` and `--unlock` must be provied",
	ExitCode: exit.ExitCommandArguments,
}

func (l instanceLock) AfterParse() error {
	if l.Lock == l.Unlock {
		return errLockUnlockExcluded
	}
	return nil
}

var errNotUnlock = exit.Error{
	Message:  "Unable to unlock instance: Not locked",
	ExitCode: exit.ExitCommandArguments,
}

func (l instanceLock) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(l.Positionals.Slug)
	if err != nil {
		return err
	}

	if l.Unlock {
		if !instance.Unlock() {
			return errNotUnlock
		}
		context.Println("unlocked")
		return nil
	}

	if err := instance.TryLock(); err != nil {
		return err
	}

	context.Println("locked")
	return nil
}
