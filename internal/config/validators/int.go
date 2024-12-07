//spellchecker:words validators
package validators

//spellchecker:words strconv github errors
import (
	"strconv"

	"github.com/pkg/errors"
)

func ValidatePositive(value *int, dflt string) (err error) {
	if *value == 0 && dflt != "" {
		v, err := strconv.ParseInt(dflt, 10, 64)
		if err != nil {
			return err
		}
		*value = int(v)
	}
	if *value <= 0 {
		return errors.Errorf("%d is not a positive value", *value)
	}
	return nil
}

func ValidatePort(value *uint16, dflt string) (err error) {
	if *value == 0 && dflt != "" {
		v, err := strconv.ParseUint(dflt, 10, 16)
		if err != nil {
			return err
		}
		*value = uint16(v)
	}
	return nil
}
