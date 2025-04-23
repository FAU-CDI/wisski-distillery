//spellchecker:words validators
package validators

//spellchecker:words regexp github errors
import (
	"fmt"
	"regexp"

	"errors"
)

var regexpEmail = regexp.MustCompile(`^([-a-zA-Z0-9]+)\@([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

var errNotAValidEmail = errors.New("not a valid email")

// ValidateEmail checks that s represents an email, and then returns it as is.
func ValidateEmail(email *string, dflt string) error {
	if *email == "" {
		*email = dflt
	}
	if *email == "" { // no email provided => ok
		return nil
	}

	if !regexpEmail.MatchString(*email) {
		return fmt.Errorf("%q: %w", *email, errNotAValidEmail)
	}
	return nil
}
