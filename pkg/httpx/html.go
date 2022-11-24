package httpx

import (
	"html/template"
	"net/http"
)

type HTMLHandler[T any] struct {
	Handler func(r *http.Request) (T, error)

	Template     *template.Template // called with T
	TemplateName string             // name of template to render, defaults to root
}

// ServeHTTP calls j(r) and returns json
func (h HTMLHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	result, err := h.Handler(r)

	// intercept any errors
	if htmlInterceptor.Intercept(w, r, err) {
		return
	}

	// write out the response as json
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if h.TemplateName != "" {
		h.Template.ExecuteTemplate(w, h.TemplateName, result)
	} else {
		h.Template.Execute(w, result)
	}
}
