package httpx

import (
	"html/template"
	"net/http"
)

// WriteHTML writes a html response of type T to w.
// If an error occured, writes an error response instead.
func WriteHTML[T any](result T, err error, template *template.Template, templateName string, w http.ResponseWriter, r *http.Request) {
	// intercept any errors
	if HTMLInterceptor.Intercept(w, r, err) {
		return
	}

	// write out the response as html
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// minify html!
	minifier := MinifyHTMLWriter(w)
	defer minifier.Close()

	// and return the template
	if templateName != "" {
		template.ExecuteTemplate(minifier, templateName, result)
	} else {
		template.Execute(minifier, result)
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
