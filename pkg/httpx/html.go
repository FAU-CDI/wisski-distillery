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

var htmlInternalServerErr = []byte(`<!DOCTYPE HTML><title>Internal Server Error</title>Internal Server Error`)
var htmlNotFound = []byte(`<!DOCTYPE HTML><title>Not Found</title>Not Found`)

// ServeHTTP calls j(r) and returns json
func (h HTMLHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call the function
	result, err := h.Handler(r)

	// entity not found
	if err == ErrNotFound {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		w.Write(htmlNotFound)
		return
	}

	// handle other errors
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(htmlInternalServerErr)
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
