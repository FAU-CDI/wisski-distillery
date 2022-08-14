package distillery

import (
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var templateRegexp = regexp.MustCompile(`\${[^}]+}`)

// InstallTemplates open the resource src, and installs it into dst.
// the template resource must fit into memory.
//
// For each variable ${THING} inside dest, a key 'THING' must exist in context.
// Extra or missing template keys are an error.
func InstallTemplate(dst, src string, context map[string]string) error {
	bytes, srcMode, err := doTemplate(src, context)
	if err != nil {
		return err
	}

	// determine if we need to create the destination file, or if it already exists
	dstStat, dstErr := os.Stat(dst)
	flag := os.O_WRONLY
	switch {
	case os.IsNotExist(dstErr):
		flag |= os.O_CREATE
	case dstErr != nil:
		return errors.Wrapf(dstErr, "Error calling stat on destination %s", dst)
	case dstStat.IsDir():
		return errors.Wrapf(errExpectedFileButGotDirectory, "Error processing destination %s", dst)
	}

	// open and write the destination file
	dstFile, err := os.OpenFile(dst, flag, srcMode)
	if err != nil {
		return errors.Wrapf(err, "Unable to open file %s", dst)
	}
	_, err = dstFile.Write(bytes)
	return errors.Wrapf(err, "Unable to write destination %s", dst)
}

// ReadTemplate is like InstallTemplate, except that it writes template into a byte slice and returns it.
func ReadTemplate(src string, context map[string]string) ([]byte, error) {
	bytes, _, err := doTemplate(src, context)
	return bytes, err
}

func doTemplate(src string, context map[string]string) (bytes []byte, mode fs.FileMode, err error) {
	// open the source file!
	srcFile, err := resourceEmbed.Open(src)
	if err != nil {
		return nil, mode, errors.Wrapf(err, "Error opening source file %s", src)
	}
	defer srcFile.Close()

	// stat the source file to install
	srcStat, srcErr := srcFile.Stat()
	if srcErr != nil {
		return nil, mode, errors.Wrapf(srcErr, "Error calling stat on source %s", src)
	}

	// should not be a directory
	if srcStat.IsDir() {
		return nil, mode, errors.Wrapf(errExpectedFileButGotDirectory, "Error calling stat on source %s", src)
	}

	// read the template and replace
	templates, err := io.ReadAll(srcFile)
	if err != nil {
		return nil, mode, errors.Wrapf(err, "Unable to read src file %s", src)
	}

	// keep track of context keys that have not been used
	unuusedContext := make(map[string]struct{}, len(context))
	for key := range context {
		unuusedContext[key] = struct{}{}
	}

	// replace the template regexp
	// keeping track of unuused errors
	var hadError error
	templates = templateRegexp.ReplaceAllFunc(templates, func(b []byte) []byte {
		name := string(b[2 : len(b)-1]) // remove the leading ${ and trailing }
		delete(unuusedContext, name)    // mark the key as having been read

		value, ok := context[name]
		if hadError != nil && !ok {
			hadError = errors.Errorf("key %s missing in context", name)
		}
		return []byte(value)
	})

	if hadError != nil {
		return nil, mode, hadError
	}

	if len(unuusedContext) != 0 {
		keys := maps.Keys(unuusedContext)
		slices.Sort(keys)
		return nil, mode, errors.Errorf("additional keys %s in context", strings.Join(keys, ","))
	}

	// return the data and the mode!
	return templates, srcStat.Mode(), nil
}
