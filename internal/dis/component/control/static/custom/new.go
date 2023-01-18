package custom

import (
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
)

// Parsed represents a parsed template that receives an underlying context of type C
type Parsed[C any] struct {
	template *template.Template
}

// Parse creates a new Parsed from a template source.
// Parse calls panic() when parsing fails.
func Parse[C any](name string, source []byte, Assets static.Assets) Parsed[C] {
	return Parsed[C]{
		template: Assets.MustParseShared(name, string(source)),
	}
}

// Prepare prepares this template for use inside a concrete handler.
// gaps must either be of length 0 or length 1 and may pre-fill gaps to be used when executing the template later.
func (p *Parsed[C]) Prepare(custom *Custom, gaps ...BaseContextGaps) *Template[C] {
	wrap := Template[C]{
		custom:   custom,
		template: custom.Template(p.template),
	}
	if len(gaps) > 1 {
		panic("WrapTemplate: must provide either 1 or no gaps")
	}
	if len(gaps) == 1 {
		wrap.gaps = gaps[0]
	}
	return &wrap
}

// Tempalte represents an executable template.
type Template[C any] struct {
	custom   *Custom
	template *template.Template
	gaps     BaseContextGaps
}

// Template returns a template that, if executed together with the context by the Context method, produces the desired result.
func (tw *Template[C]) Template() *template.Template {
	return tw.template
}

// Context generates a context for a given request that can be used to execute the provided template.
func (tw *Template[C]) Context(r *http.Request, c C, gaps ...BaseContextGaps) any {
	// make the gaps something
	if len(gaps) > 1 {
		panic("Context: must provide either 1 or no gaps")
	}

	// update the context with gaps
	{
		g := tw.gaps
		if len(gaps) == 1 {
			g = gaps[0]
		}
		tw.custom.update(&c, r, g)
	}

	return c
}

// ParseForm is like Parse[BaseFormContext]
var ParseForm = Parse[BaseFormContext]

// FormTemplateContext returns a new handler for a form with the given base context
func FormTemplateContext(tw *Template[BaseFormContext]) func(ctx httpx.FormContext, r *http.Request) any {
	return func(ctx httpx.FormContext, r *http.Request) any {
		return tw.Context(r, BaseFormContext{FormContext: ctx})
	}
}

// MappedHandler returns a new handler that maps the incoming context via f
func MappedHandler[In, Out any](tw *Template[Out], f func(ctx In, r *http.Request) (Out, BaseContextGaps)) func(ctx In, r *http.Request) any {
	// TODO: Should this one be removed?
	return func(ctx In, r *http.Request) any {
		c, g := f(ctx, r)
		return tw.Context(r, c, g)
	}
}

// Hander returns a function that returns a context for the given template
func (tw *Template[C]) Handler(f func(r *http.Request) (C, error)) func(r *http.Request) (any, error) {
	// TODO: Should this one be removed?
	return tw.HandlerWithGaps(func(r *http.Request, gaps *BaseContextGaps) (C, error) {
		return f(r)
	})
}

// HTMLHandler returns a new HTMLHandler for this request
func (tw *Template[C]) HTMLHandler(f func(r *http.Request) (C, error)) httpx.HTMLHandler[any] {
	return httpx.HTMLHandler[any]{
		Handler:  tw.Handler(f),
		Template: tw.Template(),
	}
}

// HandlerWithGaps works like handler, but additionally receives a gaps object to update.
func (tw *Template[C]) HandlerWithGaps(f func(r *http.Request, gaps *BaseContextGaps) (C, error)) func(r *http.Request) (any, error) {
	// TODO: Drop this variant?
	var zero C
	return func(r *http.Request) (any, error) {
		g := tw.gaps.clone()
		c, err := f(r, &g)
		if err != nil {
			return zero, err
		}

		// update the context
		return tw.Context(r, c, g), nil
	}
}

func (tw *Template[C]) HTMLHandlerWithGaps(f func(r *http.Request, gaps *BaseContextGaps) (C, error)) httpx.HTMLHandler[any] {
	return httpx.HTMLHandler[any]{
		Handler:  tw.HandlerWithGaps(f),
		Template: tw.Template(),
	}
}

// Execute executes this template with the given context
func (tw *Template[C]) Execute(w http.ResponseWriter, r *http.Request, c C, gaps ...BaseContextGaps) error {
	return tw.ExecuteWithError(w, r, c, nil, gaps...)
}

// ExecuteWithError executes this template, or the default error handler if err != nil
func (tw *Template[C]) ExecuteWithError(w http.ResponseWriter, r *http.Request, c C, err error, gaps ...BaseContextGaps) error {
	// TODO: Drop this variant?
	// TODO: This should be removed!
	return httpx.WriteHTML(tw.Context(r, c, gaps...), err, tw.template, "", w, r)
}
