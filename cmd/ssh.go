package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"
	"net"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewSSHCommand() *cobra.Command {
	impl := new(ssh)

	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "starts the ssh server to allow clients to connect to this distillery",
		Args:    cobra.NoArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.StringVar(&impl.Bind, "bind", "127.0.0.1:2223", "address to listen on")
	flags.StringVar(&impl.PrivateKeyPath, "private-key-path", "", "path to store private host keys in")
	if err := cmd.MarkFlagRequired("private-key-path"); err != nil {
		panic("failed to mark flag as required")
	}

	return cmd
}

type ssh struct {
	Bind           string
	PrivateKeyPath string
}

func (s *ssh) ParseArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func (*ssh) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "ssh",
		Description: "starts the ssh server to allow clients to connect to this distillery",
	}
}

var (
	errSSHServer    = exit.NewErrorWithCode("unable to listen server", exit.ExitGeneric)
	errSSHListen    = exit.NewErrorWithCode("unable to listen", exit.ExitGeneric)
	errServerListen = exit.NewErrorWithCode("server listen error", exit.ExitGeneric)
)

func (s *ssh) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errSSHServer, err)
	}

	server, err := dis.SSH().Server(cmd.Context(), s.PrivateKeyPath, cmd.ErrOrStderr())
	if err != nil {
		return fmt.Errorf("%w: %w", errSSHServer, err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Listening on %s\n", s.Bind)

	// make a new listener
	listener, err := net.Listen("tcp", s.Bind)
	if err != nil {
		return fmt.Errorf("%w: %w", errSSHListen, err)
	}

	go func() {
		<-cmd.Context().Done()
		_ = listener.Close() // it is either closed or it isn't
	}()

	// and serve that listener
	err = server.Serve(listener)
	if err != nil {
		return fmt.Errorf("%w: %w", errServerListen, err)
	}
	return nil
}
