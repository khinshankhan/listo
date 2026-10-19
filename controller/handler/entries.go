package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/views"
)

// EntryFinder is the slice of the generated queries the entry page needs.
type EntryFinder interface {
	EntryByPublicID(ctx context.Context, publicID int64) (sqlcgen.EntryByPublicIDRow, error)
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

	e, err := p.entries.EntryByPublicID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		// also a group's id under /entries: the query only finds entries
		p.NotFound(w, r)
		return
	}
	if err != nil {
		logx.Default().Named("entries").Error("finding entry",
			fields.Int64("public_id", id),
			fields.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	head := views.Head{
		Title:       fmt.Sprintf("%s | listo", e.Title),
		Description: fmt.Sprintf("%s on listo.", e.Title),
		Canonical:   canonical,
	}
	p.render(w, r, http.StatusOK, head, views.EntryPage(entryView(e)))
}

func entryView(e sqlcgen.EntryByPublicIDRow) views.Entry {
	v := views.Entry{Title: e.Title, Format: e.FormatLabel}
	if e.StartDate != nil && e.StartPrecision != nil {
		v.Released = released(*e.StartDate, *e.StartPrecision)
	}
	if e.UnitTotal != nil && e.UnitKind != "none" {
		v.Units = units(int(*e.UnitTotal), e.UnitKind)
	}
	if e.Country != nil {
		v.Country = *e.Country
	}
	return v
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
