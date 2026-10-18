package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/views"
	"github.com/khinshankhan/listo/web"
)

// Pages renders the server-side pages.
type Pages struct {
	tags web.Tags
}

func NewPages(tags web.Tags) *Pages {
	return &Pages{tags: tags}
}

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	head := views.Head{Title: "listo", Description: "The 1 list for all your content."}
	p.render(w, r, http.StatusOK, head, views.Hello("world"))
}

func (p *Pages) There(w http.ResponseWriter, r *http.Request) {
	head := views.Head{Title: "there | listo", Description: "Hello there."}
	p.render(w, r, http.StatusOK, head, views.Hello("there"))
}

func (p *Pages) NotFound(w http.ResponseWriter, r *http.Request) {
	head := views.Head{Title: "not found | listo", Description: "There is no page here."}
	p.render(w, r, http.StatusNotFound, head, views.NotFound())
}

// render writes a page. Pages are the same for everyone, so they're cacheable.
func (p *Pages) render(w http.ResponseWriter, r *http.Request, status int, head views.Head, body templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=60")
	w.WriteHeader(status)

	if err := views.Layout(head, p.tags).Render(templ.WithChildren(r.Context(), body), w); err != nil {
		logx.Default().Named("render").Error("rendering page",
			fields.String("path", r.URL.Path),
			fields.Error(err))
	}
}
