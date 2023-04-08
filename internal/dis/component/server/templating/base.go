package templating

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/timex"
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

// LazyContext is like the lazy context
func (tpl *Template[C]) LazyContext(r *http.Request, f func() (C, error), funcs ...FlagFunc) (ctx *tContext[C]) {
	ctx = tpl.context(r, funcs...)
	ctx.startLazy(f)
	return ctx
}

// Context generates the context to pass to an instance of the template returned by Template.
func (tpl *Template[C]) Context(r *http.Request, c C, funcs ...FlagFunc) (ctx *tContext[C]) {
	ctx = tpl.context(r, funcs...)
	ctx.start(c, nil) // setup the request
	return ctx
}

func (tpl *Template[C]) context(r *http.Request, funcs ...FlagFunc) (ctx *tContext[C]) {
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
	Runtime        RuntimeFlags // underlying flags
	updateEmbedded bool         // should we automatically update an embedded RuntimeFlags inside the context?

	ctx context.Context // underlying context for render

	// the main template and context
	eMain    chan error // are we done?
	cWaiting bool
	tMain    *template.Template
	cMain    C

	// the footer template and context
	tFooter *template.Template
	cFooter RuntimeFlags
}

func (ctx *tContext[C]) start(c C, err error) {
	ctx.cMain = c
	ctx.eMain = make(chan error, 1)
	ctx.eMain <- err
}

func (ctx *tContext[C]) startLazy(f func() (C, error)) {
	ctx.eMain = make(chan error, 1)
	go func() {
		defer close(ctx.eMain)

		// compute the result, storing the error
		var err error
		ctx.cMain, err = f()
		ctx.eMain <- err
	}()
}

const mainDelay = time.Second

// Main renders the main template.
func (ctx *tContext[C]) Main() (template.HTML, error) {
	timer := timex.NewTimer()
	defer timex.ReleaseTimer(timer)

	timer.Reset(mainDelay)
	select {

	case err := <-ctx.eMain:
		// we received the result within the given time
		// so we can render it immediatly
		ctx.cWaiting = false
		return ctx.doMain(err)

	case <-timer.C:
		// the template is taking longer than expected.
		// we should display a spinner, and do something later
		ctx.cWaiting = true
		return timeWait, nil
	}
}

// Footer renders the footer template
func (ctx *tContext[C]) Footer() (template.HTML, error) {
	return ctx.renderSafe("footer", ctx.tFooter, ctx.cFooter)
}

const (
	timeWait   = "Loading"
	errUnknown = "An unknown error occured, see the server log for details. "
)

func (ctx *tContext[C]) doMain(err error) (template.HTML, error) {
	if err != nil {
		zerolog.Ctx(ctx.ctx).Err(err).Msg("error lazy loading template")
		return errUnknown, nil
	}

	// if the context has a runtime flags embed, then set the field properly
	if ctx.updateEmbedded {
		reflect.ValueOf(&ctx.cMain).Elem().
			FieldByName(runtimeFlagsName).
			Set(reflect.ValueOf(ctx.Runtime))
	}

	return ctx.renderSafe("main", ctx.tMain, ctx.cMain)
}

func (ctx *tContext[C]) AfterBody() (template.HTML, error) {
	// everything was done already
	if !ctx.cWaiting {
		return "", nil
	}

	// wait for the result to appear
	res, err := ctx.doMain(<-ctx.eMain)
	if err != nil {
		return "", err
	}

	str, err := json.Marshal(string(res))
	if err != nil {
		return "", err
	}

	fix := "<script>document.getElementById('main').innerHTML=" + string(str) + "</script>"

	// hook that is called after the body is complete
	return template.HTML(fix), nil
}

const renderSafeError = "Error displaying page. See server log for details. "

func (ctx *tContext[C]) renderSafe(name string, t *template.Template, c any) (template.HTML, error) {

	// already done
	if err := ctx.ctx.Err(); err != nil {
		return "", err
	}

	value, err, panicked := func() (value template.HTML, err error, panicked bool) {
		var builder strings.Builder

		defer func() {
			if panicked {
				r := recover()
				zerolog.Ctx(ctx.ctx).Error().
					Str("uri", ctx.Runtime.RequestURI).
					Str("name", name).
					Str("panic", fmt.Sprint(r)).
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

		return template.HTML(builder.String()), err, false
	}()

	if err != nil || panicked {
		return renderSafeError, httpx.ErrInternalServerError
	}
	return value, nil
}
