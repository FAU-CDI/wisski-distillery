//spellchecker:words validators
package validators

//spellchecker:words github errors
import "github.com/pkg/errors"

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
