package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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
	p.render(w, r, http.StatusOK, head, views.GroupPage(v))
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

// orderings folds the step rows back into orderings. Rows arrive sorted by
// ordering, position, branch, step; a position's first branch is numbered,
// later ones are "or".
func orderings(rows []sqlcgen.GroupOrderingsRow) []views.Ordering {
	var out []views.Ordering
	for i, row := range rows {
		if i == 0 || row.OrderingID != rows[i-1].OrderingID {
			out = append(out, views.Ordering{
				Kind:      row.Kind + " order",
				Name:      row.Name,
				Rationale: deref(row.Rationale),
			})
		}

		step := views.Step{
			Title: row.Title,
			Href:  EntryPath(row.PublicID, row.Slug),
			Units: unitRange(row.UnitFrom, row.UnitTo),
			Note:  deref(row.Note),
		}
		if row.Necessity != "required" {
			step.Necessity = necessity(row.Necessity)
		}
		switch {
		case row.Step > 1:
		case row.Branch > 1:
			step.Number = "or"
		default:
			step.Number = fmt.Sprintf("%02d", row.Position)
		}

		o := &out[len(out)-1]
		o.Steps = append(o.Steps, step)
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
