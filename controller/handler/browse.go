package handler

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/views"
)

// perShelf is how many entries each medium's shelf holds: more than the
// widest row of covers, so every shelf has more to scroll to.
const perShelf = 24

// browsePage is how many entries a page of results holds. Every column count
// the covers take divides it, so the last row is always full.
const browsePage = 48

// viewCookie remembers how someone likes to browse.
const viewCookie = "listo_view"

// the filters, in the order the form and its URLs use them
var browseFilters = []string{"medium", "format", "genre", "year", "season"}

// seasons are a season filter's values, winter first: the query takes the
// index
var seasons = []string{"winter", "spring", "summer", "fall"}

// browseViews are the ways to show entries, the default first.
var browseViews = []struct{ key, label string }{{"cards", "cards"}, {"list", "list"}, {"text", "text only"}}

// Entries is the catalog: a shelf per medium, or with ?medium=, ?format=,
// ?genre=, ?year= or ?season=, a page of everything that matches. ?view=
// picks cards, a list or plain text, and is remembered. A value the filters
// don't offer, or a page past the end, is a 404, so each view has one
// address; a format goes with its medium, so one without it is dropped.
func (p *Pages) Entries(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	clean := url.Values{}
	dirty := false
	for _, name := range append(slices.Clone(browseFilters), "view", "page") {
		switch v := query.Get(name); {
		case !query.Has(name):
		// "any" from the form, or the first page: the plain address has it
		case v == "" || (name == "page" && v == "1"):
			dirty = true
		default:
			clean.Set(name, v)
		}
	}
	if dirty {
		http.Redirect(w, r, browsePath(clean), http.StatusFound)
		return
	}

	view, ok := p.browseView(w, r, clean.Get("view"))
	if !ok {
		p.NotFound(w, r)
		return
	}
	page := 1
	if s := clean.Get("page"); s != "" {
		n, err := strconv.Atoi(s)
		// past this, the page's offset no longer fits the query's int32
		if err != nil || n < 2 || n > math.MaxInt32/browsePage {
			p.NotFound(w, r)
			return
		}
		page = n
	}

	ctx := r.Context()
	medium := clean.Get("medium")
	all, err := p.catalog.BrowseFormats(ctx)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	var formats []sqlcgen.CatalogMediumFormat
	for _, f := range all {
		if f.Medium == medium {
			formats = append(formats, f)
		}
	}
	// one left over from another type in the form, or with no type at all
	if f := clean.Get("format"); f != "" && !slices.ContainsFunc(formats, func(row sqlcgen.CatalogMediumFormat) bool { return row.Format == f }) {
		clean.Del("format")
		http.Redirect(w, r, browsePath(clean), http.StatusFound)
		return
	}

	media, err := p.mediumLinks(r, &medium)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	filters, args, ok, err := p.browseFilters(ctx, media, formats, clean)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	filtered := slices.ContainsFunc(browseFilters, clean.Has)
	// the shelves are one page
	if !ok || (!filtered && page > 1) {
		p.NotFound(w, r)
		return
	}

	v := views.Browse{Filters: filters, Filtered: filtered, View: view, KeepView: clean.Get("view"), Formats: formatOptions(all)}
	for _, bv := range browseViews {
		v.Views = append(v.Views, views.SortLink{Label: bv.label, Href: browsePath(with(clean, "view", bv.key)), Active: bv.key == view})
	}
	if filtered {
		args.PageSize, args.Skip = browsePage, int32((page-1)*browsePage)
		rows, err := p.catalog.BrowseEntries(ctx, args)
		if err != nil {
			p.serverError(w, r, err)
			return
		}
		if len(rows) == 0 && page > 1 {
			p.NotFound(w, r)
			return
		}
		v.Cards = []views.Card{}
		for _, e := range rows {
			v.Cards = append(v.Cards, card(p.store, pagedCard(e)))
			v.Total = int(e.Total)
		}
		v.Page, v.Pages = page, max(1, (v.Total+browsePage-1)/browsePage)
		if page > 1 {
			v.Prev = browsePath(with(clean, "page", strconv.Itoa(page-1)))
		}
		if page < v.Pages {
			v.Next = browsePath(with(clean, "page", strconv.Itoa(page+1)))
		}
	} else {
		if v.Shelves, err = p.shelves(ctx, media, clean.Get("view")); err != nil {
			p.serverError(w, r, err)
			return
		}
	}

	head := views.Head{Title: "browse | listo", Description: "Every entry in listo's catalog.", Canonical: "/entries"}
	switch {
	case medium != "" && !clean.Has("format") && !clean.Has("genre") && !clean.Has("year") && !clean.Has("season"):
		head.Canonical = browsePath(url.Values{"medium": {medium}, "page": clean["page"]})
	case filtered:
		// every mix of filters is a page; only the plain ones are worth indexing
		head.NoIndex = true
	}
	if page > 1 {
		head.Title = fmt.Sprintf("browse, page %d | listo", page)
	}
	// reuse the links rather than counting again
	chrome, err := p.chromeWith(r, media, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	chrome.Wide = true
	// the view can come from a cookie, so no shared cache may keep this
	w.Header().Set("Cache-Control", "private, max-age=60")
	w.Header().Add("Vary", "Cookie")
	p.render(w, r, http.StatusOK, head, chrome, views.BrowsePage(v))
}

// browseView is the view asked for, remembered for next time, or else the one
// remembered, or else cards. ok is false for a view that doesn't exist.
func (p *Pages) browseView(w http.ResponseWriter, r *http.Request, asked string) (string, bool) {
	known := func(v string) bool {
		return slices.ContainsFunc(browseViews, func(bv struct{ key, label string }) bool { return bv.key == v })
	}
	if asked != "" {
		if !known(asked) {
			return "", false
		}
		http.SetCookie(w, &http.Cookie{Name: viewCookie, Value: asked, Path: "/entries", MaxAge: 365 * 24 * 60 * 60,
			HttpOnly: true, SameSite: http.SameSiteLaxMode})
		return asked, true
	}
	if c, err := r.Cookie(viewCookie); err == nil && known(c.Value) {
		return c.Value, true
	}
	return browseViews[0].key, true
}

// browsePath is /entries with these values, in a stable order.
func browsePath(values url.Values) string {
	if len(values) == 0 {
		return "/entries"
	}
	return "/entries?" + values.Encode()
}

// formatPath browses one format, which only means something with its medium.
func formatPath(medium, format string) string {
	return browsePath(url.Values{"medium": {medium}, "format": {format}})
}

// with is values with one more set. Page one is the plain address, so it
// drops the page instead.
func with(values url.Values, key, value string) url.Values {
	out := url.Values{}
	for k, v := range values {
		out[k] = slices.Clone(v)
	}
	if key == "page" && value == "1" {
		out.Del(key)
		return out
	}
	out.Set(key, value)
	return out
}

// browseFilters is the form's dropdowns with the chosen values selected, and
// the query those values make. ok is false when a value isn't on offer.
// formats are the chosen medium's; with none chosen, there are none to pick.
func (p *Pages) browseFilters(ctx context.Context, media []views.MediumLink, formats []sqlcgen.CatalogMediumFormat, chosen url.Values) ([]views.Filter, sqlcgen.BrowseEntriesParams, bool, error) {
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
	var formatValues, formatLabels []string
	for _, f := range formats {
		formatValues = append(formatValues, f.Format)
		formatLabels = append(formatLabels, f.Label)
	}
	formatFilter, _ := pick("format", "format", formatValues, formatLabels)
	if len(formats) == 0 {
		formatFilter.Options[0].Label, formatFilter.Disabled = noType, true
	}
	genreFilter, okGenre := pick("genre", "genre", genres, genres)
	yearFilter, okYear := pick("year", "year", yearValues, yearValues)
	seasonFilter, okSeason := pick("season", "season", seasons, seasons)
	if !okType || !okGenre || !okYear || !okSeason {
		return nil, args, false, nil
	}

	if m := chosen.Get("medium"); m != "" {
		args.Medium = &m
	}
	if f := chosen.Get("format"); f != "" {
		args.Format = &f
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
	return []views.Filter{typeFilter, formatFilter, genreFilter, yearFilter, seasonFilter}, args, true, nil
}

// noType is the format filter with no type to have formats.
const noType = "pick a type first"

// formatOptions is the format filter's options for each type, and for no
// type at all. formats come grouped by medium.
func formatOptions(formats []sqlcgen.CatalogMediumFormat) []views.FormatOptions {
	out := []views.FormatOptions{{Options: []views.Option{{Label: noType}}}}
	for i, f := range formats {
		if i == 0 || f.Medium != formats[i-1].Medium {
			out = append(out, views.FormatOptions{Medium: f.Medium, Options: []views.Option{{Label: "any"}}})
		}
		last := &out[len(out)-1]
		last.Options = append(last.Options, views.Option{Value: f.Format, Label: f.Label})
	}
	return out
}

// shelves is each medium's newest entries, the media in the sidebar's order.
func (p *Pages) shelves(ctx context.Context, media []views.MediumLink, view string) ([]views.Shelf, error) {
	rows, err := p.catalog.BrowseShelves(ctx, perShelf)
	if err != nil {
		return nil, fmt.Errorf("shelves: %w", err)
	}
	byMedium := map[string][]views.Card{}
	for _, e := range rows {
		byMedium[e.Medium] = append(byMedium[e.Medium], card(p.store, e))
	}
	var out []views.Shelf
	for _, m := range media[1:] {
		if cards := byMedium[m.Medium]; len(cards) > 0 {
			all := url.Values{"medium": {m.Medium}}
			if view != "" {
				all.Set("view", view)
			}
			out = append(out, views.Shelf{ID: m.Medium, Title: m.Label, Href: browsePath(all), Cards: cards})
		}
	}
	return out, nil
}

// pagedCard is a page row as a shelf row: the same columns, less the count.
func pagedCard(e sqlcgen.BrowseEntriesRow) sqlcgen.BrowseShelvesRow {
	return sqlcgen.BrowseShelvesRow{PublicID: e.PublicID, Slug: e.Slug, Title: e.Title, Medium: e.Medium, MediumLabel: e.MediumLabel,
		Format: e.Format, FormatLabel: e.FormatLabel, StartDate: e.StartDate, StartPrecision: e.StartPrecision,
		UnitKind: e.UnitKind, UnitTotal: e.UnitTotal,
		CoverKey: e.CoverKey, CoverWidth: e.CoverWidth, CoverHeight: e.CoverHeight, CoverFit: e.CoverFit}
}

// card is an entry for browsing: its format, and the rest of its meta eg
// "spring 2011 · 24 episodes".
func card(store media.Store, e sqlcgen.BrowseShelvesRow) views.Card {
	var meta []string
	if e.StartDate != nil && e.StartPrecision != nil {
		meta = append(meta, when(e.Medium, e.Format, *e.StartDate, *e.StartPrecision))
	}
	if e.UnitTotal != nil && e.UnitKind != "none" {
		meta = append(meta, plural(int(*e.UnitTotal), e.UnitKind, e.UnitKind+"s"))
	}
	return views.Card{Title: e.Title, Href: EntryPath(e.PublicID, e.Slug), Format: e.FormatLabel,
		FormatHref: formatPath(e.Medium, e.Format), Meta: strings.Join(meta, " · "),
		Cover: coverOf(store, e.MediumLabel, e.Title, e.StartDate, e.CoverKey, e.CoverWidth, e.CoverHeight, e.CoverFit)}
}
