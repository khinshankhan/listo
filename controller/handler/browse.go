package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// perShelf is how many entries each medium's shelf shows.
const perShelf = 6

// the filters, in the order the form and its URLs use them
var browseFilters = []string{"medium", "genre", "year", "season"}

// seasons are a season filter's values, winter first: the query takes the
// index
var seasons = []string{"winter", "spring", "summer", "fall"}

// Entries is the catalog: a shelf per medium, or with ?medium=, ?genre=,
// ?year= or ?season=, everything that matches. A value the filters don't
// offer is a 404, so each view has one address.
func (p *Pages) Entries(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// the form sends "any" as empty; drop those for the clean address
	clean := url.Values{}
	for _, name := range browseFilters {
		if v := query.Get(name); v != "" {
			clean.Set(name, v)
		}
	}
	if dirty := slices.ContainsFunc(browseFilters, func(name string) bool { return query.Has(name) && query.Get(name) == "" }); dirty {
		to := "/entries"
		if len(clean) > 0 {
			to += "?" + clean.Encode()
		}
		http.Redirect(w, r, to, http.StatusFound)
		return
	}

	ctx := r.Context()
	medium := clean.Get("medium")
	media, err := p.mediumLinks(r, &medium)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	filters, args, ok, err := p.browseFilters(ctx, media, clean)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	if !ok {
		p.NotFound(w, r)
		return
	}

	v := views.Browse{Filters: filters, Filtered: len(clean) > 0}
	if v.Filtered {
		rows, err := p.catalog.BrowseEntries(ctx, args)
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		v.Cards = []views.Card{}
		for _, e := range rows {
			v.Cards = append(v.Cards, card(e))
		}
	} else {
		if v.Shelves, err = p.shelves(ctx, media); err != nil {
			p.serverError(w, r, err)
			return
		}
	}

	head := views.Head{Title: "browse | listo", Description: "Every entry in listo's catalog.", Canonical: "/entries"}
	switch {
	case medium != "" && len(clean) == 1:
		head.Canonical = "/entries?medium=" + url.QueryEscape(medium)
	case len(clean) > 0:
		// every mix of filters is a page; only the plain ones are worth indexing
		head.NoIndex = true
	}
	// reuse the links rather than counting again
	chrome, err := p.chromeWith(r, media, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	chrome.HideSearch = true
	p.render(w, r, http.StatusOK, head, chrome, views.BrowsePage(v))
}

// browseFilters is the form's dropdowns with the chosen values selected, and
// the query those values make. ok is false when a value isn't on offer.
func (p *Pages) browseFilters(ctx context.Context, media []views.MediumLink, chosen url.Values) ([]views.Filter, sqlcgen.BrowseEntriesParams, bool, error) {
	var args sqlcgen.BrowseEntriesParams
	pick := func(name, label string, values, labels []string) (views.Filter, bool) {
		f := views.Filter{Name: name, Label: label, Options: []views.Option{{Label: "any", Selected: !chosen.Has(name)}}}
		found := !chosen.Has(name)
		for i, val := range values {
			selected := chosen.Get(name) == val
			found = found || selected
			f.Options = append(f.Options, views.Option{Value: val, Label: labels[i], Selected: selected})
		}
		return f, found
	}

	var mediaValues, mediaLabels []string
	for _, m := range media[1:] {
		mediaValues = append(mediaValues, m.Medium)
		mediaLabels = append(mediaLabels, m.Label)
	}
	genres, err := p.catalog.BrowseGenres(ctx)
	if err != nil {
		return nil, args, false, fmt.Errorf("genres: %w", err)
	}
	years, err := p.catalog.BrowseYears(ctx)
	if err != nil {
		return nil, args, false, fmt.Errorf("years: %w", err)
	}
	var yearValues []string
	for _, y := range years {
		yearValues = append(yearValues, strconv.Itoa(int(y)))
	}

	typeFilter, okType := pick("medium", "type", mediaValues, mediaLabels)
	genreFilter, okGenre := pick("genre", "genre", genres, genres)
	yearFilter, okYear := pick("year", "year", yearValues, yearValues)
	seasonFilter, okSeason := pick("season", "season", seasons, seasons)
	if !okType || !okGenre || !okYear || !okSeason {
		return nil, args, false, nil
	}

	if m := chosen.Get("medium"); m != "" {
		args.Medium = &m
	}
	if g := chosen.Get("genre"); g != "" {
		args.Genre = &g
	}
	if y := chosen.Get("year"); y != "" {
		n, _ := strconv.Atoi(y) // offered, so a number
		year := int32(n)
		args.Year = &year
	}
	if s := chosen.Get("season"); s != "" {
		i := int32(slices.Index(seasons, s))
		args.Season = &i
	}
	return []views.Filter{typeFilter, genreFilter, yearFilter, seasonFilter}, args, true, nil
}

// shelves is each medium's newest entries, the media in the sidebar's order.
func (p *Pages) shelves(ctx context.Context, media []views.MediumLink) ([]views.Shelf, error) {
	rows, err := p.catalog.BrowseShelves(ctx, perShelf)
	if err != nil {
		return nil, fmt.Errorf("shelves: %w", err)
	}
	byMedium := map[string][]views.Card{}
	for _, e := range rows {
		// the same columns, so the same card
		byMedium[e.Medium] = append(byMedium[e.Medium], card(sqlcgen.BrowseEntriesRow(e)))
	}
	var out []views.Shelf
	for _, m := range media[1:] {
		if cards := byMedium[m.Medium]; len(cards) > 0 {
			out = append(out, views.Shelf{ID: m.Medium, Title: m.Label, Href: m.Href, Cards: cards})
		}
	}
	return out, nil
}

// card is an entry as a cover, its meta eg "anime, spring 2011, 24 episodes".
func card(e sqlcgen.BrowseEntriesRow) views.Card {
	meta := []string{e.FormatLabel}
	if e.StartDate != nil && e.StartPrecision != nil {
		meta = append(meta, when(e.Medium, e.Format, *e.StartDate, *e.StartPrecision))
	}
	if e.UnitTotal != nil && e.UnitKind != "none" {
		meta = append(meta, plural(int(*e.UnitTotal), e.UnitKind, e.UnitKind+"s"))
	}
	return views.Card{Title: e.Title, Href: EntryPath(e.PublicID, e.Slug), Meta: strings.Join(meta, " · ")}
}
