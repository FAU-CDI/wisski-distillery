package validators

import (
	"regexp"

	"github.com/pkg/errors"
)

var regexpEmail = regexp.MustCompile(`^([-a-zA-Z0-9]+)\@([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

// ValidateEmail checks that s represents an email, and then returns it as is.
func ValidateEmail(email *string, dflt string) error {
	if *email == "" {
		*email = dflt
	}
	if *email == "" { // no email provided => ok
		return nil
	}

	if !regexpEmail.MatchString(*email) {
		return errors.Errorf("%q is not a valid email", *email)
	}
	return nil
}
