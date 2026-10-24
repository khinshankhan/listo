package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

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

	description := describe(v)
	if v.Description != "" {
		description = clip(v.Description, metaDescription)
	}
	head := views.Head{Title: fmt.Sprintf("%s | listo", e.Title), Description: description, Canonical: canonical}
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

	v.Description = deref(e.Description)
	status := deref(e.Status)

	v.Facts = append(v.Facts, views.Fact{Label: "format", Value: e.FormatLabel})
	if e.UnitKind != "none" {
		if count := units(e.UnitTotal, e.UnitAvailable, status); count != "" {
			v.Facts = append(v.Facts, views.Fact{Label: e.UnitKind + "s", Value: count})
		}
	}
	if d := duration(e.UnitMinutes, e.LengthMinutes); d != "" {
		v.Facts = append(v.Facts, views.Fact{Label: "duration", Value: d})
	}
	if status != "" {
		v.Facts = append(v.Facts, views.Fact{Label: "status", Value: status})
	}

	companies, err := p.catalog.EntryCompanies(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("companies: %w", err)
	}
	// one line per role, studios first as the query orders them
	var roles []string
	byRole := map[string][]string{}
	for _, c := range companies {
		if len(byRole[c.Role]) == 0 {
			roles = append(roles, c.Role)
		}
		byRole[c.Role] = append(byRole[c.Role], c.Name)
	}
	for _, role := range roles {
		v.Facts = append(v.Facts, views.Fact{Label: role, Value: strings.Join(byRole[role], ", ")})
	}

	if e.StartDate != nil && e.StartPrecision != nil {
		v.Year = e.StartDate.Format("2006")
		v.Eyebrow = append(v.Eyebrow, when(e.Medium, e.Format, *e.StartDate, *e.StartPrecision))
		v.Facts = append(v.Facts, views.Fact{Label: startLabel(e.Medium, e.Format), Value: released(*e.StartDate, *e.StartPrecision)})
	}
	if e.EndDate != nil && e.EndPrecision != nil {
		v.Facts = append(v.Facts, views.Fact{Label: "ended", Value: released(*e.EndDate, *e.EndPrecision)})
	}
	country := deref(e.Country)
	if country != "" {
		v.Facts = append(v.Facts, views.Fact{Label: "country", Value: country})
	}

	parts, err := p.nodeParts(ctx, e.ID, country)
	if err != nil {
		return v, err
	}
	// an adaptation_of relation names the source better than its type does
	if source := deref(e.SourceType); source != "" && !slices.ContainsFunc(parts.relations, func(r views.Related) bool { return r.Reading == "adaptation of" }) {
		v.Facts = append(v.Facts, views.Fact{Label: "source", Value: strings.ReplaceAll(source, "_", " ")})
	}

	unitRows, err := p.catalog.EntryUnits(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("units: %w", err)
	}
	v.Next, v.Episodes = episodeList(unitRows, status, time.Now())
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

// metaDescription is about as much as a search result shows.
const metaDescription = 160

// clip cuts text to at most n runes, at a word, with an ellipsis.
func clip(text string, n int) string {
	runes := []rune(text)
	if len(runes) <= n {
		return text
	}
	cut := string(runes[:n-1])
	// back up to a space only if the cut split a word
	if runes[n-1] != ' ' {
		if i := strings.LastIndexByte(cut, ' '); i > 0 {
			cut = cut[:i]
		}
	}
	return strings.TrimRight(cut, " ,.;:") + "…"
}

// units is how many units there are: "28", or "19 of 28" while releasing.
func units(total *int32, available pgtype.Numeric, status string) string {
	out := ""
	if available.Valid && status == "releasing" {
		if f, err := available.Float64Value(); err == nil && f.Valid {
			out = strconv.FormatFloat(f.Float64, 'f', -1, 64)
			if total == nil {
				return out + " so far"
			}
			return out + " of " + fmt.Sprint(*total)
		}
	}
	if total != nil {
		out = fmt.Sprint(*total)
	}
	return out
}

// duration is an episode's length, "24 min", or else the whole work's,
// "1 hr 30 min".
func duration(perUnit, whole *int32) string {
	m := perUnit
	if m == nil {
		m = whole
	}
	if m == nil {
		return ""
	}
	hours, mins := *m/60, *m%60
	switch {
	case hours == 0:
		return fmt.Sprintf("%d min", mins)
	case mins == 0:
		return fmt.Sprintf("%d hr", hours)
	}
	return fmt.Sprintf("%d hr %d min", hours, mins)
}

// episodeList is what airs next and the list to show. Nothing is listed unless
// some unit has a title or an air date: a column of bare numbers says nothing
// the count doesn't. While releasing, the list runs newest first from what's
// next; otherwise in order.
func episodeList(rows []sqlcgen.EntryUnitsRow, status string, now time.Time) (*views.NextUnit, []views.Episode) {
	var next *views.NextUnit
	known := false
	var list []views.Episode
	for _, r := range rows {
		known = known || r.Title != nil || r.AirDate != nil
		ep := views.Episode{Number: r.Label, Title: deref(r.Title)}
		if r.AirDate != nil {
			ep.When = relative(*r.AirDate, now)
			ep.Upcoming = r.AirDate.After(now)
			if ep.Upcoming && next == nil {
				next = &views.NextUnit{Label: r.Label, In: ep.When}
			}
		}
		list = append(list, ep)
	}
	if !known {
		return nil, nil
	}
	if status == "releasing" {
		// the next one, then everything out, newest first
		var out []views.Episode
		for i := len(list) - 1; i >= 0; i-- {
			if !list[i].Upcoming || (next != nil && list[i].Number == next.Label) {
				out = append(out, list[i])
			}
		}
		list = out
	}
	return next, list
}

// relative is when something airs or aired, from now: "in 3d 4h", "in 40m",
// "2 days ago", "3 weeks ago". Past a month it's the date.
func relative(t, now time.Time) string {
	if d := t.Sub(now); d > 0 {
		days, hours, mins := int(d.Hours())/24, int(d.Hours())%24, int(d.Minutes())%60
		switch {
		case days > 0:
			return fmt.Sprintf("in %dd %dh", days, hours)
		case hours > 0:
			return fmt.Sprintf("in %dh %dm", hours, mins)
		case mins > 0:
			return fmt.Sprintf("in %dm", mins)
		}
		return "any minute"
	}
	ago := now.Sub(t)
	switch days := int(ago.Hours()) / 24; {
	case ago < time.Hour:
		return "just now"
	case days == 0:
		return plural(int(ago.Hours()), "hour", "hours") + " ago"
	case days == 1:
		return "yesterday"
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	case days < 28:
		return plural(days/7, "week", "weeks") + " ago"
	}
	return t.UTC().Format("Jan 2, 2006")
}
