package httpx

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/csrf"
)

// DefaultFieldTemplate is the default template to render fields.
var DefaultFieldTemplate = template.Must(template.New("").Parse(`<input type="{{.Type}}" value="{{.Value}}" name="{{.Name}}" placeholder={{.Placeholder}}>`))
var PureCSSFieldTemplate = template.Must(template.New("").Parse(`
<div class="pure-control-group"><label for="{{.Name}}">{{.Label}}</label><input type="{{.Type}}" value="{{.Value}}" name="{{.Name}}" id="{{.Name}}" placeholder="{{.Placeholder}}"></div>`))

// Form implements a user-submittable form
type Form[D any] struct {
	Fields []Field

	// FieldTemplate is executed for each field.
	// Defaults to a [DefaultFieldTemplate]
	FieldTemplate *template.Template

	// CSRF holds an optional reference to a CSRF.Protect call.
	// It must be set before any other functions on this Form are called, and may not be changed.
	CSRF func(http.Handler) http.Handler
	csrf lazy.Lazy[http.Handler]

	// SkipForm, if non-nil, is called on every get request to determine if form parsing should be skipped entirely.
	// If skip is true, RenderSuccess is directly called with the given values map.
	SkipForm func(r *http.Request) (data D, skip bool)

	// RenderForm handles rendering a form into a request.
	//
	// template holds pre-rendered html fields.
	// err is a non-nil error returned from Validate, or the r.ParseForm() method.
	// It is nil on the initial render.
	RenderForm func(template template.HTML, err error, w http.ResponseWriter, r *http.Request)

	// Validate, if non-nil, validates the given submitted values.
	// There is no guarantee that the values are set.
	Validate func(r *http.Request, values map[string]string) (D, error)

	// RenderSuccess handles rendering a success result into a response.
	RenderSuccess func(data D, values map[string]string, w http.ResponseWriter, r *http.Request) error
}

// Template renders this form as a HTML string for insertion into a template.
func (form *Form[D]) Template(values map[string]string, isError bool) template.HTML {
	var builder strings.Builder

	for _, field := range form.Fields {
		value := values[field.Name]
		if isError && field.EmptyOnError {
			value = ""
		}

		field.WriteTo(&builder, form.FieldTemplate, value)
	}

	return template.HTML(builder.String())
}

// Values returns (validated) form values contained in the given request
func (form *Form[D]) Values(r *http.Request) (v map[string]string, d D, err error) {
	// parse the form
	if err := r.ParseForm(); err != nil {
		return nil, d, err
	}

	// pick each of the values
	values := make(map[string]string, len(form.Fields))
	for _, field := range form.Fields {
		values[field.Name] = r.PostForm.Get(field.Name)
	}

	// validate the form
	if form.Validate != nil {
		d, err = form.Validate(r, values)
		if err != nil {
			return nil, d, err
		}
	}

	// and return them
	return values, d, nil
}

func (form *Form[D]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := form.csrf.Get(func() (handler http.Handler) {
		handler = http.HandlerFunc(form.serveHTTP)
		if form.CSRF != nil {
			handler = form.CSRF(handler)
		}
		return
	})
	handler.ServeHTTP(w, r)
}

func (form *Form[D]) serveHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		TextInterceptor.Intercept(w, r, ErrMethodNotAllowed)
		return
	case r.Method == http.MethodPost:
		values, data, err := form.Values(r)
		if err != nil {
			form.renderForm(err, values, w, r)
		} else {
			form.renderSuccess(data, values, w, r)
		}
	case r.Method == http.MethodGet && form.SkipForm != nil:
		if data, skip := form.SkipForm(r); skip {
			form.renderSuccess(data, nil, w, r)
			return
		}
		fallthrough
	case r.Method == http.MethodGet:
		form.renderForm(nil, nil, w, r)
	}
}

func (form *Form[D]) renderForm(err error, values map[string]string, w http.ResponseWriter, r *http.Request) {
	template := form.Template(values, err != nil)
	if form.CSRF != nil {
		template += csrf.TemplateField(r)
	}
	form.RenderForm(template, err, w, r)
}

func (form *Form[D]) renderSuccess(data D, values map[string]string, w http.ResponseWriter, r *http.Request) {
	err := form.RenderSuccess(data, values, w, r)
	if err == nil {
		return
	}
	form.renderForm(err, values, w, r)
}

// Field represents a field
type Field struct {
	Name string
	Type InputType

	Placeholder string // Optional placeholder
	Label       string // Label for the template. Not used by the default template.

	EmptyOnError bool // indicates if the field should be reset on error
}

type fieldContext struct {
	Field
	Value string
}

func (field Field) WriteTo(w io.Writer, template *template.Template, value string) {
	if template == nil {
		template = DefaultFieldTemplate
	}
	template.Execute(w, fieldContext{Field: field, Value: value})
}

// InputType represents the type of input
type InputType string

const (
	TextField     InputType = "text"
	PasswordField InputType = "password"
)
