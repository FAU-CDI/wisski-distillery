// Package envreader
package envreader

import (
	"bufio"
	"io"
	"strings"
)

// Scanner is a scanner for environment files.
// To create a new scanner use [NewScanner].
//
// It scans through a reader and reads environment variables from it.
// Reads may be internally buffered.
//
// An environment variable is of the form:
//   KEY=VALUE
// on a separate line.
// Keys and values are case-sensitive and may contain anything except for newline characters.
// Spaces around key and value are trimmed using [strings.TrimSpace].
// Keys may not contain an '='.
// Lines not containing a '=' (e.g. blank lines) and those starting with '#' and '//' are ignored.
//
// To advance the scanner to the next key, value pair use [Scan].
// To get the current (key, value) pair, use [Data].
//
// A typical use-case of a scanner is as follows:
//
//  scanner := NewScanner(r)
//  for scanner.Scan() {
//      // process any data ....
//      fmt.Println(scanner.Data())
//  }
//  if err := scanner.Err(); err != nil {
//    	// handle errors
//  }
//
// For the common use case of reading a set of distinct keys from a file see [ReadAll].
type Scanner struct {
	s *bufio.Scanner

	// current key and value
	key   string
	value string
}

// NewScanner creates a new scanner from the underlying Reader
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		s: bufio.NewScanner(r),
	}
}

// Scanner advances the scanner until the next KEY=VALUE pair.
//
// If there are no more values left (e.g. the underlying reader returned io.EOF)
// or when an unexpected error occured, returns false.
//
// A caller should always check Err() to see if there was an error.
func (scanner *Scanner) Scan() bool {
	var found bool
	for scanner.s.Scan() {
		// check that we don't have an empty or comment only line
		tokens := strings.TrimSpace(scanner.s.Text())
		if len(tokens) == 0 || tokens[0] == '#' || strings.HasPrefix(tokens, "//") {
			continue
		}

		// check that we have a 'key=value' pair
		scanner.key, scanner.value, found = strings.Cut(tokens, "=")
		if !found {
			continue
		}

		// got a key = value
		scanner.key = strings.TrimSpace(scanner.key)
		scanner.value = strings.TrimSpace(scanner.value)
		return true
	}

	// nothing found
	scanner.key = ""
	scanner.value = ""
	return false
}

// Data reads the current value from the scanner.
// When Scan() has not been called, or returned false, returns two empty strings.
func (scanner Scanner) Data() (key, value string) {
	return scanner.key, scanner.value
}

// Err returns any error that occured on the underlying read.
//
// When no error occured, or the underlying read is io.EOF, returns nil.
func (scanner Scanner) Err() error {
	return scanner.s.Err()
}

// ReadAll creates a new [Scanner], and then reads all key/value pairs from r.
// If a key occurs more than once, only the last value is set in the returned map.
func ReadAll(r io.Reader) (values map[string]string, err error) {
	scanner := NewScanner(r)

	// read and store all values
	values = make(map[string]string)
	for scanner.Scan() {
		key, value := scanner.Data()
		values[key] = value
	}

	// check if there was an error!
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
