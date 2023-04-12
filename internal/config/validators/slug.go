package validators

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var regexpSlug = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

var ErrInvalidSlug = errors.New("invalid slug")

// ValidateSlug validates a slug and normalizes it.
func ValidateSlug(s *string, dflt string) error {
	if *s == "" {
		*s = dflt
	}
	*s = strings.ToLower(*s)
	if !regexpSlug.MatchString(*s) {
		return ErrInvalidSlug
	}
	if strings.HasSuffix(*s, "_") {
		return ErrInvalidSlug
	}
	return nil
}
