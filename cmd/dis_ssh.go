package cmd

import (
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/tkw1536/goprogram/exit"

	gossh "golang.org/x/crypto/ssh"
)

// DisSSH is the 'dis_ssh' command
var DisSSH wisski_distillery.Command = disSSH{}

type disSSH struct {
	Add     bool   `short:"a" long:"add" description:"add key to user"`
	Remove  bool   `short:"r" long:"remove" description:"remove key from user"`
	Comment string `short:"c" long:"comment" description:"comment of new key"`

	Positionals struct {
		User string `positional-arg-name:"USER" required:"1-1" description:"distillery username"`
		Path string `positional-arg-name:"PATH" required:"1-1" description:"path of key to add"`
	} `positional-args:"true"`
}

func (disSSH) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_ssh",
		Description: "add or remove an ssh key from a user",
	}
}

func (ds disSSH) AfterParse() error {
	var counter int
	for _, action := range []bool{
		ds.Add,
		ds.Remove,
	} {
		if action {
			counter++
		}
	}

	if counter != 1 {
		return errNoActionSelected
	}

	return nil
}

var errSSHManageFailed = exit.Error{
	Message:  "unable to manage ssh keys",
	ExitCode: exit.ExitCommandArguments,
}

func (ds disSSH) Run(context wisski_distillery.Context) error {
	switch {
	case ds.Add:
		return ds.runAdd(context)
	case ds.Remove:
		return ds.runRemove(context)
	}
	panic("never reached")
}

var errNoKey = exit.Error{
	Message:  "unable to parse key",
	ExitCode: exit.ExitCommandArguments,
}

func (ds disSSH) parseOpts(context wisski_distillery.Context) (user *auth.AuthUser, key gossh.PublicKey, err error) {
	user, err = context.Environment.Auth().User(context.Context, ds.Positionals.User)
	if err != nil {
		return nil, nil, errSSHManageFailed.Wrap(err)
	}

	content, err := os.ReadFile(ds.Positionals.Path)
	if err != nil {
		return nil, nil, errSSHManageFailed.Wrap(err)
	}

	pk, _, _, _, err := gossh.ParseAuthorizedKey(content)
	if pk == nil || err != nil {
		return nil, nil, errNoKey
	}

	return user, pk, nil
}

func (ds disSSH) runAdd(context wisski_distillery.Context) error {
	user, key, err := ds.parseOpts(context)
	if err != nil {
		return err
	}

	if err := context.Environment.Keys().Add(context.Context, user.User.User, ds.Comment, key); err != nil {
		return errSSHManageFailed.Wrap(err)
	}
	return nil
}

func (ds disSSH) runRemove(context wisski_distillery.Context) error {
	user, key, err := ds.parseOpts(context)
	if err != nil {
		return err
	}

	if err := context.Environment.Keys().Remove(context.Context, user.User.User, key); err != nil {
		return errSSHManageFailed.Wrap(err)
	}
	return nil
}
