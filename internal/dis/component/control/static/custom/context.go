package custom

import (
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/gorilla/csrf"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

// baseContextName is the name of the [BaseContext] type
var baseContextName = reflectx.TypeOf[BaseContext]().Name()

// BaseContext represents a context shared by all the templates.
//
// This context should always be initialized using the [custom.Update], [custom.New] or [custom.NewForm] functions.
// Other invocations might cause an error at runtime.
type BaseContext struct {
	inited bool // has this context been inited

	GeneratedAt time.Time // time this page was generated at

	CSRF template.HTML // CSRF Field
}

// constants that are used in various parts of the template to render stuff
const (
	errorPrefix template.HTML = `<div style="z-index:10000;position:fixed;top:0;left:0;width:100vh;height:100vw;background:red;text-align:center;padding:10vh 10vw;font-size:xx-large;font-weight:bold">`
	errorSuffix template.HTML = "</div>"

	csrfError template.HTML = errorPrefix + "CSRF used but not provided" + errorSuffix
	initError template.HTML = errorPrefix + "BaseContext not initialized" + errorSuffix
)

// Use updates this context to use the values from the given base.
// For convenience the passed context is also returned.
func (tc *BaseContext) use(base component.Base, r *http.Request) *BaseContext {
	tc.inited = true

	tc.GeneratedAt = time.Now().UTC()

	// setup the CSRF field
	tc.CSRF = csrfError
	if r != nil {
		tc.CSRF = csrf.TemplateField(r)
	}

	return tc
}

func (bc BaseContext) DoInitCheck() template.HTML {
	if !bc.inited {
		return initError
	}
	return ""
}

// NewForm is like New, but returns a new BaseFormContext
func (custom *Custom) NewForm(context httpx.FormContext, r *http.Request) (ctx BaseFormContext) {
	ctx.FormContext = context
	ctx.use(custom.Base, r)
	return
}

// RenderContext is exactly like NewForm, but returns any to be used as httpx.Form.RenderTemplateContext
func (custom *Custom) RenderContext(ctx httpx.FormContext, r *http.Request) any {
	return custom.NewForm(ctx, r)
}

// Update updates an embedded BaseContext field in context.
//
// Assumes that context is a pointer to a struct type.
// If this is not the case, might call panic().
func (custom *Custom) Update(context any, r *http.Request) *BaseContext {
	ctx := reflect.ValueOf(context).
		Elem().FieldByName(baseContextName).Addr().
		Interface().(*BaseContext)
	ctx.use(custom.Base, r)
	return ctx
}

// BaseFormContext combines BaseContext and FormContext
type BaseFormContext struct {
	BaseContext
	httpx.FormContext
}
