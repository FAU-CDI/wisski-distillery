// Package unpack unpacks files and templates to a target directory.
package unpack

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// ts represents state of the template parser
type ts int

const (
	tsGobble    ts = iota // gobble into dst
	tsSawDollar           // saw a '$'
	tsGobbleVar           // gobble into var
)

// MissingTemplateKeyError indicates [WriteTemplate] found missing keys in the context
type MissingTemplateKeyError struct {
	Keys []string
}

func (mtke MissingTemplateKeyError) Error() string {
	return fmt.Sprintf("missing template keys from context: %v", mtke.Keys)
}

// UnusuedTemplateKeyError indicates [WriteTemplate] found unusued keys in the context
type UnusuedTemplateKeyError struct {
	Keys []string
}

func (utke UnusuedTemplateKeyError) Error() string {
	return fmt.Sprintf("unused template keys from context: %v", utke.Keys)
}

// WriteTemplate writes the template defined by src with the given context into reader.
//
// To run the template, variables of the form ${NAME} are replaced with their corresponding value from the context.
//
// If an underlying read or write fails, it is returned as is.
// Missing template keys return a [MissingTemplateKeyError], but are replaced with the empty string.
// Unused template keys return a [UnusuedTemplateKeyError], but are replaced with the empty string.
//
// Reader / Writer errors are always returned first; next missing template keys, and finally unused template keys.
func WriteTemplate(dst io.Writer, context map[string]string, src io.Reader) error {

	// We keep track of contect keys that have not been used.
	//
	// We first fill the map with all the keys from the context.
	// Then when we use a key, we delete it from the map.
	// If there are any keys left at the end of the replacement, that is an error.
	unusedKeys := make(map[string]struct{}, len(context))
	for key := range context {
		unusedKeys[key] = struct{}{}
	}

	// When we encounter a missing key, put it into this map.
	// This is so that we can build an error message below.
	missingKeys := make(map[string]struct{}, 0)

	// We use a new bufio reader to read data from the input.
	// This is a cheap trick to get a ReadRune() method.
	reader := bufio.NewReader(src)

	//
	// MAIN PARSING LOOP
	//

	// start out in gobble mode!
	mode := tsGobble

	// keep track of variable names
	var varB strings.Builder

parseloop:
	for {
		r, _, err := reader.ReadRune()
		switch {
		case err == io.EOF:
			// finished parsing the source
			break parseloop
		case err != nil:
			// the reader broke
			return err

		case mode == tsGobble && r == '$':
			// saw a '$' in gobble mode
			mode = tsSawDollar
		case mode == tsGobble:
			// normal gobbleing
			// => pass it through
			if _, err := dst.Write([]byte{byte(r)}); err != nil {
				return err
			}

		case mode == tsSawDollar && r == '{':
			// saw '{', following the '$'
			// => read everything else into the buffer
			mode = tsGobbleVar
		case mode == tsSawDollar && r == '$':
			// saw a '$' following the '$'
			// => write the first '$', and handle the case $${stuff}
			if _, err := dst.Write([]byte("$")); err != nil {
				return err
			}
		case mode == tsSawDollar:
			// saw anything else following the '$'
			// => write both back and switch back to gobble mode
			if _, err := dst.Write([]byte{byte('$'), byte(r)}); err != nil {
				return err
			}
			mode = tsGobble

		case mode == tsGobbleVar && r != '}':
			// saw anything except for closing bracket
			// => keep it in the buffer
			if _, err := varB.WriteRune(r); err != nil {
				return err
			}

		case mode == tsGobbleVar:
			// saw a closing '}' inside tsGobbleVar mode
			// => use the variable

			name := varB.String()

			// get the variable from the context
			value, ok := context[name]

			delete(unusedKeys, name) // mark the variable as used!
			if !ok {
				// store unusued variables
				missingKeys[name] = struct{}{}
				value = ""
			}

			// write the replacement into the string
			if _, err := io.WriteString(dst, value); err != nil {
				return err
			}

			// reset the builder and go back into normal mode
			varB.Reset()
			mode = tsGobble
		}
	}

	//
	// CLEANUP UNUSUED INPUT
	//

	switch mode {
	case tsSawDollar:
		// we had a '$', but no '{'
		// => write the trailing '$' into dest
		if _, err := dst.Write([]byte("$")); err != nil {
			return err
		}
	case tsGobbleVar:
		// we had a "${", followed by somthing unclosed
		// => write everything back into the dst
		if _, err := dst.Write([]byte("${")); err != nil {
			return err
		}
		if _, err := io.WriteString(dst, varB.String()); err != nil {
			return err
		}
	}

	// Check if there were missing template keys.
	// If so, we sort them and return an appropriate error.
	if len(missingKeys) != 0 {
		keys := maps.Keys(unusedKeys)
		slices.Sort(keys)
		return MissingTemplateKeyError{
			Keys: keys,
		}
	}

	// Check if there were unused template keys.
	// If so, we sort them and return an appropriate error.
	if len(unusedKeys) != 0 {
		keys := maps.Keys(unusedKeys)
		slices.Sort(keys)
		return UnusuedTemplateKeyError{
			Keys: keys,
		}
	}

	return nil
}

// InstallTemplate unpacks the resource located at src in fsys, then processes it as a template, and eventually writes it to dst.
// Any existing file is truncated and overwritten.
//
// See [WriteTemplate] for possible errors.
func InstallTemplate(env environment.Environment, dst string, context map[string]string, src string, fsys fs.FS) error {

	// open the srcFile
	srcFile, err := fsys.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// stat it
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// check if it is a directory
	if srcInfo.IsDir() {
		return errExpectedFileButGotDirectory
	}

	// open the destination file
	file, err := env.Create(dst, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// write the file!
	return WriteTemplate(file, context, srcFile)
}
