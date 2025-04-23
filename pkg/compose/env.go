//spellchecker:words compose
package compose

//spellchecker:words strings github pkglib collection
import (
	"fmt"
	"io"
	"strings"

	"github.com/tkw1536/pkglib/collection"
)

const (
	EnvFileHeader = "# This file is automatically created and updated by the distillery; DO NOT EDIT.\n\n"

	EnvEqualChar   = '='  // assignment
	EnvReplaceChar = '#'  // replacement for invalid characters
	EnvEscapeChar  = '\\' // escaping
	EnvQuoteChar   = '"'  // quoting
)

type invalidNameError string

func (ei invalidNameError) Error() string {
	return fmt.Sprintf("invalid variable name: %q", string(ei))
}

// WriteEnvFile writes a .env file to io.Writer.
// Variables are written in consistent order.
//
// Variable names may only contain ascii letters, numbers or the character "_".
// Invalid variable names are an error.
//
// Variables values are escaped using EscapeEnvValue.
//
// count contains the number of bytes written to writer.
// In case of an error, partial content may already have been written to writer, as indicated by count.
func WriteEnvFile(writer io.Writer, env map[string]string) (count int, err error) {
	var n int

	// write the header to the file
	n, err = fmt.Fprint(writer, EnvFileHeader)
	count += n
	if err != nil {
		return
	}

	for key, value := range collection.IterSorted(env) {
		// if we don't have a valid name, break
		if !isValidVariable(key) {
			return count, invalidNameError(key)
		}

		// write write key = EscapeEnvValue(value) followed by a new line
		n, err = fmt.Fprintf(writer, "%s%s%s\n", key, string(EnvEqualChar), EscapeEnvValue(value))
		if err != nil {
			return count, fmt.Errorf("failed to format variable %q: %w", key, err)
		}
		count += n
	}
	return count, nil
}

// isValidVariable checks if name is a valid variable name.
func isValidVariable(name string) bool {
	for _, r := range name {
		if r != '_' && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// Escape escapes the given value to be written to an environment variable.
// If the value does not need escaping, it may return it unchanged.
//
// EscapeEnvValue allows ASCII characters from ' ' to '~' (inclusive) as well as  '\t', '\r', '\n'.
// Other characters are automatically replaced by EnvReplaceChar.
func EscapeEnvValue(value string) (escaped string) {
	// first check if we need to escape at all.
	var changed bool
	for _, r := range value {
		if !isValidEnvChar(r) || r == '\n' || r == '\r' || r == '\t' || r == '$' || r == EnvEscapeChar || r == EnvQuoteChar {
			changed = true
		}
	}
	if !changed {
		return value
	}

	// make a new builder and make space for the original value
	var builder strings.Builder
	builder.Grow(len(value) + 2)

	// begin the quoting
	builder.WriteRune(EnvQuoteChar)

	// iterate over it
	for _, r := range value {
		// if the character is invalid, it is replaced with an '_'
		if !isValidEnvChar(r) {
			builder.WriteRune(EnvReplaceChar)
			continue
		}

		switch r {
		// custom escape for '\n', '\r', '\t'
		case '\n':
			builder.WriteRune(EnvEscapeChar)
			builder.WriteRune('n')
		case '\r':
			builder.WriteRune(EnvEscapeChar)
			builder.WriteRune('r')
		case '\t':
			builder.WriteRune(EnvEscapeChar)
			builder.WriteRune('t')

		// standard escape for special characters
		case '$', EnvEscapeChar, EnvQuoteChar:
			builder.WriteRune(EnvEscapeChar)
			fallthrough

		// that's it
		default:
			builder.WriteRune(r)
		}
	}

	// close the quote
	builder.WriteRune(EnvQuoteChar)

	return builder.String()
}

// isValidEnvChar checks if the rune r is allowed in environment variables.
func isValidEnvChar(r rune) bool {
	return r == '\n' || r == '\r' || r == '\t' || (r >= ' ' && r <= '~')
}
