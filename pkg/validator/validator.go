package validator

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	validateTag = "validate"
	recurseTag  = "recurse"
	dfltTag     = "default"
)

// Validate validates an object of type T, setting defaults where appropriate.
//
// T must be a struct type, when this is not the case, returns ErrNotAStruct.
// validators should contain a set of validators.
//
// Validate iterates over the fields and tags of those fields as follows:
//   - If the 'validate' tag is not the empty string, read the appropriate validator from the map, and call the function.
//     If the element in the validators map does not exist, returns an error that unwraps to type UnknownValidator.
//     If the element in the validators map is not a validator, returns an error that unwraps to type NotAValidator.
//     If the type of validator function does not match the field type, returns an error that unwraps to type IncompatibleValidator.
//   - If the 'recurse' tag is not the empty string, recurse into the struct type by calling Validate on it.
//     If the annotated field is not a struct, return an error.
//
// Any error is wrapped in a FieldError, indicating the field they occured in.
// Recursive validate calls may result in FieldError wraps.
// For a description of struct tags, see [reflect.StructTag].
func Validate[T any](data *T, validators map[string]any) error {
	return validate(reflect.ValueOf(data).Elem(), validators)
}

// FieldError wraps an error to indicate which field it occured in.
type FieldError struct {
	Field string
	Err   error
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("field %q: %s", fe.Field, fe.Err)
}

func (fe FieldError) Unwrap() error {
	return fe.Err
}

var ErrNotAStruct = errors.New("validate called on non-struct type")

func validate(datum reflect.Value, validators Collection) error {
	// make sure that we have a struct type
	typ := datum.Type()
	if typ.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	fieldC := typ.NumField()
	for i := 0; i < fieldC; i++ {
		field := typ.Field(i)

		// if the recurse tag is set, do the recursion!
		if field.Tag.Get(recurseTag) != "" {
			if err := validate(datum.FieldByName(field.Name), validators); err != nil {
				return FieldError{Field: field.Name, Err: err}
			}
			continue
		}

		// check if there is a validator associated with this tag
		// and if not, skip it!
		validator := field.Tag.Get(validateTag)
		if validator == "" {
			continue
		}

		// call the actual validator
		if err := validators.Call(
			validator,
			datum.FieldByName(field.Name),
			field.Tag.Get(dfltTag),
		); err != nil {
			return FieldError{Field: field.Name, Err: err}
		}
	}

	return nil
}
