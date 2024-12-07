//spellchecker:words validators
package validators

//spellchecker:words regexp strings github errors
import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var regexpDomain = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

func ValidateDomain(domain *string, dflt string) error {
	if *domain == "" {
		*domain = dflt
	}
	if !regexpDomain.MatchString(*domain) {
		return errors.Errorf("%q is not a valid domain", *domain)
	}
	*domain = strings.ToLower(*domain)
	return nil
}
