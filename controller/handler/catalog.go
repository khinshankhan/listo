package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
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

	// sidebarFranchises is how many franchises the sidebar lists, so it's
	// the same height on every page however big the catalog gets; the rest
	// are on /groups.
	sidebarFranchises = 6
)

// Search serves /search?q=, fuzzy over every title in the local catalog.
func (p *Pages) Search(w http.ResponseWriter, r *http.Request) {
	q := clampQuery(r.URL.Query().Get("q"))
	v := views.Search{Query: q}

	if q != "" {
		ctx := r.Context()
		counts, err := p.catalog.SearchMedia(ctx, q)
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		var medium string
		v.Media, medium, v.Total = searchFilters(q, r.URL.Query().Get("medium"), counts)
		args := sqlcgen.SearchNodesParams{Q: q, MaxResults: searchResults}
		if medium != "" {
			args.Medium = &medium
		}
		hits, err := p.catalog.SearchNodes(ctx, args)
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		v.Searched = true
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

// searchFilters is a link per medium the search finds, from counts over every
// hit, not just the page shown; the medium it narrows to; and how many hits
// that leaves. There are no links when filtering wouldn't change anything,
// and a medium with no hits shows them all.
func searchFilters(q, chosen string, counts []sqlcgen.SearchMediaRow) ([]views.MediumLink, string, int) {
	all, groups := 0, 0
	byMedium := map[string]int{}
	var media []sqlcgen.SearchMediaRow
	for _, c := range counts {
		all += int(c.Hits)
		if c.Medium == nil {
			groups += int(c.Hits)
			continue
		}
		byMedium[*c.Medium] = int(c.Hits)
		media = append(media, c)
	}
	medium, total := "", all
	if n := byMedium[chosen]; n > 0 {
		medium, total = chosen, n
	}
	if len(media) < 2 && (len(media) == 0 || groups == 0) {
		return nil, medium, total
	}

	links := []views.MediumLink{{Label: "all", Href: searchPath(q, ""), Count: all, Active: medium == ""}}
	for _, m := range media {
		links = append(links, views.MediumLink{Medium: *m.Medium, Label: deref(m.MediumLabel), Href: searchPath(q, *m.Medium),
			Count: int(m.Hits), Active: *m.Medium == medium})
	}
	return links, medium, total
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

	side := &views.Sidebar{Media: media, Franchises: biggest(groupCards(groups), sidebarFranchises)}
	return views.Chrome{Query: query, Sidebar: side}, nil
}

// groupCards is the top-level groups as cards, in the order they came.
func groupCards(groups []sqlcgen.TopGroupsRow) []views.GroupCard {
	cards := make([]views.GroupCard, 0, len(groups))
	for _, g := range groups {
		cards = append(cards, views.GroupCard{
			Name:    g.Name,
			Kind:    g.Kind,
			Href:    GroupPath(g.PublicID, g.Slug),
			Entries: int(g.Entries),
		})
	}
	return cards
}

// biggest is the n cards with the most entries, most first; a tie keeps the
// order they came in.
func biggest(cards []views.GroupCard, n int) []views.GroupCard {
	out := slices.Clone(cards)
	slices.SortStableFunc(out, func(a, b views.GroupCard) int { return b.Entries - a.Entries })
	return out[:min(n, len(out))]
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
			res.FormatHref = formatPath(deref(h.Medium), deref(h.Format))
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
