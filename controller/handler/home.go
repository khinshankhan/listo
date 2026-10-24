package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/views"
)

// the guides the home page previews, by public id, and why the franchise
// needs one
const (
	guideFranchise = 2
	guideBlurb     = "Steins;Gate 0 doesn't continue the ending most people know. It follows a timeline that splits off late in the original, and you need the alternate episode 23β before you start it."
	guideGroup     = 103
	guideOrdering  = "first-timer order"

	// guideSteps and guideRows are how much of each guide a preview shows
	guideSteps = 5
	guideRows  = 3
)

// Home is what listo is, the search, and a preview of the guides. A guide
// that isn't in the catalog is left out.
func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	counts, err := p.catalog.MediumCounts(ctx)
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	var v views.Home
	for _, c := range counts {
		chip := views.MediumChip{Label: c.Label}
		if c.Entries > 0 {
			chip.Href = "/entries?medium=" + url.QueryEscape(c.Medium)
		}
		v.Media = append(v.Media, chip)
	}

	if v.Franchise, err = p.franchiseGuide(ctx); err != nil {
		p.serverError(w, r, err)
		return
	}
	if v.Ordering, err = p.orderingGuide(ctx); err != nil {
		p.serverError(w, r, err)
		return
	}

	head := views.Head{Title: "listo", Description: "Keep one list for everything you watch, read and play, and find out where to start."}
	// top bar only, and the hero has the search
	p.render(w, r, http.StatusOK, head, views.Chrome{HideSearch: true}, views.HomePage(v))
}

// franchiseGuide is the featured franchise's first two continuities.
func (p *Pages) franchiseGuide(ctx context.Context) (*views.FranchiseGuide, error) {
	g, err := p.catalog.GroupByPublicID(ctx, guideFranchise)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("guide franchise: %w", err)
	}
	rows, err := p.catalog.GroupEntries(ctx, g.ID)
	if err != nil {
		return nil, fmt.Errorf("guide entries: %w", err)
	}

	guide := &views.FranchiseGuide{Name: g.Name, Href: GroupPath(g.PublicID, g.Slug), Blurb: guideBlurb, Stats: entryStats(rows)}
	for _, card := range continuityCards(rows) {
		if len(guide.Continuities) == 2 {
			break
		}
		col := views.GuideColumn{Name: card.Name, Standing: card.Standing}
		seen := map[int64]bool{}
		for _, r := range rows {
			if GroupPath(r.BranchPublicID, r.BranchSlug) != card.Href || seen[r.PublicID] || len(col.Rows) == guideRows {
				continue
			}
			seen[r.PublicID] = true
			col.Rows = append(col.Rows, views.GuideRow{Title: r.Title, Href: EntryPath(r.PublicID, r.Slug), Meta: withYear(r.FormatLabel, r.StartDate)})
		}
		guide.Continuities = append(guide.Continuities, col)
	}
	return guide, nil
}

// orderingGuide is the featured ordering's first few steps.
func (p *Pages) orderingGuide(ctx context.Context) (*views.OrderingGuide, error) {
	g, err := p.catalog.GroupByPublicID(ctx, guideGroup)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("guide group: %w", err)
	}
	rows, err := p.catalog.GroupOrderings(ctx, g.ID)
	if err != nil {
		return nil, fmt.Errorf("guide orderings: %w", err)
	}

	for _, o := range orderings(rows, GroupPath(g.PublicID, g.Slug)) {
		if o.Name != guideOrdering {
			continue
		}
		guide := &views.OrderingGuide{Name: g.Name, Href: o.Href, Total: len(o.Positions),
			Stats: []string{plural(len(o.Positions), "step", "steps")}}
		if n, rough := episodes(rows, o.Name); n > 0 {
			s := plural(n, "episode", "episodes")
			if rough {
				s = "~" + s
			}
			guide.Stats = append(guide.Stats, s)
		}
		for _, pos := range o.Positions[:min(guideSteps, len(o.Positions))] {
			// a fork shows its first way
			st := pos.Choices[0].Steps[0]
			need := st.Necessity
			if need == "" {
				need = "required"
			}
			guide.Steps = append(guide.Steps, views.GuideStep{Number: pos.Ordinal, Title: st.Title, Href: st.Href, Note: st.Note, Necessity: need})
		}
		return guide, nil
	}
	return nil, nil
}
