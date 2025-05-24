//spellchecker:words validators
package validators

//spellchecker:words errors
import "errors"

//spellchecker:words github errors

var errEmpty = errors.New("value is empty")

func ValidateNonempty(value *string, dflt string) error {
	if *value == "" {
		*value = dflt
	}

	if *value == "" {
		return errEmpty
	}
	return nil
}
