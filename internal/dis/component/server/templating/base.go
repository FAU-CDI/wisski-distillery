package templating

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx/content"
	"github.com/tkw1536/pkglib/httpx/form"
	"github.com/tkw1536/pkglib/httpx/wrap"
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
	ctx = tpl.context(r, funcs...)
	ctx.cMain = c
	return ctx
}

func (tpl *Template[C]) context(r *http.Request, funcs ...FlagFunc) (ctx *tContext[C]) {
	// create a new context
	ctx = new(tContext[C])

	// setup the basic properties
	ctx.ctx = r.Context()
	ctx.Runtime.RequestURI = r.URL.RequestURI()
	ctx.Runtime.StartedAt = wrap.TimeStart(r).UTC()
	ctx.Runtime.GeneratedAt = time.Now().UTC()
	ctx.Runtime.CSRF = csrf.TemplateField(r)
	ctx.Runtime.Menu = tpl.templating.buildMenu(r)

	// generate the rest of the options
	ctx.Runtime.Flags = ctx.Runtime.Flags.Apply(r, tpl.p.funcs...)
	ctx.Runtime.Flags = ctx.Runtime.Flags.Apply(r, funcs...)
	ctx.updateEmbedded = tpl.p.hasRuntimeFlagsEmbed

	// the main template
	ctx.tMain = tpl.p.tpl

	// the footer template
	ctx.tFooter = tpl.templating.GetCustomizable(footerTemplate)
	ctx.cFooter = ctx.Runtime

	return
}

// ParseForm is like Parse[BaseFormContext]
var ParseForm = Parse[FormContext]

type FormContext struct {
	form.FormContext
	RuntimeFlags
}

// NewFormContext returns a new FormContext from an underlying context
func NewFormContext(context form.FormContext) FormContext {
	return FormContext{FormContext: context}
}

// FormTemplateContext returns a new handler for a form with the given base context
func FormTemplateContext(tw *Template[FormContext]) func(ctx form.FormContext, r *http.Request) any {
	// TODO: Is this needed?
	return func(ctx form.FormContext, r *http.Request) any {
		return tw.Context(r, FormContext{FormContext: ctx})
	}
}

// HandlerWithFlags returns a function that, given a request, generates context and error to pass to the generated template.
// The worker implements the actual buisness logic, it takes a request, and returns the content for the main template, and any error.
// See also HandlerWithFlags.
func (tw *Template[C]) Handler(f func(r *http.Request) (C, error)) func(r *http.Request) (any, error) {
	return tw.HandlerWithFlags(func(r *http.Request) (C, []FlagFunc, error) {
		c, err := f(r)
		return c, nil, err
	})
}

// HandlerWithFlags returns a function that, given a request, generates context and error to pass to the generated template.
// The worker implements the actual buisness logic, it takes a request, and returns the content for the main template, flag functions and error.
// See also Handler.
func (tw *Template[C]) HandlerWithFlags(worker func(r *http.Request) (C, []FlagFunc, error)) func(r *http.Request) (any, error) {
	return func(r *http.Request) (any, error) {
		c, funcs, err := worker(r)
		if err != nil {
			return nil, err
		}

		return tw.Context(r, c, funcs...), nil
	}
}

// HTMLHandler creates a new httpx.HTMLHandler that calls tw.Handler(worker) and tw.Template.
// See also Handler.
func (tw *Template[C]) HTMLHandler(handling *handling.Handling, worker func(r *http.Request) (C, error)) content.HTMLHandler[any] {
	return content.HTMLHandler[any]{
		Handler:     tw.Handler(worker),
		Template:    tw.Template(),
		Interceptor: handling.HTMLInterceptor(),
	}
}

// HTMLHandlerWithFlags creates a new httpx.HTMLHandler that calls tw.HandlerWithFlags(worker) and tw.Template.
// See also HandlerWithFlags.
func (tw *Template[C]) HTMLHandlerWithFlags(handling *handling.Handling, worker func(r *http.Request) (C, []FlagFunc, error)) content.HTMLHandler[any] {
	return content.HTMLHandler[any]{
		Handler:     tw.HandlerWithFlags(worker),
		Template:    tw.Template(),
		Interceptor: handling.HTMLInterceptor(),
	}
}

// tContext is passed to the underlying template.
//
// Callers may not retain references beyond the invocation of the template.
// Callers must not rely on the internal structure of this tContext.
type tContext[C any] struct {
	Runtime        RuntimeFlags // underlying flags
	updateEmbedded bool         // should we automatically update an embedded RuntimeFlags inside the context?

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
	// if the context has a runtime flags embed, then set the field properly
	if ctx.updateEmbedded {
		reflect.ValueOf(&ctx.cMain).Elem().
			FieldByName(runtimeFlagsName).
			Set(reflect.ValueOf(ctx.Runtime))
	}

	return ctx.renderSafe("main", ctx.tMain, ctx.cMain)
}

// Footer renders the footer template
func (ctx *tContext[C]) Footer() (template.HTML, error) {
	return ctx.renderSafe("footer", ctx.tFooter, ctx.cFooter)
}

const renderSafeError = "Error displaying page. See server log for details. "
const renderPanicError = "Panic displaying page. See server log for details. "

func (ctx *tContext[C]) renderSafe(name string, t *template.Template, c any) (template.HTML, error) {

	// already done with context => return
	if err := ctx.ctx.Err(); err != nil {
		return "", err
	}

	value, panicked, panik, stack, err := func() (value template.HTML, panicked bool, panik any, stack []byte, err error) {
		var builder strings.Builder

		defer func() {
			if panicked {
				panik = recover()
				stack = debug.Stack()

				zerolog.Ctx(ctx.ctx).Error().
					Str("uri", ctx.Runtime.RequestURI).
					Str("name", name).
					Str("panic", fmt.Sprint(panik)).
					Str("stack", string(stack)).
					Msg("renderSafe: template panic()ed")
			}
		}()

		panicked = true
		err = t.Execute(&builder, c)
		panicked = false

		if err != nil {
			zerolog.Ctx(ctx.ctx).Err(err).
				Str("uri", ctx.Runtime.RequestURI).
				Str("name", name).
				Msg("template errored")
		}

		return template.HTML(builder.String()), false, nil, nil, err
	}()

	if err != nil {
		return renderSafeError, err
	}
	if panicked {
		return renderPanicError, panicErr{value: panik, stack: stack}
	}
	return value, nil
}

// panicErr is returned by renderSafe when a panic occurs
type panicErr struct {
	value any
	stack []byte
}

func (pe panicErr) Error() string {
	return fmt.Sprintf("panic: %v", pe.value)
}
