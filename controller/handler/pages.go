package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/views"
	"github.com/khinshankhan/listo/web"
)

// Pages renders the server-side pages.
type Pages struct {
	tags    web.Tags
	catalog Catalog
	// store only names images, for their URLs; it never reads them
	store media.Store
}

func NewPages(tags web.Tags, catalog Catalog, store media.Store) *Pages {
	return &Pages{tags: tags, catalog: catalog, store: store}
}

func (p *Pages) NotFound(w http.ResponseWriter, r *http.Request) {
	head := views.Head{Title: "not found | listo", Description: "There is no page here."}
	// top bar only: a 404 shouldn't cost database reads
	p.render(w, r, http.StatusNotFound, head, views.Chrome{}, views.NotFound())
}

// render writes a page. Pages are the same for everyone, so they're
// cacheable, unless the handler already said otherwise.
func (p *Pages) render(w http.ResponseWriter, r *http.Request, status int, head views.Head, chrome views.Chrome, body templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=60")
	}
	w.WriteHeader(status)

	if err := views.Layout(head, p.tags, chrome).Render(templ.WithChildren(r.Context(), body), w); err != nil {
		logx.Default().Named("render").Error("rendering page",
			fields.String("path", r.URL.Path),
			fields.Error(err))
	}
}
