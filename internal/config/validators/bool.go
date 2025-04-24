//spellchecker:words validators
package validators

// TODO: Figure out if there is an existing package for this!

//spellchecker:words strconv gopkg yaml
import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

// NullableBool represents a bool that can be null.
type NullableBool struct {
	Set, Value bool
}

func (nb *NullableBool) UnmarshalYAML(value *yaml.Node) error {
	nb.Set = true
	if err := value.Decode(&nb.Value); err != nil {
		nb.Set = false
		nb.Value = false
	}

	return nil
}

func (nb NullableBool) MarshalYAML() (interface{}, error) {
	if !nb.Set {
		return nil, nil
	}
	return nb.Value, nil
}

func ValidateBool(value *NullableBool, dflt string) (err error) {
	if !value.Set {
		res, err := strconv.ParseBool(dflt)
		if err != nil {
			return fmt.Errorf("failed to parse boolean: %w", err)
		}
		value.Set = true
		value.Value = res
	}
	return err
}
