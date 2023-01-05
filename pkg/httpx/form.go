package httpx

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
)

// DefaultFieldTemplate is the default template to render fields.
var DefaultFieldTemplate = template.Must(template.New("").Parse(`<input type="{{.Type}}" value="{{.Value}}" name="{{.Name}}" placeholder={{.Placeholder}}>`))
var PureCSSFieldTemplate = template.Must(template.New("").Parse(`
<div class="pure-control-group"><label for="{{.Name}}">{{.Label}}</label><input type="{{.Type}}" value="{{.Value}}" name="{{.Name}}" id="{{.Name}}" placeholder="{{.Placeholder}}"></div>`))

// Form provides a form that a user can submit via a http POST method call.
// It implements [http.Handler].
type Form[D any] struct {
	// Fields are the fields this form consists of.
	Fields []Field

	// FieldTemplate is an optional template to be executed for each field.
	// FieldTemplate may be nil; in which case [DefaultFieldTemplate] is used.
	FieldTemplate *template.Template

	// SkipCSRF if CSRF should be explicitly omitted
	SkipCSRF bool

	// SkipForm, if non-nil, is called on every get request to determine if form parsing should be skipped entirely.
	// If skip is true, RenderSuccess is directly called with the given values map.
	SkipForm func(r *http.Request) (data D, skip bool)

	// RenderForm handles rendering a form into a request.
	// If RenderForm is nil, RenderTemplate is invoked with an appropriate [FormContext] instance.
	// Either RenderForm or RenderTemplate must be non-nil.
	//
	// template holds pre-rendered html fields.
	// err is a non-nil error returned from Validate, or the r.ParseForm() method.
	// It is nil on the initial render.
	RenderForm func(context FormContext, w http.ResponseWriter, r *http.Request)

	// RenderTemplate represents an optional form to display to the user when RenderForm is nil
	// It is passed the return value of [RenderTemplateContext], or a [FormContext] instance if this does not exist.
	RenderTemplate *template.Template

	// RenderTemplateContext is the context to be used for RenderTemplate.
	// When nil, assumed to be the identify function
	RenderTemplateContext func(ctx FormContext, r *http.Request) any

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

// ServeHTTP implements [http.Handler] and serves the form
func (form *Form[D]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// renderForm renders the form into a request
func (form *Form[D]) renderForm(err error, values map[string]string, w http.ResponseWriter, r *http.Request) {
	template := form.Template(values, err != nil)
	if !form.SkipCSRF {
		template += csrf.TemplateField(r)
	}

	ctx := FormContext{Err: err, Form: template}

	if form.RenderForm != nil {
		form.RenderForm(ctx, w, r)
		return
	}

	// must have a form or a RenderForm
	if form.RenderTemplate == nil {
		panic("form.RenderForm and form.Form are nil")
	}

	// get the template context
	var tplctx any
	if form.RenderTemplateContext == nil {
		tplctx = ctx
	} else {
		tplctx = form.RenderTemplateContext(ctx, r)
	}

	// render the form
	WriteHTML(tplctx, nil, form.RenderTemplate, "", w, r)
}

// FormContext is passed to Form.Form when used
type FormContext struct {
	// Error is the underlying error (if any)
	Err error

	// Template is the underlying template rendered as html
	Form template.HTML
}

// Error returns the underlying error string
func (fc FormContext) Error() string {
	if fc.Err == nil {
		return ""
	}
	return fc.Err.Error()
}

// renderSuccess renders a successfull pass of the form
// if an error occurs during rendering, renderForm is called instead
func (form *Form[D]) renderSuccess(data D, values map[string]string, w http.ResponseWriter, r *http.Request) {
	err := form.RenderSuccess(data, values, w, r)
	if err == nil {
		return
	}
	form.renderForm(err, values, w, r)
}

// Field represents a field inside a form.
type Field struct {
	Name string    // Name is the name of the field
	Type InputType // Type is the type of the field. It corresponds to the "name" attribute in html.

	Placeholder string // Value for the "placeholder" attribute
	Label       string // (External) Label for the field. Not used by the default template.

	EmptyOnError bool // indicates if the field should be reset on error
}

// fieldContext is passed to the template context
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
	CheckboxField InputType = "checkbox"
)
