package cli

//spellchecker:words context github wisski distillery internal cobra
import (
	"context"

	"github.com/spf13/cobra"
)

type cobraKey int

const (
	flagsKey cobraKey = iota
	parametersKey
)

// SetFlags sets the value for a cobra command from a set of flags.
func SetFlags(cmd *cobra.Command, flags *Flags) {
	set(cmd, flagsKey, flags)
}

// SetParameters sets parameters for a cobra command.
func SetParameters(cmd *cobra.Command, params *Params) {
	set(cmd, parametersKey, params)
}

func set[T any](cmd *cobra.Command, key cobraKey, data *T) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(context.WithValue(ctx, key, data))
}

func get[T any](cmd *cobra.Command, key cobraKey) T {
	flags := cmd.Context().Value(key)
	data, ok := flags.(*T)
	if !ok {
		var zero T
		return zero
	}
	return *data
}
