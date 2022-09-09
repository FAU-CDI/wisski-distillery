// Package unpack unpacks files and templates to a target directory.
package unpack

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var errExpectedFileButGotDirectory = errors.New("Expected a file, but got a directory")

// UnpackTemplate unpacks the given file template and template
func UnpackTemplate(context map[string]string, src fs.File) ([]byte, fs.FileMode, error) {
	// stat the source file to install
	srcStat, srcErr := src.Stat()
	if srcErr != nil {
		return nil, 0, errors.Wrapf(srcErr, "Error calling stat on source")
	}

	// should not be a directory
	if srcStat.IsDir() {
		return nil, 0, errors.Wrapf(errExpectedFileButGotDirectory, "Error calling stat on source %s", srcStat.Name())
	}

	// read all the bytes into a buffer
	var buffer bytes.Buffer
	WriteTemplate(&buffer, context, src)
	return buffer.Bytes(), srcStat.Mode(), nil
}

type templateMode int

const (
	templateModeNormal templateMode = iota // normal mode
	templateModeDollar                     // saw '$'
	templateModeOpen                       // saw '${'
)

// WriteTemplate writes the template defined by src with the given context into reader.
//
// To run the template, variables of the form ${NAME} are replaced with their corresponding value from the context.
//
// Extra or missing variables from the context are an error.
func WriteTemplate(dst io.Writer, context map[string]string, src io.Reader) error {
	// keep track of context keys that have not been used
	unuusedContext := make(map[string]struct{}, len(context))
	for key := range context {
		unuusedContext[key] = struct{}{}
	}

	reader := bufio.NewReader(src) // a new reader
	var missingKeyErr error        // error for missing keys
	var builder strings.Builder    // holding variable names
	mode := templateModeNormal     // the current mode of the reader
parseloop:
	for {
		r, _, err := reader.ReadRune()
		switch {
		case err == io.EOF:
			/* finished the source, see below */
			break parseloop
		case err != nil:
			/* something went wrong */
			return err

		case mode == templateModeNormal && r == '$':
			// saw a '$' in normal mode
			// => switch to the dollar mode
			mode = templateModeDollar
		case mode == templateModeNormal:
			// saw anything else
			// => just pass it through
			if _, err := dst.Write([]byte{byte(r)}); err != nil {
				return err
			}

		case mode == templateModeDollar && r == '{':
			// saw '{', following the '$'
			// => read everything else into the buffer
			mode = templateModeOpen
		case mode == templateModeDollar && r == '$':
			// saw a '$' following the '$'
			// => write the first '$', and handle the case $${stuff}
			if _, err := dst.Write([]byte("$")); err != nil {
				return err
			}
		case mode == templateModeDollar:
			// saw anything else following the '$'
			// => write both back and switch back to normal mode
			if _, err := dst.Write([]byte{byte('$'), byte(r)}); err != nil {
				return err
			}
			mode = templateModeNormal

		case mode == templateModeOpen && r != '}':
			// saw anything except for closing bracket
			// => keep it in the buffer
			if _, err := builder.WriteRune(r); err != nil {
				return err
			}

		case mode == templateModeOpen:
			// saw a closing '}' inside the open mode
			// => use the variable

			name := builder.String()
			delete(unuusedContext, name) // mark the variable as used

			// get the variable from the context
			value, ok := context[name]
			if missingKeyErr != nil && !ok {
				missingKeyErr = errors.Errorf("key %s missing in context", name)
			}

			// write the replacement into the string
			if _, err := io.WriteString(dst, value); err != nil {
				return err
			}

			// reset the builder and go back into normal mode
			builder.Reset()
			mode = templateModeNormal

		default:
			panic("never reached")
		}
	}

	// cleanup at end of input

	switch mode {
	case templateModeNormal:
		// => everything is fine
	case templateModeDollar:
		// we had a '$', but no '{'
		// => write the trailing '$' into dest
		if _, err := dst.Write([]byte("$")); err != nil {
			return err
		}
	case templateModeOpen:
		// we had a "${", followed by somthing unclosed
		// => write everything back into the dst
		if _, err := dst.Write([]byte("${")); err != nil {
			return err
		}
		if _, err := io.WriteString(dst, builder.String()); err != nil {
			return err
		}
	default:
		panic("never reached")
	}

	// check if there was a missing key!
	if missingKeyErr != nil {
		return missingKeyErr
	}

	// check if there was an unused key!
	if len(unuusedContext) != 0 {
		keys := maps.Keys(unuusedContext)
		slices.Sort(keys)
		return errors.Errorf("additional keys %s in context", strings.Join(keys, ","))
	}

	return nil
}
