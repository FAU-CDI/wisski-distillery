package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit
import (
	"encoding/json"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// DrupalSetting is then 'drupal_setting' command.
var DrupalSetting wisski_distillery.Command = setting{}

type setting struct {
	Positionals struct {
		Slug    string `description:"slug of instance to get or set value for" positional-arg-name:"SLUG"    required:"1-1"`
		Setting string `description:"name of setting to read or write"         positional-arg-name:"SETTING" require:"1-1"`
		Value   string `description:"json serialization of value to write"     positional-arg-name:"VALUE"`
	} `positional-args:"true"`
}

func (setting) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "drupal_setting",
		Description: "get or set a drupal setting",
	}
}

var (
	errSettingGet    = exit.NewErrorWithCode("unable to get setting", exit.ExitGeneric)
	errSettingSet    = exit.NewErrorWithCode("unable to set setting", exit.ExitGeneric)
	errSettingWissKI = exit.NewErrorWithCode("unable to get WissKI", exit.ExitGeneric)
)

func (ds setting) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, ds.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errSettingWissKI, err)
	}

	if ds.Positionals.Value == "" {
		// get the setting
		value, err := instance.Settings().Get(context.Context, nil, ds.Positionals.Setting)
		if err != nil {
			return fmt.Errorf("%w: %w", errSettingGet, err)
		}

		// and print it
		if err := json.NewEncoder(context.Stdout).Encode(value); err != nil {
			return fmt.Errorf("%w: %w", errSettingGet, err)
		}

		// finish with a newline
		_, _ = context.Println("")
		return nil
	}

	// serialize the setting into json
	var data any
	if err := json.Unmarshal([]byte(ds.Positionals.Value), &data); err != nil {
		return fmt.Errorf("%w: %w", errSettingSet, err)
	}

	// set the serialized value!
	if err := instance.Settings().Set(context.Context, nil, ds.Positionals.Setting, data); err != nil {
		return fmt.Errorf("%w: %w", errSettingSet, err)
	}

	// and we're done
	return nil
}
