package validators

import (
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/pkglib/validator"
)

// New creates a new set of standard validators for the configuration
func New(env environment.Environment) validator.Collection {
	coll := make(validator.Collection)

	validator.Add(coll, "nonempty", ValidateNonempty)

	validator.Add(coll, "directory", func(value *string, dflt string) error {
		return ValidateDirectory(env, value, dflt)
	})
	validator.Add(coll, "file", func(value *string, dflt string) error {
		return ValidateFile(env, value, dflt)
	})

	validator.Add(coll, "domain", ValidateDomain)
	validator.AddSlice(coll, "domains", ",", ValidateDomain)
	validator.Add(coll, "https", ValidateHTTPSURL)
	validator.Add(coll, "slug", ValidateSlug)
	validator.Add(coll, "email", ValidateEmail)

	validator.Add(coll, "positive", ValidatePositive)
	validator.Add(coll, "port", ValidatePort)

	validator.Add(coll, "duration", ValidateDuration)
	return coll
}
