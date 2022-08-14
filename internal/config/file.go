package config

import (
	"bufio"
	"io"
	"strings"
)

// Scanner scans an io.Reader for a source file
type Scanner struct {
	src *bufio.Scanner

	key   string
	value string
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		src: bufio.NewScanner(r),
	}
}

// Scanner advances the scanner to the next variable
func (scanner *Scanner) Scan() bool {
	for scanner.src.Scan() {
		// check that we don't have an empty or comment only line
		tokens := strings.TrimSpace(scanner.src.Text())
		if len(tokens) == 0 || tokens[0] == '#' || strings.HasPrefix(tokens, "//") {
			continue
		}

		// check that we have a 'key=value' pair
		values := strings.SplitN(tokens, "=", 2)
		if len(values) != 2 {
			continue
		}

		// got a key = value
		scanner.key = strings.TrimSpace(values[0])
		scanner.value = strings.TrimSpace(values[1])
		return true
	}
	scanner.key = ""
	scanner.value = ""
	return false
}

// Data reads the current value from the scanner.
// When Scan() has not been called, or returned false, returns two empty strings.
func (scanner Scanner) Data() (key, value string) {
	return scanner.key, scanner.value
}

// Error returns an error (if any)
func (scanner Scanner) Error() error {
	return scanner.src.Err()
}

// ReadAll reads all key-value pairs from r.
// If a key occurs more than once, a later occurance overwrites a previous one.
func ReadAll(r io.Reader) (values map[string]string, err error) {
	scanner := NewScanner(r)

	// read and store all values
	values = make(map[string]string)
	for scanner.Scan() {
		key, value := scanner.Data()
		values[key] = value
	}

	// check if there was an error!
	if err := scanner.Error(); err != nil {
		return nil, err
	}
	return values, nil
}
