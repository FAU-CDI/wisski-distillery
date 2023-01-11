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
	"golang.org/x/exp/slices"
)

// baseContextName is the name of the [BaseContext] type
var baseContextName = reflectx.TypeOf[BaseContext]().Name()

// BaseContext represents a context shared by all the templates.
//
// This context should always be initialized using the [custom.Update], [custom.New] or [custom.NewForm] functions.
// Other invocations might cause an error at runtime.
type BaseContext struct {
	inited        bool // has this context been inited?
	requestWasNil bool // was the passed request nil

	GeneratedAt time.Time // time this page was generated at

	// Menu and breadcrumbs
	Menu   []component.MenuItem
	Crumbs []component.MenuItem

	CSRF template.HTML // CSRF Field
}

// constants that are used in various parts of the template to render stuff
const (
	errorPrefix template.HTML = `<div style="z-index:10000;position:fixed;top:0;left:0;width:100vh;height:100vw;background:red;text-align:center;padding:10vh 10vw;font-size:xx-large;font-weight:bold">`
	errorSuffix template.HTML = "</div>"

	csrfError       template.HTML = errorPrefix + "CSRF used but not provided" + errorSuffix
	initError       template.HTML = errorPrefix + "<code>BaseContext.use()</code> not called" + errorSuffix
	requestNilError template.HTML = errorPrefix + "<code>BaseContext.use()</code> called with nil request" + errorSuffix
)

// Use updates this context to use the values from the given base.
//
// The given request *must not* be nil.
//
// For convenience the passed context is also returned.
func (tc *BaseContext) use(custom *Custom, r *http.Request, crumbs []component.MenuItem) *BaseContext {
	// tc.custom = custom
	tc.inited = true
	tc.requestWasNil = r == nil

	tc.GeneratedAt = time.Now().UTC()

	// setup the CSRF field
	tc.CSRF = csrfError
	if r != nil {
		tc.CSRF = csrf.TemplateField(r)
	}

	// build the menu
	tc.Menu = custom.BuildMenu(r)

	// build the breadcrumbs
	tc.Crumbs = slices.Clone(crumbs)
	last := len(tc.Crumbs) - 1
	for i := range tc.Crumbs {
		tc.Crumbs[i].Active = i == last
	}

	return tc
}

// DoInitCheck is called by the template to check that the BaseContext was initialized properly
func (bc BaseContext) DoInitCheck() template.HTML {
	if !bc.inited {
		return initError
	}
	if bc.requestWasNil {
		return requestNilError
	}
	return ""
}

// NewForm is like New, but returns a new BaseFormContext
func (custom *Custom) NewForm(context httpx.FormContext, r *http.Request, crumbs []component.MenuItem) (ctx BaseFormContext) {
	ctx.FormContext = context
	ctx.use(custom, r, crumbs)
	return
}

// Update updates an embedded BaseContext field in context.
//
// Assumes that context is a pointer to a struct type.
// If this is not the case, might call panic().
func (custom *Custom) Update(context any, r *http.Request, crumbs []component.MenuItem) *BaseContext {
	ctx := reflect.ValueOf(context).
		Elem().FieldByName(baseContextName).Addr().
		Interface().(*BaseContext)
	ctx.use(custom, r, crumbs)
	return ctx
}

// BaseFormContext combines BaseContext and FormContext
type BaseFormContext struct {
	BaseContext
	httpx.FormContext
}
