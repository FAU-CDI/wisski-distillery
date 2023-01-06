//go:build !nominify

package httpx

import "io"

// MinifyHTMLWriter wraps the given io.Writer to minify the given html instead.
// The writer must be closed explicitly.
//
// Specific environments may chose to disable http minification, in which case MinifyHTMLWriter becomes the identity function.
func MinifyHTMLWriter(dest io.Writer) io.WriteCloser {
	return noop{Writer: dest}
}

type noop struct {
	io.Writer
}

func (noop) Close() error {
	return nil
}

// MinifyHTML minifies the html source.
// If an error occurs, returns the unmodified source instead.
func MinifyHTML(source []byte) []byte {
	return source
}
