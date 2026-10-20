package handler

import (
	"context"
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

	if g.Kind == "franchise" {
		p.franchise(w, r, g, canonical)
		return
	}

	v, err := p.groupView(ctx, g)
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	head := views.Head{
		Title:       fmt.Sprintf("%s | listo", g.Name),
		Description: fmt.Sprintf("%s, %s %s on listo.", g.Name, article(g.Kind), g.Kind),
		Canonical:   canonical,
	}
	chrome, err := p.chrome(r, nil, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	p.render(w, r, http.StatusOK, head, chrome, views.GroupPage(v))
}

// franchise renders a franchise's own page, with its sections in the sidebar.
func (p *Pages) franchise(w http.ResponseWriter, r *http.Request, g sqlcgen.GroupByPublicIDRow, canonical string) {
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

	head := views.Head{
		Title:       fmt.Sprintf("%s | listo", g.Name),
		Description: fmt.Sprintf("%s, a franchise on listo: %s.", g.Name, strings.Join(v.Stats, ", ")),
		Canonical:   canonical,
	}
	p.render(w, r, http.StatusOK, head, chrome, views.FranchisePage(v))
}

func (p *Pages) groupView(ctx context.Context, g sqlcgen.GroupByPublicIDRow) (views.Group, error) {
	v := views.Group{Name: g.Name, Kind: g.Kind, Standing: deref(g.Standing)}

	members, err := p.catalog.GroupMembers(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("members: %w", err)
	}
	for _, m := range members {
		member := views.Member{
			Position: deref(m.PositionLabel),
			Title:    m.Title,
			Href:     NodePath(m.Type, m.PublicID, m.Slug),
			Meta:     deref(m.GroupKind),
		}
		if m.Type == "entry" {
			member.Meta = withYear(deref(m.MediumLabel), m.StartDate)
		}
		v.Members = append(v.Members, member)
	}

	steps, err := p.catalog.GroupOrderings(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("orderings: %w", err)
	}
	v.Orderings = orderings(steps)

	parts, err := p.nodeParts(ctx, g.ID, "")
	if err != nil {
		return v, err
	}
	v.Native, v.NativeLang = parts.native, parts.nativeLang
	v.Titles, v.Relations, v.Groups, v.Links = parts.titles, parts.relations, parts.groups, parts.links
	return v, nil
}

// orderings folds the step rows back into orderings, positions, choices and
// steps. Rows arrive sorted by ordering, position, branch, step.
func orderings(rows []sqlcgen.GroupOrderingsRow) []views.Ordering {
	var out []views.Ordering
	var section string
	for i, row := range rows {
		newOrdering := i == 0 || row.OrderingID != rows[i-1].OrderingID
		if newOrdering {
			out = append(out, views.Ordering{
				Kind:      row.Kind + " order",
				Name:      row.Name,
				Rationale: deref(row.Rationale),
			})
			section = ""
		}
		o := &out[len(out)-1]

		if newOrdering || row.Position != rows[i-1].Position {
			pos := views.Position{Number: fmt.Sprintf("%02d", row.Position)}
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
			Title: row.Title,
			Href:  EntryPath(row.PublicID, row.Slug),
			Units: unitRange(row.UnitFrom, row.UnitTo),
			Note:  deref(row.Note),
		}
		if row.Necessity != "required" {
			step.Necessity = necessity(row.Necessity)
		}
		choice.Steps = append(choice.Steps, step)
	}
	return out
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
