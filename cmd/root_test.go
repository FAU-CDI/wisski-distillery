package cmd_test

import (
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/FAU-CDI/wisski-distillery/cmd"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/stream"
)

func Test_NewCommand(t *testing.T) {
	t.Parallel()

	// find all the commands that are available
	root := cmd.NewCommand(t.Context(), cli.Params{})
	commands := findCommands(root, func(cmd *cobra.Command) bool { return !cmd.DisableFlagParsing })

	// run each of the commanbds with the --help flag
	for _, argv := range commands {
		t.Run(strings.Join(argv, " "), func(t *testing.T) {
			t.Parallel()

			root := cmd.NewCommand(t.Context(), cli.Params{})
			root.SetArgs(append(argv[1:], "--help"))
			root.SetIn(stream.Null)
			root.SetOut(io.Discard)
			root.SetErr(io.Discard)

			err := root.Execute()
			if err != nil {
				t.Errorf("error executing command: %v", err)
			}
		})
	}
}

func findCommands(cmd *cobra.Command, include func(cmd *cobra.Command) bool) (commands [][]string) {
	var walkTree func(cmd *cobra.Command, path []string)
	walkTree = func(cmd *cobra.Command, path []string) {
		us := slices.Clone(path)
		us = append(us, cmd.Name())

		if include(cmd) {
			commands = append(commands, us)
		}
		for _, sub := range cmd.Commands() {
			walkTree(sub, us)
		}
	}
	walkTree(cmd, nil)
	return
}
