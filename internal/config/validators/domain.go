//spellchecker:words validators
package validators

//spellchecker:words regexp strings github errors
import (
	"fmt"
	"regexp"
	"strings"
)

var regexpDomain = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

type invalidDomainError string

func (ide invalidDomainError) Error() string {
	return fmt.Sprintf("%q is not a valid domain", string(ide))
}

func ValidateDomain(domain *string, dflt string) error {
	if *domain == "" {
		*domain = dflt
	}
	if !regexpDomain.MatchString(*domain) {
		return invalidDomainError(*domain)
	}
	*domain = strings.ToLower(*domain)
	return nil
}
