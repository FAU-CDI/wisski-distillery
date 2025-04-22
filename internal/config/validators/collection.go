//spellchecker:words validators
package validators

//spellchecker:words github pkglib validator
import (
	"github.com/tkw1536/pkglib/validator"
)

// New creates a new set of standard validators for the configuration.
func New() validator.Collection {
	coll := make(validator.Collection)

	validator.Add(coll, "nonempty", ValidateNonempty)

	validator.Add(coll, "bool", ValidateBool)

	validator.Add(coll, "directory", ValidateDirectory)
	validator.Add(coll, "file", ValidateFile)

	validator.Add(coll, "domain", ValidateDomain)
	validator.AddSlice(coll, "domains", ",", ValidateDomain)
	validator.Add(coll, "https", ValidateHTTPSURL)
	validator.Add(coll, "slug", ValidateSlug)
	validator.Add(coll, "email", ValidateEmail)

	validator.Add(coll, "positive", ValidatePositive)
	validator.Add(coll, "port", ValidatePort)
	validator.AddSlice(coll, "ports", ",", ValidatePort)

	validator.Add(coll, "duration", ValidateDuration)
	return coll
}
