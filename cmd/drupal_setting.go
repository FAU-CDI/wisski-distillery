package cmd

//spellchecker:words encoding json github wisski distillery internal cobra pkglib exit
import (
	"encoding/json"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewDrupalSettingCommand() *cobra.Command {
	impl := new(setting)

	cmd := &cobra.Command{
		Use:     "drupal_setting SLUG SETTING [VALUE]",
		Short:   "get or set a drupal setting",
		Args:    cobra.RangeArgs(2, 3),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type setting struct {
	Positionals struct {
		Slug    string
		Setting string
		Value   string
	}
}

func (ds *setting) ParseArgs(cmd *cobra.Command, args []string) error {
	ds.Positionals.Slug = args[0]
	ds.Positionals.Setting = args[1]
	if len(args) >= 3 {
		ds.Positionals.Value = args[2]
	}
	return nil
}

var (
	errSettingGet    = exit.NewErrorWithCode("unable to get setting", exit.ExitGeneric)
	errSettingSet    = exit.NewErrorWithCode("unable to set setting", exit.ExitGeneric)
	errSettingWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
)

func (ds *setting) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errSettingWissKI, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), ds.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errSettingWissKI, err)
	}

	if ds.Positionals.Value == "" {
		// get the setting
		value, err := instance.Settings().Get(cmd.Context(), nil, ds.Positionals.Setting)
		if err != nil {
			return fmt.Errorf("%w: %w", errSettingGet, err)
		}

		// and print it
		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(value); err != nil {
			return fmt.Errorf("%w: %w", errSettingGet, err)
		}

		// finish with a newline
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "")
		return nil
	}

	// serialize the setting into json
	var data any
	if err := json.Unmarshal([]byte(ds.Positionals.Value), &data); err != nil {
		return fmt.Errorf("%w: %w", errSettingSet, err)
	}

	// set the serialized value!
	if err := instance.Settings().Set(cmd.Context(), nil, ds.Positionals.Setting, data); err != nil {
		return fmt.Errorf("%w: %w", errSettingSet, err)
	}

	// and we're done
	return nil
}
