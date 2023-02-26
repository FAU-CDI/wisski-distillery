package httpx

import (
	"html/template"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/timex"
)

const HTMLFlushInterval = time.Second / 10

// WriteHTML writes a html response of type T to w.
// If an error occured, writes an error response instead.
func WriteHTML[T any](result T, err error, template *template.Template, templateName string, w http.ResponseWriter, r *http.Request) (e error) {
	// log any error that occurs;
	defer func() {
		if e != nil {
			zerolog.Ctx(r.Context()).Err(e).Str("path", r.URL.String()).Msg("error rendering template")
		}
	}()

	// create a synced respone writer
	sw := &SyncedResponseWriter{ResponseWriter: w}

	done := make(chan struct{})
	defer close(done)

	// and regularly flush it until the end of the function
	go func() {
		timer := timex.NewTimer()
		defer timex.ReleaseTimer(timer)

		for {
			timer.Reset(HTMLFlushInterval)
			select {
			case <-timer.C:
				sw.Flush()
			case <-done:
				return
			}
		}
	}()

	// intercept any errors
	if HTMLInterceptor.Intercept(sw, r, err) {
		return nil
	}

	// write out the response as html
	sw.Header().Set("Content-Type", "text/html")
	sw.WriteHeader(http.StatusOK)

	// minify html!
	minifier := MinifyHTMLWriter(sw)
	defer minifier.Close()

	// and return the template
	if templateName != "" {
		return template.ExecuteTemplate(minifier, templateName, result)
	} else {
		return template.Execute(minifier, result)
	}
}

type HTMLHandler[T any] struct {
	Handler func(r *http.Request) (T, error)

	Template     *template.Template // called with T
	TemplateName string             // name of template to render, defaults to root
}

// ServeHTTP calls j(r) and returns json
func (h HTMLHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	result, err := h.Handler(r)
	WriteHTML(result, err, h.Template, h.TemplateName, w, r)
}
