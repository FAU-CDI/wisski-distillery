//spellchecker:words config
package config

//spellchecker:words github wisski distillery internal config validators pkglib validator gopkg yaml
import (
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
	"github.com/tkw1536/pkglib/validator"
	"gopkg.in/yaml.v3"
)

// Unmarshal reads configuration from the provided io.Reader, and then validates it.
// Configuration is read in yaml format.
func (config *Config) Unmarshal(src io.Reader) error {
	// read yaml!
	{
		decoder := yaml.NewDecoder(src)
		decoder.KnownFields(true)
		if err := decoder.Decode(config); err != nil {
			return err
		}
	}

	// TODO: should this be done seperatly?
	return config.Validate()
}

// Validate validates this configuration file and sets appropriate defaults
func (config *Config) Validate() error {
	return validator.Validate(config, validators.New())
}

func (config *Config) Marshal(dest io.Writer) error {
	encoder := yaml.NewEncoder(dest)
	encoder.SetIndent(4)
	return encoder.Encode(config)
}
