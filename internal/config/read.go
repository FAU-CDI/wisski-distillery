package config

import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/validator"
	"gopkg.in/yaml.v3"
)

// Unmarshal reads configuration from the provided io.Reader, and then validates it.
// Configuration is read in yaml format.
func (config *Config) Unmarshal(env environment.Environment, src io.Reader) error {
	// read yaml!
	{
		decoder := yaml.NewDecoder(src)
		decoder.KnownFields(true)
		if err := decoder.Decode(config); err != nil {
			return err
		}
	}

	// TODO: should this be done seperatly?
	return config.Validate(env)
}

// Validate validates this configuration file and sets appropriate defaults
func (config *Config) Validate(env environment.Environment) error {
	return validator.Validate(config, validators.New(env))
}

func (config *Config) Marshal(dest io.Writer) error {
	encoder := yaml.NewEncoder(dest)
	encoder.SetIndent(4)
	return encoder.Encode(config)
}
