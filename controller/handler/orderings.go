package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// shownSteps is how many steps a long ordering shows before "show more".
const shownSteps = 5

var origins = map[string]string{"staff": "maintained by staff", "community": "from the community", "auto": "auto-seeded"}

// Ordering serves /groups/{id}/{slug}/orderings/{ordering}: one path through
// a group, step by step.
func (p *Pages) Ordering(w http.ResponseWriter, r *http.Request) {
	id, ok := publicID(r)
	orderingID, err := uuid.Parse(r.PathValue("ordering"))
	if !ok || err != nil {
		p.NotFound(w, r)
		return
	}

	ctx := r.Context()
	g, err := p.catalog.GroupByPublicID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		p.NotFound(w, r)
		return
	}
	if err != nil {
		p.serverError(w, r, err)
		return
	}

	group := GroupPath(g.PublicID, g.Slug)
	rows, err := p.catalog.GroupOrderings(ctx, g.ID)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	all := orderings(p.store, rows, group)
	at := -1
	for i, o := range all {
		if o.ID == orderingID {
			at = i
			break
		}
	}
	if at < 0 {
		p.NotFound(w, r)
		return
	}
	o := all[at]
	// a wrong or missing slug, the group's or the ordering's, is a redirect
	if r.URL.Path != o.Href {
		moved(w, r, o.Href)
		return
	}

	parents, err := p.catalog.NodeGroups(ctx, g.ID)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	v := orderingView(o, all, rows, g, parents)

	chrome, err := p.chrome(r, nil, "")
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	// the group and whatever holds it count as where we are
	here := map[string]bool{group: true}
	for _, parent := range parents {
		here[GroupPath(parent.PublicID, parent.Slug)] = true
	}
	for i := range chrome.Sidebar.Franchises {
		chrome.Sidebar.Franchises[i].Active = here[chrome.Sidebar.Franchises[i].Href]
	}

	head := views.Head{
		Title:       fmt.Sprintf("%s, %s | listo", o.Name, g.Name),
		Description: fmt.Sprintf("%s: a %s order for %s, %s.", o.Name, o.Kind, g.Name, plural(len(o.Positions), "step", "steps")),
		Canonical:   o.Href,
	}
	p.render(w, r, http.StatusOK, head, chrome, views.OrderingPage(v))
}

func orderingView(o views.Ordering, all []views.Ordering, rows []sqlcgen.GroupOrderingsRow, g sqlcgen.GroupByPublicIDRow, parents []sqlcgen.NodeGroupsRow) views.OrderingDetail {
	v := views.OrderingDetail{
		Name:      o.Name,
		Origin:    origins[o.Origin],
		Rationale: o.Rationale,
		Positions: o.Positions,
		Shown:     len(o.Positions),
	}
	// folding a couple of steps away saves nothing
	if len(o.Positions) > shownSteps+2 {
		v.Shown = shownSteps
	}

	for i := len(parents) - 1; i >= 0; i-- {
		v.Trail = append(v.Trail, views.Crumb{Kind: parents[i].Kind, Name: parents[i].Name, Href: GroupPath(parents[i].PublicID, parents[i].Slug)})
	}
	v.Trail = append(v.Trail, views.Crumb{Kind: g.Kind, Name: g.Name, Href: GroupPath(g.PublicID, g.Slug)}, views.Crumb{Kind: "ordering"})

	required := 0
	for _, pos := range o.Positions {
		if !pos.Skip() && pos.Choices[0].Steps[0].Necessity == "" {
			required++
		}
	}
	v.Stats = []string{plural(len(o.Positions), "step", "steps"), fmt.Sprintf("%d required", required)}
	if n, rough := episodes(rows, o.Name); n > 0 {
		s := plural(n, "episode", "episodes")
		if rough {
			s = "~" + s
		}
		v.Stats = append(v.Stats, s)
	}

	// one ordering has nothing to switch to
	if len(all) > 1 {
		for _, other := range all {
			v.Orderings = append(v.Orderings, views.OrderingLink{
				Name:   other.Name,
				Meta:   other.Kind + " · " + origin(other.Origin),
				Href:   other.Href,
				Active: other.Href == o.Href,
			})
		}
	}
	v.Alternative = alternative(o, all)
	return v
}

// episodes totals the episodes a first watch of an ordering covers, taking
// the first way through each fork. It's rough when there's a fork or a count
// is missing.
func episodes(rows []sqlcgen.GroupOrderingsRow, ordering string) (int, bool) {
	total, rough := 0, false
	for _, r := range rows {
		if r.Name != ordering {
			continue
		}
		if r.Branch > 1 {
			rough = true
			continue
		}
		if r.UnitKind != "episode" || r.Necessity == "skip_first_time" {
			continue
		}
		last := r.UnitTotal
		if r.ToOrdinal != nil {
			last = r.ToOrdinal
		}
		if last == nil {
			rough = true
			continue
		}
		first := int32(1)
		if r.FromOrdinal != nil {
			first = *r.FromOrdinal
		}
		total += int(*last - first + 1)
	}
	return total, rough
}

// alternative points a first-timer's ordering at the one for people who
// know the story, and every other ordering back at the recommended one.
func alternative(o views.Ordering, all []views.Ordering) *views.Alternative {
	find := func(kinds ...string) *views.Ordering {
		for _, kind := range kinds {
			for i := range all {
				if all[i].Kind == kind && all[i].Href != o.Href {
					return &all[i]
				}
			}
		}
		return nil
	}

	if o.Kind == "recommended" {
		other := find("chronological", "release")
		if other == nil {
			return nil
		}
		body := fmt.Sprintf("%s follows the story's own timeline.", other.Name)
		if other.Kind == "release" {
			body = fmt.Sprintf("%s goes in the order they came out.", other.Name)
		}
		return &views.Alternative{Title: "rewatching, or already know the story?", Body: body, Label: "view " + other.Name, Href: other.Href}
	}
	other := find("recommended")
	if other == nil {
		return nil
	}
	return &views.Alternative{
		Title: "watching for the first time?",
		Body:  fmt.Sprintf("%s is the recommended way in.", other.Name),
		Label: "view " + other.Name,
		Href:  other.Href,
	}
}

func origin(o string) string {
	if o == "auto" {
		return "auto-seeded"
	}
	return o
}
