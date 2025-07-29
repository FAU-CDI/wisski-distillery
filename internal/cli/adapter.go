package cli

//spellchecker:words context github cobra
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/spf13/cobra"
)

type cobraKey int

const (
	flagsKey cobraKey = iota
	parametersKey
)

// GetDistillery gets the distillery for the currently running command.
// [SetFlags] and [SetParameters] must have been called.
func GetDistillery(cmd *cobra.Command, requirements Requirements) (*dis.Distillery, error) {
	// TODO: merge these functions together
	return NewDistillery(get[Params](cmd, parametersKey), get[Flags](cmd, flagsKey), requirements)
}

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
