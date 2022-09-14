// Package envreader
package envreader

import (
	"fmt"
	"strings"
)

func ExampleNewScanner() {
	scanner := NewScanner(strings.NewReader(`
lines without an equal sign are ignored

// this line is a comment, even with an = sign
KEY=VALUE

# this is also a comment =
spaces in keys = spaces in values
multiple=equal=signs
CaSe = SenSitiVe
empty value=
=empty key
`))

	for scanner.Scan() {
		key, value := scanner.Data()
		fmt.Printf("%q %q\n", key, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(scanner.Err())
	} else {
		fmt.Println("no error")
	}

	// Output: "KEY" "VALUE"
	// "spaces in keys" "spaces in values"
	// "multiple" "equal=signs"
	// "CaSe" "SenSitiVe"
	// "empty value" ""
	// "" "empty key"
	// no error
}
