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

	// suggestions is how many nearby titles a search that found nothing
	// offers, and nearTitle how close one has to be. Titles we don't have
	// score up to about 0.27 against the catalog; real near misses, 0.3 up.
	suggestions = 3
	nearTitle   = 0.3

	// maxQuery caps what's searched, in characters; a title is never longer.
	maxQuery = 200
)

// Search serves /search?q=, fuzzy over every title in the local catalog.
func (p *Pages) Search(w http.ResponseWriter, r *http.Request) {
	q := clampQuery(r.URL.Query().Get("q"))
	v := views.Search{Query: q}

	if q != "" {
		ctx := r.Context()
		hits, err := p.catalog.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: q, MaxResults: searchResults})
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		v.Searched = true
		v.Media, hits = searchFilters(q, r.URL.Query().Get("medium"), hits)
		v.Results = resultsFor(hits)

		if len(hits) == 0 {
			near, err := p.catalog.TitleSuggestions(ctx, sqlcgen.TitleSuggestionsParams{Q: q, AtLeast: nearTitle, MaxResults: suggestions})
			if err != nil {
				p.serverError(w, r, err)
				return
			}
			for _, t := range near {
				v.Suggestions = append(v.Suggestions, views.Suggestion{Title: t, Href: searchPath(t, "")})
			}
		}
	}

	title := "search | listo"
	if q != "" {
		title = fmt.Sprintf("%s | search | listo", q)
	}
	// results pages stay out of search engines
	head := views.Head{Title: title, Description: "Search listo's catalog.", NoIndex: true}
	// the page has its own search box, so the top bar's stays empty
	chrome, err := p.chrome(r, nil, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	p.render(w, r, http.StatusOK, head, chrome, views.SearchPage(v))
}

// searchFilters is a link per medium among the hits, and the hits the chosen
// one keeps. There are no links when filtering wouldn't change anything, and
// a medium that isn't among the hits shows them all.
func searchFilters(q, medium string, hits []sqlcgen.SearchNodesRow) ([]views.MediumLink, []sqlcgen.SearchNodesRow) {
	var order []string
	labels, counts := map[string]string{}, map[string]int{}
	groups := 0
	for _, h := range hits {
		if h.Medium == nil {
			groups++
			continue
		}
		if counts[*h.Medium] == 0 {
			order = append(order, *h.Medium)
			labels[*h.Medium] = deref(h.MediumLabel)
		}
		counts[*h.Medium]++
	}
	if len(order) < 2 && (len(order) == 0 || groups == 0) {
		return nil, hits
	}
	if counts[medium] == 0 {
		medium = ""
	}

	links := []views.MediumLink{{Label: "all", Href: searchPath(q, ""), Count: len(hits), Active: medium == ""}}
	for _, m := range order {
		links = append(links, views.MediumLink{Medium: m, Label: labels[m], Href: searchPath(q, m), Count: counts[m], Active: m == medium})
	}
	if medium == "" {
		return links, hits
	}
	kept := make([]sqlcgen.SearchNodesRow, 0, counts[medium])
	for _, h := range hits {
		if h.Medium != nil && *h.Medium == medium {
			kept = append(kept, h)
		}
	}
	return links, kept
}

func searchPath(q, medium string) string {
	v := url.Values{"q": {q}}
	if medium != "" {
		v.Set("medium", medium)
	}
	return "/search?" + v.Encode()
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
		// nothing to browse yet
		if c.Entries == 0 {
			continue
		}
		links[0].Count += int(c.Entries)
		links = append(links, views.MediumLink{
			Medium: c.Medium,
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
	return p.chromeWith(r, media, query)
}

// chromeWith is chrome for a page that already has its medium links.
func (p *Pages) chromeWith(r *http.Request, media []views.MediumLink, query string) (views.Chrome, error) {
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

// resultsFor is how hits show, on the search page and in the top bar.
func resultsFor(hits []sqlcgen.SearchNodesRow) []views.Result {
	results := make([]views.Result, 0, len(hits))
	for _, h := range hits {
		res := views.Result{
			Title:  h.Title,
			Href:   NodePath(h.Type, h.PublicID, h.Slug),
			Meta:   deref(h.GroupKind),
			Medium: deref(h.GroupKind),
		}
		if h.Type == "entry" {
			res.Meta = withYear(deref(h.MediumLabel), h.StartDate)
			res.Medium = deref(h.FormatLabel)
			res.Detail = hitDetail(h)
		}
		if h.Matched != h.Title {
			res.Matched = h.Matched
			res.MatchedAs = matchedAs(h.MatchedKind)
		}
		results = append(results, res)
	}
	return results
}

// hitDetail is when an entry came out and how long it is: "spring 2011, 24
// episodes".
func hitDetail(h sqlcgen.SearchNodesRow) string {
	var parts []string
	if h.StartDate != nil && h.StartPrecision != nil {
		parts = append(parts, when(deref(h.Medium), deref(h.Format), *h.StartDate, *h.StartPrecision))
	}
	if kind := deref(h.UnitKind); h.UnitTotal != nil && kind != "none" && kind != "" {
		parts = append(parts, plural(int(*h.UnitTotal), kind, kind+"s"))
	}
	return strings.Join(parts, " · ")
}

// matchedAs names the kind of title a search matched.
func matchedAs(kind string) string {
	switch kind {
	case "native", "romanized", "english", "localized":
		return kind + " title"
	}
	return "alias"
}
