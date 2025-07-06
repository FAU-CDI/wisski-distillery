//spellchecker:words handling
package handling

//spellchecker:words html template http github wisski distillery internal component wdlog pkglib httpx content lazy
import (
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/httpx"
	"go.tkw01536.de/pkglib/httpx/content"
	"go.tkw01536.de/pkglib/lazy"
)

type Handling struct {
	component.Base

	text lazy.Lazy[httpx.ErrInterceptor]
	html lazy.Lazy[httpx.ErrInterceptor]
}

func (h *Handling) TextInterceptor() httpx.ErrInterceptor {
	return h.text.Get(func() httpx.ErrInterceptor {
		return h.interceptor(httpx.TextInterceptor)
	})
}

func (h *Handling) HTMLInterceptor() httpx.ErrInterceptor {
	return h.html.Get(func() httpx.ErrInterceptor {
		return h.interceptor(httpx.TextInterceptor)
	})
}

// Interceptor returns a copy of the parent interceptor with global distillery interceptor options enabled.
func (h *Handling) interceptor(parent httpx.ErrInterceptor) httpx.ErrInterceptor {
	pf := parent.OnFallback
	if pf == nil {
		pf = func(r *http.Request, err error) {}
	}

	config := component.GetStill(h).Config
	parent.RenderError = config.HTTP.Debug.Set && config.HTTP.Debug.Value
	parent.OnFallback = func(r *http.Request, err error) {
		pf(r, err)

		wdlog.Of(r.Context()).Error(
			"unknown error",
			"error", err,

			"path", r.URL.Path,
		)
	}
	return parent
}

func (h *Handling) Redirect(handler content.RedirectFunc) http.Handler {
	r := content.Redirect(handler)
	r.Interceptor = h.TextInterceptor()
	return r
}

func (h *Handling) WriteHTML(context any, err error, template *template.Template, w http.ResponseWriter, r *http.Request) error {
	return LogTemplateError(r, content.WriteHTMLI(context, err, template, h.HTMLInterceptor(), w, r))
}

func LogTemplateError(r *http.Request, err error) error {
	if err != nil {
		wdlog.Of(r.Context()).Error(
			"error rendering template",
			"error", err,
			"path", r.URL.String(),
		)
	}
	return err
}
