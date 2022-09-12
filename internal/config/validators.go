package config

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/pkg/errors"
)

// Validator reads from the configuration file
type Validator func(s string) (interface{}, error)

var knownValidators map[string]Validator = map[string]Validator{
	"is_valid_abspath":   IsValidAbspath,
	"is_valid_domain":    IsValidDomain,
	"is_valid_domains":   IsValidDomains,
	"is_valid_number":    IsValidNumber,
	"is_valid_https_url": IsValidHttpsURL,
	"is_valid_slug":      IsValidSlug,
	"is_valid_file":      IsValidFile,
	"is_valid_email":     IsValidEmail,
	"is_nonempty":        IsNonEmpty,
}

func IsValidAbspath(s string) (interface{}, error) {
	if !fsx.IsDirectory(s) {
		return nil, errors.Errorf("%q does not exist or is not a directory", s)
	}
	return s, nil
}

func IsValidFile(s string) (interface{}, error) {
	if !fsx.IsFile(s) {
		return nil, errors.Errorf("%q does not exist or is not a regular file", s)
	}
	return s, nil
}

func IsNonEmpty(s string) (interface{}, error) {
	if s == "" {
		return nil, errors.New("value is empty")
	}
	return s, nil
}

var regexpDomain = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

func IsValidDomain(s string) (interface{}, error) {
	if !regexpDomain.MatchString(s) {
		return nil, errors.Errorf("%q is not a valid domain", s)
	}
	return s, nil
}
func IsValidDomains(s string) (interface{}, error) {
	if len(s) == 0 {
		return []string{}, nil
	}
	domains := strings.Split(s, ",")
	for _, d := range domains {
		if !regexpDomain.MatchString(d) {
			return nil, errors.Errorf("%q is not a valid domain", d)
		}
	}
	return domains, nil
}

func IsValidNumber(s string) (interface{}, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	return int(value), err
}

func IsValidHttpsURL(s string) (interface{}, error) {
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

func IsValidEmail(s string) (interface{}, error) {
	if s == "" { // no email provided
		return "", nil
	}
	if !regexpEmail.MatchString(s) {
		return nil, errors.Errorf("%q is not a valid email", s)
	}
	return s, nil
}

var regexpSlug = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]*$`) // TODO: Make this regexp nicer!

func IsValidSlug(s string) (interface{}, error) {
	if !regexpSlug.MatchString(s) {
		return nil, errors.Errorf("%q is not a valid slug", s)
	}
	return s, nil
}
