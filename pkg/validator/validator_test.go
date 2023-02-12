package validator

import (
	"errors"
	"fmt"
	"strconv"
)

func ExampleValidate() {
	var value struct {
		Number    int    `validate:"positive" default:"234"`
		String    string `validate:"nonempty" default:"stuff"`
		Recursive struct {
			Number int    `validate:"positive" default:"45"`
			String string `validate:"nonempty" default:"more"`
		} `recurse:"true"`
	}

	collection := make(Collection, 2)
	Add(collection, "positive", func(value *int, dflt string) error {
		if *value == 0 {
			i, err := strconv.ParseInt(dflt, 10, 64)
			if err != nil {
				return err
			}
			*value = int(i)
			return nil
		}
		if *value < 0 {
			return errors.New("not positive")
		}
		return nil
	})
	Add(collection, "nonempty", func(value *string, dflt string) error {
		if *value == "" {
			*value = dflt
		}
		if *value == "" {
			return errors.New("empty string")
		}
		return nil
	})

	err := Validate(&value, collection)
	fmt.Printf("%v\n", value)
	fmt.Println(err)
	// Output: {234 stuff {45 more}}
	// <nil>
}

func ExampleValidate_fail() {
	var value struct {
		Number    int    `validate:"positive" default:"12"`
		String    string `validate:"nonempty" default:"stuff"`
		Recursive struct {
			Number int    `validate:"positive" default:"12"`
			String string `validate:"nonempty"`
		} `recurse:"true"`
	}

	collection := make(Collection, 2)
	Add(collection, "positive", func(value *int, dflt string) error {
		if *value == 0 {
			i, err := strconv.ParseInt(dflt, 10, 64)
			if err != nil {
				return err
			}
			*value = int(i)
			return nil
		}
		if *value < 0 {
			return errors.New("not positive")
		}
		return nil
	})
	Add(collection, "nonempty", func(value *string, dflt string) error {
		if *value == "" {
			*value = dflt
		}
		if *value == "" {
			return errors.New("empty string")
		}
		return nil
	})

	err := Validate(&value, collection)
	fmt.Printf("%v\n", value)
	fmt.Println(err)
	// Output: {12 stuff {12 }}
	// field "Recursive": field "String": empty string
}

func ExampleValidate_notastruct() {
	var value int
	err := Validate(&value, nil)
	fmt.Println(err)
}

func ExampleValidate_notavalidator() {
	var value struct {
		Field int `validate:"generic"`
	}
	collection := make(Collection, 2)
	collection["generic"] = func(x, y int) error {
		panic("never reached")
	}
	err := Validate(&value, collection)
	fmt.Println(err)
	// Output: field "Field": entry "generic" in validators is not a valiator
}

func ExampleValidate_invalid() {
	var value struct {
		Field int `validate:"string"`
	}
	collection := make(Collection, 2)
	collection["string"] = func(value *string, dflt string) error {
		panic("never reached")
	}
	err := Validate(&value, collection)
	fmt.Println(err)
	// Output: field "Field": validator "string": got type string, expected type int
}
