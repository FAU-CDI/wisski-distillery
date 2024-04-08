package handling

import (
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/rs/zerolog"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/content"
	"github.com/tkw1536/pkglib/lazy"
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

		zerolog.Ctx(r.Context()).
			Err(err).
			Str("path", r.URL.Path).
			Msg("unknown error")
	}
	return parent
}

func (h *Handling) Redirect(Handler content.RedirectFunc) http.Handler {
	r := content.Redirect(Handler)
	r.Interceptor = h.TextInterceptor()
	return r
}

func (h *Handling) WriteHTML(context any, err error, template *template.Template, w http.ResponseWriter, r *http.Request) error {
	return content.WriteHTMLI(context, err, template, h.HTMLInterceptor(), w, r)
}
