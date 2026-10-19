package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// Entry serves /entries/{id}[/{slug}]. The id is the key; a wrong or missing
// slug redirects to the canonical URL.
func (p *Pages) Entry(w http.ResponseWriter, r *http.Request) {
	id, ok := publicID(r)
	if !ok {
		p.NotFound(w, r)
		return
	}

	ctx := r.Context()
	e, err := p.catalog.EntryByPublicID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		p.elsewhere(w, r, id)
		return
	}
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	canonical := EntryPath(e.PublicID, e.Slug)
	if redirectCanonical(w, r, canonical, e.Slug) {
		return
	}

	v, err := p.entryView(ctx, e)
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	head := views.Head{
		Title:       fmt.Sprintf("%s | listo", e.Title),
		Description: describe(v),
		Canonical:   canonical,
	}
	chrome, err := p.chrome(r, &e.Medium, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	p.render(w, r, http.StatusOK, head, chrome, views.EntryPage(v))
}

// entryView gathers everything the page shows into display form.
func (p *Pages) entryView(ctx context.Context, e sqlcgen.EntryByPublicIDRow) (views.Entry, error) {
	v := views.Entry{
		Title:      e.Title,
		Medium:     e.MediumLabel,
		MediumHref: "/entries?medium=" + url.QueryEscape(e.Medium),
		Eyebrow:    []string{e.MediumLabel},
	}
	if e.FormatLabel != e.MediumLabel {
		v.Format = e.FormatLabel
		v.Eyebrow = append(v.Eyebrow, e.FormatLabel)
	}

	v.Facts = append(v.Facts, views.Fact{Label: "format", Value: e.FormatLabel})
	if e.UnitTotal != nil && e.UnitKind != "none" {
		v.Facts = append(v.Facts, views.Fact{Label: e.UnitKind + "s", Value: fmt.Sprint(*e.UnitTotal)})
	}
	if e.StartDate != nil && e.StartPrecision != nil {
		v.Year = e.StartDate.Format("2006")
		v.Eyebrow = append(v.Eyebrow, when(e.Medium, e.Format, *e.StartDate, *e.StartPrecision))
		v.Facts = append(v.Facts, views.Fact{Label: startLabel(e.Medium, e.Format), Value: released(*e.StartDate, *e.StartPrecision)})
	}
	country := deref(e.Country)
	if country != "" {
		v.Facts = append(v.Facts, views.Fact{Label: "country", Value: country})
	}

	parts, err := p.nodeParts(ctx, e.ID, country)
	if err != nil {
		return v, err
	}
	// the heading already shows the display title
	for _, t := range parts.titles {
		if t.Title != e.Title {
			v.Titles = append(v.Titles, t)
		}
	}
	v.Relations, v.Groups, v.Links = parts.relations, parts.groups, parts.links

	tags, err := p.catalog.NodeTags(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("tags: %w", err)
	}
	for _, t := range tags {
		v.Tags = append(v.Tags, t.Name)
	}
	return v, nil
}

// when is the eyebrow's date: the season for TV anime ("spring 2011"), else
// the year.
func when(medium, format string, start time.Time, precision string) string {
	if medium == "anime" && format == "tv" && precision != "year" {
		return season(start.Month()) + " " + start.Format("2006")
	}
	return start.Format("2006")
}

func season(m time.Month) string {
	switch {
	case m <= time.March:
		return "winter"
	case m <= time.June:
		return "spring"
	case m <= time.September:
		return "summer"
	}
	return "fall"
}

// startLabel is "aired", "published" or "released", by medium and format.
func startLabel(medium, format string) string {
	switch {
	case format == "movie":
		return "released"
	case medium == "anime" || medium == "tv":
		return "aired"
	case medium == "manga" || medium == "comic" || medium == "book":
		return "published"
	}
	return "released"
}

// describe is the meta description: "Detective Conan, an anime from 1996."
func describe(v views.Entry) string {
	what := v.Medium
	if v.Format != "" {
		what = v.Medium + " " + v.Format
	}
	s := fmt.Sprintf("%s, %s %s", v.Title, article(what), what)
	if v.Year != "" {
		s += " from " + v.Year
	}
	return s + "."
}

func article(word string) string {
	// "uni" sounds like "you": a universe
	if strings.HasPrefix(word, "uni") {
		return "a"
	}
	switch word[0] {
	case 'a', 'e', 'i', 'o', 'u':
		return "an"
	}
	return "a"
}

// lang is a native title's HTML lang: its own, else guessed from country.
func lang(language *string, country string) string {
	if language != nil {
		return *language
	}
	switch country {
	case "JP":
		return "ja"
	case "KR":
		return "ko"
	case "CN", "TW", "HK":
		return "zh"
	}
	return ""
}

// released prints a start date only as precisely as it's known.
func released(d time.Time, precision string) string {
	switch precision {
	case "day":
		return d.Format("Jan 2, 2006")
	case "month":
		return d.Format("Jan 2006")
	default:
		return d.Format("2006")
	}
}

// withYear is "anime · 2011", or just "anime" without a date.
func withYear(label string, start *time.Time) string {
	if start == nil {
		return label
	}
	return label + " · " + start.Format("2006")
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
