package fsx

import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// SameFile checks if path1 and path2 refer to the same file.
// If both files exist, they are compared using [env.SameFile].
// If both files do not exist, the paths are first compared syntactically and then via recursion on [filepath.Dir].
func SameFile(env environment.Environment, path1, path2 string) bool {

	// initial attempt: check if directly
	same, certain := couldBeSameFile(env, path1, path2)
	if certain {
		return same
	}

	// second attempt: find the directory names and base paths
	d1, n1 := filepath.Dir(path1), filepath.Base(path1)
	d2, n2 := filepath.Dir(path2), filepath.Base(path2)

	// if we have different file names (and they don't exist)
	// we don't need to continue
	if n1 != n2 {
		return false
	}

	// compare the base names!
	{
		same, _ := couldBeSameFile(env, d1, d2)
		return same
	}
}

// couldBeSameFile checks if path1 might be the same as path2.
//
// If both files exist, compares using [env.SameFile].
// Otherwise compares absolute paths using string comparison.
//
// same indicates if they might be the same file.
// authorative indiciates if the result is authorative.
func couldBeSameFile(env environment.Environment, path1, path2 string) (same, authorative bool) {
	{
		// stat both files
		info1, err1 := os.Stat(path1)
		info2, err2 := os.Stat(path2)

		// both files exist => check using env.SameFile
		// the result is always authorative
		if err1 == nil && err2 == nil {
			same = os.SameFile(info1, info2)
			authorative = true
			return
		}

		// only 1 file errored => they could be different
		if (err1 == nil) != (err2 == nil) {
			return
		}

		// only 1 file does not exist => they could be different
		if os.IsNotExist(err1) != os.IsNotExist(err2) {
			return
		}
	}

	{
		// resolve paths absolutely
		rpath1, err1 := filepath.Abs(path1)
		rpath2, err2 := filepath.Abs(path2)

		// if either path could not be resolved absolutely
		// fallback to just using clean!
		if err1 != nil {
			rpath1 = filepath.Clean(path1)
		}
		if err2 != nil {
			rpath2 = filepath.Clean(path2)
		}

		// compare using strings
		same = rpath1 == rpath2
		authorative = same // positive result is authorative!
		return
	}
}
