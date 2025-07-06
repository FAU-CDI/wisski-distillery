package cmd

//spellchecker:words github wisski distillery internal component auth goprogram exit golang crypto gossh
import (
	"fmt"
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"go.tkw01536.de/goprogram/exit"

	gossh "golang.org/x/crypto/ssh"
)

// DisSSH is the 'dis_ssh' command.
var DisSSH wisski_distillery.Command = disSSH{}

type disSSH struct {
	Add     bool   `description:"add key to user"      long:"add"     short:"a"`
	Remove  bool   `description:"remove key from user" long:"remove"  short:"r"`
	Comment string `description:"comment of new key"   long:"comment" short:"c"`

	Positionals struct {
		User string `description:"distillery username" positional-arg-name:"USER" required:"1-1"`
		Path string `description:"path of key to add"  positional-arg-name:"PATH" required:"1-1"`
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

var errSSHManageFailed = exit.NewErrorWithCode("unable to manage ssh keys", exit.ExitCommandArguments)

func (ds disSSH) Run(context wisski_distillery.Context) error {
	switch {
	case ds.Add:
		return ds.runAdd(context)
	case ds.Remove:
		return ds.runRemove(context)
	}
	panic("never reached")
}

var errNoKey = exit.NewErrorWithCode("unable to parse key", exit.ExitCommandArguments)

func (ds disSSH) parseOpts(context wisski_distillery.Context) (user *auth.AuthUser, key gossh.PublicKey, err error) {
	user, err = context.Environment.Auth().User(context.Context, ds.Positionals.User)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}

	content, err := os.ReadFile(ds.Positionals.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", errSSHManageFailed, err)
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
		return fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}
	return nil
}

func (ds disSSH) runRemove(context wisski_distillery.Context) error {
	user, key, err := ds.parseOpts(context)
	if err != nil {
		return err
	}

	if err := context.Environment.Keys().Remove(context.Context, user.User.User, key); err != nil {
		return fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}
	return nil
}
