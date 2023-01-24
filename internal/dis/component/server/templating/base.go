package templating

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/pools"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
)

//go:embed "src/base.html"
var baseHTML string
var baseTemplate = template.Must(template.New("base.html").Parse(baseHTML))

// Tempalte represents an executable template.
type Template[C any] struct {
	templating *Templating
	p          *Parsed[C]
}

// Template returns a template that, if executed together with the context by the Context method, produces the desired result.
func (tpl *Template[C]) Template() *template.Template {
	return baseTemplate
}

// Context generates the context to pass to an instance of the template returned by Template.
func (tpl *Template[C]) Context(r *http.Request, c C, funcs ...FlagFunc) (ctx *tContext[C]) {
	// create a new context
	ctx = new(tContext[C])

	// setup the basic properties
	ctx.ctx = r.Context()
	ctx.Runtime.RequestURI = r.URL.RequestURI()
	ctx.Runtime.GeneratedAt = time.Now().UTC()
	ctx.Runtime.CSRF = csrf.TemplateField(r)
	ctx.Runtime.Menu = tpl.templating.buildMenu(r)

	// generate the rest of the options
	ctx.Runtime.Flags = ctx.Runtime.Flags.Apply(r, tpl.p.funcs...)
	ctx.Runtime.Flags = ctx.Runtime.Flags.Apply(r, funcs...)

	// if the context has a runtime flags embed, then set the field properly
	if tpl.p.hasRuntimeFlagsEmbed {
		reflect.ValueOf(&c).Elem().
			FieldByName(runtimeFlagsName).
			Set(reflect.ValueOf(ctx.Runtime))
	}

	// the main template
	ctx.cMain = c
	ctx.tMain = tpl.p.tpl

	// the footer template
	ctx.tFooter = tpl.templating.GetCustomizable(footerTemplate)
	ctx.cFooter = ctx.Runtime

	return
}

// ParseForm is like Parse[BaseFormContext]
var ParseForm = Parse[FormContext]

type FormContext struct {
	httpx.FormContext
	RuntimeFlags
}

// NewFormContext returns a new FormContext from an underlying context
func NewFormContext(context httpx.FormContext) FormContext {
	return FormContext{FormContext: context}
}

// FormTemplateContext returns a new handler for a form with the given base context
func FormTemplateContext(tw *Template[FormContext]) func(ctx httpx.FormContext, r *http.Request) any {
	// TODO: Is this needed?
	return func(ctx httpx.FormContext, r *http.Request) any {
		return tw.Context(r, FormContext{FormContext: ctx})
	}
}

// Hander returns a function that returns a context for the given template
func (tw *Template[C]) Handler(f func(r *http.Request) (C, error)) func(r *http.Request) (any, error) {
	// TODO: Should this one be removed?
	return tw.HandlerWithFlags(func(r *http.Request) (C, []FlagFunc, error) {
		c, err := f(r)
		return c, nil, err
	})
}

// HTMLHandler returns a new HTMLHandler for this request
func (tw *Template[C]) HTMLHandler(f func(r *http.Request) (C, error)) httpx.HTMLHandler[any] {
	return httpx.HTMLHandler[any]{
		Handler:  tw.Handler(f),
		Template: tw.Template(),
	}
}

// HandlerWithFlags works like handler, but additionally receive funcs to generate flags
func (tw *Template[C]) HandlerWithFlags(f func(r *http.Request) (C, []FlagFunc, error)) func(r *http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		c, funcs, err := f(r)
		if err != nil {
			return nil, err
		}

		return tw.Context(r, c, funcs...), nil
	}
}

func (tw *Template[C]) HTMLHandlerWithFlags(f func(r *http.Request) (C, []FlagFunc, error)) httpx.HTMLHandler[any] {
	return httpx.HTMLHandler[any]{
		Handler:  tw.HandlerWithFlags(f),
		Template: tw.Template(),
	}
}

// tContext is passed to the underlying template.
//
// Callers may not retain references beyond the invocation of the template.
// Callers must not rely on the internal structure of this tContext.
type tContext[C any] struct {
	Runtime RuntimeFlags // underlying flags

	ctx context.Context // underlying context for render

	// the main template and context
	tMain *template.Template
	cMain C

	// the footer template and context
	tFooter *template.Template
	cFooter RuntimeFlags
}

// Main renders the main template.
func (ctx *tContext[C]) Main() (template.HTML, error) {
	return ctx.renderSafe("main", ctx.tMain, ctx.cMain)
}

// Footer renders the footer template
func (ctx *tContext[C]) Footer() (template.HTML, error) {
	return ctx.renderSafe("footer", ctx.tFooter, ctx.cFooter)
}

const renderSafeError = "Error displaying page. See server log for details. "

func (ctx *tContext[C]) renderSafe(name string, t *template.Template, c any) (template.HTML, error) {

	// already done
	if err := ctx.ctx.Err(); err != nil {
		return "", err
	}

	value, err, panicked := func() (value template.HTML, err error, panicked bool) {
		// get a builder
		builder := pools.GetBuilder()
		defer pools.ReleaseBuilder(builder)

		defer func() {
			if panicked {
				r := recover()
				zerolog.Ctx(ctx.ctx).Error().
					Str("uri", ctx.Runtime.RequestURI).
					Str("name", name).
					Str("panic", fmt.Sprint(r)).
					Msg("templating.Main(): template panic()ed")
			}
		}()

		panicked = true
		err = t.Execute(builder, c)
		panicked = false

		if err != nil {
			zerolog.Ctx(ctx.ctx).Err(err).
				Str("uri", ctx.Runtime.RequestURI).
				Str("name", name).
				Msg("template errored")
		}

		return template.HTML(builder.String()), err, false
	}()

	if err != nil || panicked {
		return renderSafeError, httpx.ErrInternalServerError
	}
	return value, nil
}
