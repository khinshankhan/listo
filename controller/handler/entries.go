package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/views"
)

// Catalog is the slice of the generated queries the catalog pages need.
type Catalog interface {
	EntryByPublicID(ctx context.Context, publicID int64) (sqlcgen.EntryByPublicIDRow, error)
	NodeTitles(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeTitlesRow, error)
	NodeLinks(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeLinksRow, error)
	NodeRelations(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeRelationsRow, error)
	NodeGroups(ctx context.Context, memberID uuid.UUID) ([]sqlcgen.NodeGroupsRow, error)
}

// EntryPath is an entry's canonical URL: the id is the key, the slug decoration.
func EntryPath(publicID int64, slug string) string {
	return fmt.Sprintf("/entries/%d/%s", publicID, slug)
}

// Entry serves /entries/{id}[/{slug}]. The id is the key; a wrong or missing
// slug redirects to the canonical URL.
func (p *Pages) Entry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		p.NotFound(w, r)
		return
	}

	ctx := r.Context()
	e, err := p.catalog.EntryByPublicID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		// also a group's id under /entries: the query only finds entries
		p.NotFound(w, r)
		return
	}
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	canonical := EntryPath(e.PublicID, e.Slug)
	if r.PathValue("slug") != e.Slug {
		to := canonical
		if r.URL.RawQuery != "" {
			to += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, to, http.StatusMovedPermanently)
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
	p.render(w, r, http.StatusOK, head, views.EntryPage(v))
}

func (p *Pages) serverError(w http.ResponseWriter, r *http.Request, err error) {
	logx.Default().Named("pages").Error("serving page",
		fields.String("path", r.URL.Path),
		fields.Error(err))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// entryView gathers everything the page shows into display form.
func (p *Pages) entryView(ctx context.Context, e sqlcgen.EntryByPublicIDRow) (views.Entry, error) {
	v := views.Entry{Title: e.Title, Medium: e.MediumLabel}
	if e.FormatLabel != e.MediumLabel {
		v.Format = e.FormatLabel
	}
	if e.Country != nil {
		v.Country = *e.Country
	}
	if e.StartDate != nil && e.StartPrecision != nil {
		v.Year = e.StartDate.Format("2006")
		v.Released = released(*e.StartDate, *e.StartPrecision)
	}
	if e.UnitTotal != nil && e.UnitKind != "none" {
		v.Units = units(int(*e.UnitTotal), e.UnitKind)
	}

	titles, err := p.catalog.NodeTitles(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("titles: %w", err)
	}
	for _, t := range titles {
		v.Titles = append(v.Titles, views.Title{Kind: t.Kind, Title: t.Title})
		if t.Kind == "native" && v.Native == "" {
			v.Native = t.Title
			v.NativeLang = lang(t.Language, v.Country)
		}
	}

	relations, err := p.catalog.NodeRelations(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("relations: %w", err)
	}
	for _, rel := range relations {
		related := views.Related{Reading: rel.Reading, Title: rel.Title, Meta: "group"}
		if rel.Type == "entry" {
			related.Href = EntryPath(rel.PublicID, rel.Slug)
			related.Meta = deref(rel.MediumLabel)
			if rel.StartDate != nil {
				related.Meta += " · " + rel.StartDate.Format("2006")
			}
		}
		related.Note = deref(rel.Note)
		v.Relations = append(v.Relations, related)
	}

	groups, err := p.catalog.NodeGroups(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("groups: %w", err)
	}
	for _, g := range groups {
		v.Groups = append(v.Groups, views.GroupRef{Name: g.Name, Position: deref(g.PositionLabel)})
	}

	links, err := p.catalog.NodeLinks(ctx, e.ID)
	if err != nil {
		return v, fmt.Errorf("links: %w", err)
	}
	for _, l := range links {
		v.Links = append(v.Links, views.Link{Site: l.Site, Href: l.Url, ID: l.ExternalID})
	}

	return v, nil
}

// describe is the meta description: "Detective Conan, an anime from 1996, Japan."
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

// units is "24 episodes", "1 volume".
func units(n int, kind string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", kind)
	}
	return fmt.Sprintf("%d %ss", n, kind)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
