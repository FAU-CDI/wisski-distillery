package httpx

import (
	"io"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/svg"
)

// minifier holds the minfier used for all html minification
//
// NOTE(twiesing): We can't use an init function for this, because otherwise initialization order is incorrect.
var minifier = (func() *minify.M {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	return m
})()

// MinifyHTMLWriter wraps the given io.Writer to minify the given html instead.
// The writer must be closed explicitly.
//
// Specific environments may chose to disable http minification, in which case MinifyHTMLWriter becomes the identity function.
func MinifyHTMLWriter(dest io.Writer) io.WriteCloser {
	return minifier.Writer("text/html", dest)
}

// MinifyHTML minifies the html source.
// If an error occurs, returns the unmodified source instead.
func MinifyHTML(source []byte) []byte {
	result, err := minifier.Bytes("text/html", source)
	if err != nil {
		return source
	}
	return result
}
