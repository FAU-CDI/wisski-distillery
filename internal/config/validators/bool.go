package validators

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

// NullableBool represents a bool that can be null
type NullableBool struct {
	Null, Value bool
}

func (nb *NullableBool) UnmarshalYAML(value *yaml.Node) error {
	nb.Null = false
	if err := value.Decode(&nb.Value); err != nil {
		nb.Null = true
		nb.Value = false
	}

	return nil
}

func (nb *NullableBool) MarshalYAML() (interface{}, error) {
	if nb.Null {
		return nil, nil
	}
	return nb.Value, nil
}

func ValidateBool(value *NullableBool, dflt string) (err error) {
	if value.Null {
		res, err := strconv.ParseBool(dflt)
		if err != nil {
			return err
		}
		value.Null = false
		value.Value = res
	}
	return err
}
