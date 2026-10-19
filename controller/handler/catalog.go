package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

const (
	// searchResults is how many hits a search page shows.
	searchResults = 30

	// maxQuery caps what's searched, in characters; a title is never longer.
	maxQuery = 200
)

// Home is search, the top-level groups, and the media.
func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	groups, err := p.catalog.TopGroups(ctx)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	media, err := p.mediumLinks(r, nil)
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	v := views.Home{Media: media}
	for _, g := range groups {
		v.Groups = append(v.Groups, views.GroupCard{
			Name:    g.Name,
			Kind:    g.Kind,
			Href:    GroupPath(g.PublicID, g.Slug),
			Entries: int(g.Entries),
		})
	}

	head := views.Head{Title: "listo", Description: "The 1 list for all your content: anime, manga, games, books and more."}
	// the home page is top bar only: its body is the catalog
	p.render(w, r, http.StatusOK, head, views.Chrome{}, views.HomePage(v))
}

// Search serves /search?q=, fuzzy over every title in the local catalog.
func (p *Pages) Search(w http.ResponseWriter, r *http.Request) {
	q := clampQuery(r.URL.Query().Get("q"))
	v := views.Search{Query: q}

	if q != "" {
		hits, err := p.catalog.SearchNodes(r.Context(), sqlcgen.SearchNodesParams{Q: q, MaxResults: searchResults})
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		v.Searched = true
		for _, h := range hits {
			res := views.Result{
				Title: h.Title,
				Href:  NodePath(h.Type, h.PublicID, h.Slug),
				Meta:  deref(h.GroupKind),
			}
			if h.Type == "entry" {
				res.Meta = withYear(deref(h.MediumLabel), h.StartDate)
			}
			if h.Matched != h.Title {
				res.Matched = h.Matched
			}
			v.Results = append(v.Results, res)
		}
	}

	title := "search | listo"
	if q != "" {
		title = fmt.Sprintf("%s | search | listo", q)
	}
	// results pages stay out of search engines
	head := views.Head{Title: title, Description: "Search listo's catalog.", NoIndex: true}
	chrome, err := p.chrome(r, nil, q)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	p.render(w, r, http.StatusOK, head, chrome, views.SearchPage(v))
}

// Entries lists every entry, or one medium's with ?medium=. An unknown or
// empty medium is a 404.
func (p *Pages) Entries(w http.ResponseWriter, r *http.Request) {
	medium := r.URL.Query().Get("medium")

	media, err := p.mediumLinks(r, &medium)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	var filter *string
	if medium != "" {
		known := false
		for _, m := range media {
			known = known || m.Active
		}
		if !known {
			p.NotFound(w, r)
			return
		}
		filter = &medium
	}

	entries, err := p.catalog.BrowseEntries(r.Context(), filter)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	v := views.Browse{Media: media}
	for _, e := range entries {
		v.Entries = append(v.Entries, views.Row{
			Title: e.Title,
			Href:  EntryPath(e.PublicID, e.Slug),
			Meta:  withYear(e.MediumLabel, e.StartDate),
		})
	}

	head := views.Head{Title: "browse | listo", Description: "Every entry in listo's catalog.", Canonical: "/entries"}
	if filter != nil {
		head.Canonical = "/entries?medium=" + url.QueryEscape(medium)
	}
	// with no filter, medium is "" and "everything" is the active link
	chrome, err := p.chrome(r, &medium, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	p.render(w, r, http.StatusOK, head, chrome, views.BrowsePage(v))
}

// mediumLinks is "all" plus a link per medium with entries. active marks one:
// nil none, "" all.
func (p *Pages) mediumLinks(r *http.Request, active *string) ([]views.MediumLink, error) {
	counts, err := p.catalog.MediumCounts(r.Context())
	if err != nil {
		return nil, fmt.Errorf("medium counts: %w", err)
	}

	all := views.MediumLink{Href: "/entries", Active: active != nil && *active == ""}
	links := []views.MediumLink{all}
	for _, c := range counts {
		links[0].Count += int(c.Entries)
		links = append(links, views.MediumLink{
			Label:  c.Label,
			Href:   "/entries?medium=" + url.QueryEscape(c.Medium),
			Count:  int(c.Entries),
			Active: active != nil && c.Medium == *active,
		})
	}
	return links, nil
}

// chrome is the sidebar and top bar for a catalog page.
func (p *Pages) chrome(r *http.Request, medium *string, query string) (views.Chrome, error) {
	media, err := p.mediumLinks(r, medium)
	if err != nil {
		return views.Chrome{}, err
	}
	groups, err := p.catalog.TopGroups(r.Context())
	if err != nil {
		return views.Chrome{}, fmt.Errorf("top groups: %w", err)
	}

	side := &views.Sidebar{Media: media}
	for _, g := range groups {
		side.Franchises = append(side.Franchises, views.GroupCard{
			Name:    g.Name,
			Kind:    g.Kind,
			Href:    GroupPath(g.PublicID, g.Slug),
			Entries: int(g.Entries),
		})
	}
	return views.Chrome{Query: query, Sidebar: side}, nil
}

// clampQuery trims q and caps it at maxQuery runes, never splitting one.
func clampQuery(q string) string {
	q = strings.TrimSpace(q)
	if utf8.RuneCountInString(q) <= maxQuery {
		return q
	}
	return string([]rune(q)[:maxQuery])
}
