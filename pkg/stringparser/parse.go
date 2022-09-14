package stringparser

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Parse parses the provided value with the parser.
func Parse(name, value string, vField reflect.Value) error {

	// use the validator
	parser, ok := knownParsers[strings.ToLower(name)]
	if parser == nil || !ok {
		return errors.Errorf("unknown parser %q", name)
	}

	// get the parsed value
	checked, err := parser(value)
	if err != nil {
		return errors.Wrapf(err, "parser %s returned error", name)
	}

	// set the value of the field
	var errSet interface{}
	func() {
		defer func() {
			errSet = recover()
		}()
		vField.Set(reflect.ValueOf(checked))
	}()

	// capture any error
	if errSet != nil {
		return errors.Errorf("parser %s: set returned %v", name, errSet)
	}

	return nil
}

// knownParsers holds the known parsers
var knownParsers map[string]Parser[any] = map[string]Parser[any]{
	"abspath":   asGenericParser(ParseAbspath),
	"domain":    asGenericParser(ParseValidDomain),
	"domains":   asGenericParser(ParseValidDomains),
	"number":    asGenericParser(ParseNumber),
	"https_url": asGenericParser(ParseHttpsURL),
	"slug":      asGenericParser(ParseSlug),
	"file":      asGenericParser(ParseFile),
	"email":     asGenericParser(ParseEmail),
	"nonempty":  asGenericParser(ParseNonEmpty),
}

func asGenericParser[T any](parser Parser[T]) Parser[any] {
	return func(s string) (value any, err error) {
		value, err = parser(s)
		return
	}
}
