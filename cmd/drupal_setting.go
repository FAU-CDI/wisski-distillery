package cmd

import (
	"encoding/json"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// DrupalSetting is then 'drupal_setting' command
var DrupalSetting wisski_distillery.Command = setting{}

type setting struct {
	Positionals struct {
		Slug    string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to get or set value for"`
		Setting string `positional-arg-name:"SETTING" require:"1-1" description:"name of setting to read or write"`
		Value   string `positional-arg-name:"VALUE" description:"json serialization of value to write"`
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

var errSettingGet = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to get setting",
}

var errSettingSet = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to set setting",
}

var errSettingWissKI = exit.Error{
	Message:  "unable to get WissKI",
	ExitCode: exit.ExitGeneric,
}

func (ds setting) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, ds.Positionals.Slug)
	if err != nil {
		return errSettingWissKI.Wrap(err)
	}

	if ds.Positionals.Value == "" {
		// get the setting
		value, err := instance.Settings().Get(context.Context, nil, ds.Positionals.Setting)
		if err != nil {
			return errSettingGet.Wrap(err)
		}

		// and print it
		if err := json.NewEncoder(context.Stdout).Encode(value); err != nil {
			return errSettingGet.Wrap(err)
		}

		// finish with a newline
		context.Println("")
		return nil
	}

	// serialize the setting into json
	var data any
	if err := json.Unmarshal([]byte(ds.Positionals.Value), &data); err != nil {
		return errSettingSet.Wrap(err)
	}

	// set the serialized value!
	if err := instance.Settings().Set(context.Context, nil, ds.Positionals.Setting, data); err != nil {
		return errSettingSet.Wrap(err)
	}

	// and we're done
	return nil
}
