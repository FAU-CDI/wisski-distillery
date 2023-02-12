package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tkw1536/goprogram/lib/reflectx"
)

// Collection represents a set of validators.
// The zero value is not ready to use; it should be created using make().
//
// A validator is a non-nil function with signature func(value *F, dflt string) error.
// Here F is the type of a value of a field.
// The value is the initialized value to be validated.
// The validator may perform abitrary normalization on the value.
// dflt is the default value (read from the 'default' tag).
// error should be an appropriate error that occured.
//
// A validator function is applied by calling it.
type Collection map[string]any

// Add adds a Validator to the provided collection of validators.
// Any previously validator of the same name is overwritten.
func Add[F any](coll Collection, name string, validator func(value *F, dflt string) error) {
	coll[name] = validator
}

// AddSlice adds a Validator to the provided collection of validators that validates a slice of the given type. The default is seperated by seperator.
func AddSlice[F any](coll Collection, name string, sep string, validator func(value *F, dflt string) error) {
	Add(coll, name, func(value *[]F, dflt string) error {
		// some value is set, so we do not need to set the default!
		if *value != nil {
			for i := range *value {
				if err := validator(&(*value)[i], ""); err != nil {
					return err
				}
			}
		}

		// no default provided => set if to an empty slice
		if dflt == "" {
			*value = make([]F, 0)
			return nil
		}

		// some default provided => iterate over the underlying validator
		dflts := strings.Split(dflt, sep)
		*value = make([]F, len(dflts))
		for i := range *value {
			if err := validator(&(*value)[i], dflts[i]); err != nil {
				return err
			}
		}

		return nil
	})
}

var (
	errTyp = reflectx.TypeOf[error]()
	strTyp = reflectx.TypeOf[string]()
)

// UnknownValidator is an error returned from Validate if a validator does not exist
type UnknownValidator string

func (uv UnknownValidator) Error() string {
	return fmt.Sprintf("unknown validator %q", string(uv))
}

// NotAValidator is an error returned from Validate if an entry in the validators map is not a validator
type NotAValidator string

func (nv NotAValidator) Error() string {
	return fmt.Sprintf("entry %q in validators is not a valiator", string(nv))
}

// IncompatibleValidator is returned when a validator in the validators map is incompatible
type IncompatibleValidator struct {
	Validator    string
	GotType      reflect.Type
	ExpectedType reflect.Type
}

func (iv IncompatibleValidator) Error() string {
	return fmt.Sprintf("validator %q: got type %s, expected type %s", iv.Validator, iv.GotType, iv.ExpectedType)
}

// Call calls the validator with the given name, on the given value, and with the provided default.
// See documentation of [Validate] for details.
func (coll Collection) Call(name string, field reflect.Value, dflt string) error {
	validator, ok := coll[name]
	if !ok {
		return UnknownValidator(name)
	}

	// get the type of the validator
	vFunc := reflect.ValueOf(validator)
	vTyp := vFunc.Type()

	// ensure that vTyp is of type func(*F,string) error
	// where T is the type of the field
	//
	// - the first if assumes checks for some type F
	// - the second if checks if the F is the right one
	if validator == nil || vTyp.Kind() != reflect.Func || // func
		vTyp.NumIn() != 2 || vTyp.In(0).Kind() != reflect.Pointer || vTyp.In(1) != strTyp || // (*F,string)
		vTyp.NumOut() != 1 || vTyp.Out(0) != errTyp { // error
		return NotAValidator(name)
	}
	if vTyp.In(0).Elem() != field.Type() { // the correct *F
		return IncompatibleValidator{
			Validator:    name,
			GotType:      vTyp.In(0).Elem(),
			ExpectedType: field.Type(),
		}
	}

	// call the validator function, and return an error
	results := vFunc.Call([]reflect.Value{field.Addr(), reflect.ValueOf(dflt)})

	// turn the result into an error
	// NOTE: We can't just .(error) here because that panic()s on err == nil
	err := results[0].Interface()
	if err, ok := err.(error); ok {
		return err
	}
	return nil
}
