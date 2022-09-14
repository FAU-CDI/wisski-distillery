package config

import (
	"io"
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/pkg/envreader"
	"github.com/FAU-CDI/wisski-distillery/pkg/stringparser"
)

// Unmarshal updates this configuration from the provided [io.Reader].
//
// Data is read using the [envreader.ReadAll] method, see the appropriate documentation for the file format.
//
// The `env` and `parser` reflect tags of the [Config] struct determine the keys to read from, and the types to expect.
// When a key is missing, it is set to the default value.
//
// See also [stringparser.Parse].
func (config *Config) Unmarshal(src io.Reader) error {
	// read all the values!
	values, err := envreader.ReadAll(src)
	if err != nil {
		return err
	}

	vConfig := reflect.ValueOf(config).Elem()
	tConfig := vConfig.Type()

	// iterate over the types
	numValues := tConfig.NumField()
	for i := 0; i < numValues; i++ {
		tField := tConfig.Field(i)
		vField := vConfig.FieldByName(tField.Name)

		env := tField.Tag.Get("env")
		dflt := tField.Tag.Get("default")
		parser := tField.Tag.Get("parser")

		// skip it if it isn't loaded!
		if env == "" {
			continue
		}

		// read the value with a default
		value, ok := values[env]
		if !ok || value == "" {
			if dflt == "" {
				continue
			}
			value = dflt
		}

		// parse the value!
		if err := stringparser.Parse(parser, value, vField); err != nil {
			return err
		}
	}

	return nil
}
