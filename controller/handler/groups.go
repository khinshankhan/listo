package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// Group serves /groups/{id} and /groups/{id}/{slug}, the same way Entry does.
func (p *Pages) Group(w http.ResponseWriter, r *http.Request) {
	id, ok := publicID(r)
	if !ok {
		p.NotFound(w, r)
		return
	}

	ctx := r.Context()
	g, err := p.catalog.GroupByPublicID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		p.elsewhere(w, r, id)
		return
	}
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	canonical := GroupPath(g.PublicID, g.Slug)
	if redirectCanonical(w, r, canonical, g.Slug) {
		return
	}

	p.groupPage(w, r, g, canonical)
}

// groupPage renders any kind of group, with its sections in the sidebar.
func (p *Pages) groupPage(w http.ResponseWriter, r *http.Request, g sqlcgen.GroupByPublicIDRow, canonical string) {
	v, anchors, err := p.franchiseView(r.Context(), g, r.URL.Query().Get("by"))
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	chrome, err := p.chrome(r, nil, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	chrome.Sidebar.OnPage = anchors
	for i := range chrome.Sidebar.Franchises {
		chrome.Sidebar.Franchises[i].Active = chrome.Sidebar.Franchises[i].Href == canonical
	}

	description := fmt.Sprintf("%s, %s %s on listo", g.Name, article(g.Kind), g.Kind)
	if len(v.Stats) > 0 {
		description += ": " + strings.Join(v.Stats, ", ")
	}
	head := views.Head{Title: fmt.Sprintf("%s | listo", g.Name), Description: description + ".", Canonical: canonical}
	p.render(w, r, http.StatusOK, head, chrome, views.FranchisePage(v))
}

// orderings folds the step rows back into orderings, positions, choices and
// steps. Rows arrive sorted by ordering, position, branch, step.
// orderings folds a group's ordering rows into paths. group is the group's
// page, which each ordering's lives under.
func orderings(rows []sqlcgen.GroupOrderingsRow, group string) []views.Ordering {
	var out []views.Ordering
	var section string
	for i, row := range rows {
		newOrdering := i == 0 || row.OrderingID != rows[i-1].OrderingID
		if newOrdering {
			out = append(out, views.Ordering{
				Kind:      row.Kind,
				Name:      row.Name,
				Rationale: deref(row.Rationale),
				Origin:    row.Origin,
				ID:        row.OrderingID,
				Href:      OrderingPath(group, row.OrderingID, row.Name),
			})
			section = ""
		}
		o := &out[len(out)-1]

		if newOrdering || row.Position != rows[i-1].Position {
			pos := views.Position{Ordinal: int(row.Position), Number: fmt.Sprintf("%02d", row.Position)}
			// heading only where the section changes
			if s := deref(row.Section); s != section {
				pos.Section, section = s, s
			}
			o.Positions = append(o.Positions, pos)
		}
		pos := &o.Positions[len(o.Positions)-1]

		if len(pos.Choices) == 0 || row.Branch != rows[i-1].Branch {
			pos.Choices = append(pos.Choices, views.Choice{})
		}
		choice := &pos.Choices[len(pos.Choices)-1]

		step := views.Step{
			Title:  row.Title,
			Href:   EntryPath(row.PublicID, row.Slug),
			Units:  unitRange(row.UnitFrom, row.UnitTo),
			Note:   deref(row.Note),
			Medium: row.FormatLabel,
			Detail: stepDetail(row),
			Skip:   row.Necessity == "skip_first_time",
		}
		if row.Necessity != "required" {
			step.Necessity = necessity(row.Necessity)
		}
		choice.Steps = append(choice.Steps, step)
	}
	return out
}

// stepDetail is when the entry came out and how much of it the step covers:
// "spring 2011, Ep. 1 to Ep. 22", "2013".
func stepDetail(row sqlcgen.GroupOrderingsRow) string {
	var parts []string
	if row.StartDate != nil && row.StartPrecision != nil {
		parts = append(parts, when(row.Medium, row.Format, *row.StartDate, *row.StartPrecision))
	}
	if units := unitRange(row.UnitFrom, row.UnitTo); units != "" {
		parts = append(parts, units)
	} else if row.UnitTotal != nil && row.UnitKind != "none" {
		parts = append(parts, plural(int(*row.UnitTotal), row.UnitKind, row.UnitKind+"s"))
	}
	return strings.Join(parts, " · ")
}

// unitRange is "Ep. 1 to Ep. 22", "23β" when it starts and ends on one unit,
// or empty for the whole entry.
func unitRange(from, to *string) string {
	switch {
	case from == nil && to == nil:
		return ""
	case from == nil:
		return "to " + *to
	case to == nil:
		return "from " + *from
	case *from == *to:
		return *from
	}
	return *from + " to " + *to
}

func necessity(n string) string {
	if n == "skip_first_time" {
		return "skip the first time"
	}
	return n
}
