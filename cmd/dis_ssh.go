package cmd

//spellchecker:words github wisski distillery internal component auth goprogram exit golang crypto gossh
import (
	"fmt"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"

	gossh "golang.org/x/crypto/ssh"
)

func NewDisSSHCommand() *cobra.Command {
	impl := new(disSSH)

	cmd := &cobra.Command{
		Use:     "dis_ssh USER PATH",
		Short:   "add or remove an ssh key from a user",
		Args:    cobra.ExactArgs(2),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Add, "add", false, "add key to user")
	flags.BoolVar(&impl.Remove, "remove", false, "remove key from user")
	flags.StringVar(&impl.Comment, "comment", "", "comment of new key")

	return cmd
}

type disSSH struct {
	Add         bool
	Remove      bool
	Comment     string
	Positionals struct {
		User string
		Path string
	}
}

func (ds *disSSH) ParseArgs(cmd *cobra.Command, args []string) error {
	ds.Positionals.User = args[0]
	ds.Positionals.Path = args[1]

	// Validate arguments
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

func (ds *disSSH) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}

	switch {
	case ds.Add:
		return ds.runAdd(cmd, dis)
	case ds.Remove:
		return ds.runRemove(cmd, dis)
	}
	panic("never reached")
}

var errNoKey = exit.NewErrorWithCode("unable to parse key", exit.ExitCommandArguments)

func (ds *disSSH) parseOpts(cmd *cobra.Command, dis *dis.Distillery) (user *auth.AuthUser, key gossh.PublicKey, err error) {
	user, err = dis.Auth().User(cmd.Context(), ds.Positionals.User)
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

func (ds *disSSH) runAdd(cmd *cobra.Command, dis *dis.Distillery) error {
	user, key, err := ds.parseOpts(cmd, dis)
	if err != nil {
		return err
	}

	if err := dis.Keys().Add(cmd.Context(), user.User.User, ds.Comment, key); err != nil {
		return fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}
	return nil
}

func (ds *disSSH) runRemove(cmd *cobra.Command, dis *dis.Distillery) error {
	user, key, err := ds.parseOpts(cmd, dis)
	if err != nil {
		return err
	}

	if err := dis.Keys().Remove(cmd.Context(), user.User.User, key); err != nil {
		return fmt.Errorf("%w: %w", errSSHManageFailed, err)
	}
	return nil
}
