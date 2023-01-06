package custom

import (
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

const footerName = "footer"

// defaultTemplate is the default footer template
var defaultTemplate = template.Must(template.New("footer.html").Parse(`<p>Powered By WissKI Distillery</p>`))

// Template creates a copy of template with shared template parts updated accordingly.
// Any template using this should use one of the template contexts in this package.
func (custom *Custom) Template(tpl *template.Template) *template.Template {
	tree := defaultTemplate.Tree.Copy()

	clone := template.Must(tpl.Clone())                 // create a clone of the template
	template.Must(clone.AddParseTree(footerName, tree)) // add the parse tree to it
	return clone                                        // and return the tree
}

// NewContext returns a new BaseContext
func (custom *Custom) New() (ctx BaseContext) {
	ctx.Use(custom.Base)
	return
}

// NewForm is like New, but returns a new BaseFormContext
func (custom *Custom) NewForm(context httpx.FormContext) (ctx BaseFormContext) {
	ctx.FormContext = context
	ctx.Use(custom.Base)
	return
}

// RenderContext can be used as httpx.Form.RenderTemplateContext.
// It returns a new [BaseFormContext].
func (custom *Custom) RenderContext(ctx httpx.FormContext, r *http.Request) any {
	return BaseFormContext{
		FormContext: ctx,
		BaseContext: custom.New(),
	}
}

// Update updates an embedded BaseContext field in context.
//
// Assumes that context is a pointer to a struct type.
// If this is not the case, might call panic().
func (custom *Custom) Update(context any) *BaseContext {
	ctx := reflect.ValueOf(context).
		Elem().FieldByName(contextName).Addr().
		Interface().(*BaseContext)
	ctx.Use(custom.Base)
	return ctx
}

// contextName is the name of the [BaseContext] field.
var contextName = reflectx.TypeOf[BaseContext]().Name()

// BaseContext is a context struct shared by all contexts
type BaseContext struct {
	Time time.Time // time this page was generated at
}

// Use updates this context to use the values from the given base.
// For convenience the passed context is also returned.
func (tc *BaseContext) Use(base component.Base) *BaseContext {
	tc.Time = time.Now().UTC()
	return tc
}

// BaseFormContext combines BaseContext and FormContext
type BaseFormContext struct {
	BaseContext
	httpx.FormContext
}
