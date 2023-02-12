// Package stringparser provides Parser.
// It is deprecated and will be removed in a future release.
package stringparser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/pkg/errors"
)

// Parser is used to read a value from a string and turn it into a golang value.
// It is simultaniously used to validate particular setting.
//
// Parsers can be found in this package as functions called Parse*.
// They are refered to by their name, e.g. ParseNonempty can be refered to by the name 'Nonempty'.
// See [Parse].
type Parser[T any] func(env environment.Environment, s string) (T, error)

// ParseAbspath checks that s is an absolute path and returns it as-is
func ParseAbspath(env environment.Environment, s string) (string, error) {
	if !fsx.IsDirectory(env, s) {
		return "", errors.Errorf("%q does not exist or is not a directory", s)
	}
	return s, nil
}

// ParseFile checks that s is a valid file and returns it as-is
func ParseFile(env environment.Environment, s string) (string, error) {
	if !fsx.IsFile(env, s) {
		return "", errors.Errorf("%q does not exist or is not a regular file", s)
	}
	return s, nil
}

var errEmptyString = errors.New("value is empty")

// ParseNonEmpty checks that s is a non-empty string and returns it as-is
func ParseNonEmpty(env environment.Environment, s string) (string, error) {
	if s == "" {
		return "", errEmptyString
	}
	return s, nil
}

var regexpDomain = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

// ParseValidDomain checks that s is a valid domain and returns it in lowercase
func ParseValidDomain(env environment.Environment, s string) (string, error) {
	if !regexpDomain.MatchString(s) {
		return "", errors.Errorf("%q is not a valid domain", s)
	}
	return strings.ToLower(s), nil
}

// ParseValidDomains checks that s is a comma-seperated list of valid domains and returns them in lower case
func ParseValidDomains(env environment.Environment, s string) ([]string, error) {
	if len(s) == 0 {
		return []string{}, nil
	}
	domains := strings.Split(strings.ToLower(s), ",")
	for _, d := range domains {
		if !regexpDomain.MatchString(d) {
			return nil, errors.Errorf("%q is not a valid domain", d)
		}
	}
	return domains, nil
}

// ParseNumber parses s as a decimal integer
func ParseNumber(env environment.Environment, s string) (int, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	return int(value), err
}

// ParsePort parses s as a port
func ParsePort(env environment.Environment, s string) (uint16, error) {
	value, err := strconv.ParseUint(s, 10, 16)
	return uint16(value), err
}

// ParseHttpsURL parses a string into a url that starts with 'https://'
func ParseHttpsURL(env environment.Environment, s string) (*url.URL, error) {
	url, err := url.Parse(s)
	if err != nil {
		return nil, errors.Wrapf(err, "%q is not a valid URL", s)
	}
	if url.Scheme != "https" {
		return nil, errors.Errorf("%q is not a valid https URL (%q)", s, url.Scheme)
	}
	return url, nil
}

var regexpEmail = regexp.MustCompile(`^([-a-zA-Z0-9]+)\@([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

// ParseEmail checks that s represents an email, and then returns it as is.
func ParseEmail(env environment.Environment, s string) (string, error) {
	if s == "" { // no email provided
		return "", nil
	}
	if !regexpEmail.MatchString(s) {
		return "", errors.Errorf("%q is not a valid email", s)
	}
	return s, nil
}

var regexpSlug = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

// ParseSlug parses s as a slug and returns it as is.
func ParseSlug(env environment.Environment, s string) (string, error) {
	if !regexpSlug.MatchString(s) {
		return "", errors.Errorf("%q is not a valid slug", s)
	}
	return s, nil
}

// ParseDuration parses a time.Duration
func ParseDuration(env environment.Environment, s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
