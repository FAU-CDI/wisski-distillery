package stringparser

import (
	"reflect"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/pkg/errors"
)

// Parse parses the provided value with the parser.
func Parse(env environment.Environment, name, value string, vField reflect.Value) error {

	// use the validator
	parser, ok := knownParsers[strings.ToLower(name)]
	if parser == nil || !ok {
		return errors.Errorf("unknown parser %q", name)
	}

	// get the parsed value
	checked, err := parser(env, value)
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
	return func(env environment.Environment, s string) (value any, err error) {
		value, err = parser(env, s)
		return
	}
}