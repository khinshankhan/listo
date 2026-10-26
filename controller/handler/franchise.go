package handler

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// the ways a franchise's entries can be grouped; the first is the default
var franchiseSorts = []struct{ key, label string }{
	{"continuity", "by continuity"},
	{"medium", "by medium"},
	{"release", "by release"},
}

// standingRank orders continuities: the one to start with first
var standingRank = map[string]int{"primary": 0, "active": 1, "concluded": 2, "disputed": 3, "abandoned": 4}

// franchiseView is a franchise page, with links to its own sections. by is
// the ?by= grouping; anything unknown is the default. reveal is
// ?spoilers=show, which the sorting keeps.
func (p *Pages) franchiseView(ctx context.Context, g sqlcgen.GroupByPublicIDRow, by string, reveal bool) (views.Franchise, error) {
	canonical := GroupPath(g.PublicID, g.Slug)
	v := views.Franchise{Name: g.Name, Standing: deref(g.Standing)}

	parts, err := p.nodeParts(ctx, g.ID, "")
	if err != nil {
		return v, err
	}
	v.Native, v.NativeLang = parts.native, parts.nativeLang
	v.Titles, v.Relations, v.Links = parts.titles, parts.relations, parts.links

	parents, err := p.catalog.NodeGroups(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("groups: %w", err)
	}
	for i := len(parents) - 1; i >= 0; i-- {
		v.Trail = append(v.Trail, views.Crumb{Kind: parents[i].Kind, Name: parents[i].Name, Href: GroupPath(parents[i].PublicID, parents[i].Slug)})
	}
	v.Trail = append(v.Trail, views.Crumb{Kind: g.Kind, Name: g.Name})

	rows, err := p.catalog.GroupEntries(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("entries: %w", err)
	}
	v.Continuities = continuityCards(rows)
	rels, err := p.catalog.GroupRelations(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("relations: %w", err)
	}
	v.Timeline = timeline(rows, rels)
	if v.Timeline != nil {
		v.Timeline.Revealed = reveal
	}

	v.Stats = entryStats(rows)
	if len(v.Continuities) > 0 {
		v.Stats = append(v.Stats, plural(len(v.Continuities), "continuity", "continuities"))
	}

	// with no sub-groups there's nothing to group by, so release order leads
	sorts := franchiseSorts
	if !hasBranches(rows) {
		sorts = []struct{ key, label string }{{"release", "by release"}, {"medium", "by medium"}}
	}
	if !slices.ContainsFunc(sorts, func(s struct{ key, label string }) bool { return s.key == by }) {
		by = sorts[0].key
	}
	for i, s := range sorts {
		// nothing to sort
		if len(rows) == 0 {
			break
		}
		// the first is the plain address
		q := url.Values{}
		if i > 0 {
			q.Set("by", s.key)
		}
		if reveal {
			q.Set("spoilers", "show")
		}
		link := views.SortLink{Label: s.label, Href: canonical, Active: s.key == by}
		if len(q) > 0 {
			link.Href += "?" + q.Encode()
		}
		if i == 0 && s.key == "continuity" && len(v.Continuities) == 0 {
			link.Label = "by group"
		}
		v.Sorts = append(v.Sorts, link)
	}
	v.Sections = entrySections(rows, by, g.Name)

	steps, err := p.catalog.GroupOrderings(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("orderings: %w", err)
	}
	v.Orderings = orderings(steps, canonical)
	if len(v.Orderings) > 0 {
		v.StartHere = "#start-here"
	}

	siblings, err := p.catalog.GroupSiblings(ctx, g.ID)
	if err != nil {
		return v, fmt.Errorf("siblings: %w", err)
	}
	if len(siblings) > 0 && len(parents) > 0 {
		parent := parents[0]
		v.Elsewhere = parent.Name
		if parent.Kind == "universe" {
			v.ElsewhereAside = "shared setting, separate series"
		}
		for _, s := range siblings {
			v.Siblings = append(v.Siblings, views.SiblingCard{
				Name:    s.Name,
				Href:    GroupPath(s.PublicID, s.Slug),
				Meta:    s.Kind + " · " + plural(int(s.Entries), "entry", "entries"),
				Reading: "same " + parent.Kind,
			})
		}
	}

	if len(v.Continuities) > 0 {
		v.OnPage = append(v.OnPage, views.Anchor{Label: "continuities", Href: "#continuities"})
	}
	if v.Timeline != nil {
		v.OnPage = append(v.OnPage, views.Anchor{Label: "timeline", Href: "#timeline"})
	}
	v.OnPage = append(v.OnPage, views.Anchor{Label: "entries", Href: "#entries"})
	if len(v.Orderings) > 0 {
		v.OnPage = append(v.OnPage, views.Anchor{Label: "start here", Href: "#start-here"})
	}
	if len(v.Siblings) > 0 {
		v.OnPage = append(v.OnPage, views.Anchor{Label: "related", Href: "#related"})
	}
	return v, nil
}

// entryStats is "6 entries", "2 media": each entry once, however many
// branches it sits under.
func entryStats(rows []sqlcgen.GroupEntriesRow) []string {
	if len(rows) == 0 {
		return nil
	}
	entries, media := map[int64]bool{}, map[string]bool{}
	for _, r := range rows {
		entries[r.PublicID] = true
		media[r.Medium] = true
	}
	return []string{plural(len(entries), "entry", "entries"), plural(len(media), "medium", "media")}
}

// hasBranches is whether any entry sits under a sub-group rather than in the
// group itself.
func hasBranches(rows []sqlcgen.GroupEntriesRow) bool {
	for _, r := range rows {
		if r.BranchType != "entry" {
			return true
		}
	}
	return false
}

// continuityCards is a card per continuity the entries sit under, the one to
// start with first. An entry counts once per continuity.
func continuityCards(rows []sqlcgen.GroupEntriesRow) []views.ContinuityCard {
	cards := map[int64]*views.ContinuityCard{}
	seen := map[[2]int64]bool{}
	for _, r := range rows {
		if deref(r.BranchKind) != "continuity" {
			continue
		}
		c, ok := cards[r.BranchPublicID]
		if !ok {
			c = &views.ContinuityCard{Name: deref(r.BranchName), Href: GroupPath(r.BranchPublicID, r.BranchSlug), Standing: deref(r.BranchStanding)}
			cards[r.BranchPublicID] = c
		}
		if key := [2]int64{r.BranchPublicID, r.PublicID}; !seen[key] {
			seen[key] = true
			c.Entries++
		}
	}
	out := make([]views.ContinuityCard, 0, len(cards))
	for _, c := range cards {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool {
		ri, rj := standingRank[out[i].Standing], standingRank[out[j].Standing]
		if ri != rj {
			return ri < rj
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// entrySections groups the entries for the page. Rows arrive in release
// order, which each section keeps.
func entrySections(rows []sqlcgen.GroupEntriesRow, by, franchise string) []views.EntrySection {
	type section struct {
		views.EntrySection
		seen       map[int64]bool
		continuity bool
		order      int
	}
	var sections []*section
	byKey := map[string]*section{}
	add := func(key string, s *section, r sqlcgen.GroupEntriesRow) {
		got, ok := byKey[key]
		if !ok {
			s.seen = map[int64]bool{}
			s.order = len(sections)
			byKey[key] = s
			sections = append(sections, s)
			got = s
		}
		if !got.seen[r.PublicID] {
			got.seen[r.PublicID] = true
			got.Rows = append(got.Rows, entryRow(r))
		}
	}

	for _, r := range rows {
		switch by {
		case "medium":
			add(r.Medium, &section{EntrySection: views.EntrySection{Label: r.MediumLabel}}, r)
		case "release":
			add("", &section{}, r)
		default:
			if r.BranchType == "entry" {
				// directly in the franchise: its own section, last
				add("loose", &section{EntrySection: views.EntrySection{Label: "directly in " + franchise}, order: 1 << 30}, r)
				continue
			}
			add(fmt.Sprint(r.BranchPublicID), &section{
				EntrySection: views.EntrySection{
					Label:    deref(r.BranchName),
					Href:     GroupPath(r.BranchPublicID, r.BranchSlug),
					Standing: deref(r.BranchStanding),
				},
				continuity: deref(r.BranchKind) == "continuity",
			}, r)
		}
	}

	if by != "medium" && by != "release" {
		// continuities first, the one to start with leading, then other
		// groups by name, then anything directly in the franchise
		sort.SliceStable(sections, func(i, j int) bool {
			a, b := sections[i], sections[j]
			if a.order == 1<<30 || b.order == 1<<30 {
				return a.order < b.order
			}
			if a.continuity != b.continuity {
				return a.continuity
			}
			if a.continuity {
				if ra, rb := standingRank[a.Standing], standingRank[b.Standing]; ra != rb {
					return ra < rb
				}
			}
			return a.Label < b.Label
		})
	}

	out := make([]views.EntrySection, 0, len(sections))
	for _, s := range sections {
		out = append(out, s.EntrySection)
	}
	return out
}

func entryRow(r sqlcgen.GroupEntriesRow) views.EntryRow {
	var detail []string
	if r.StartDate != nil && r.StartPrecision != nil {
		detail = append(detail, when(r.Medium, r.Format, *r.StartDate, *r.StartPrecision))
	}
	if r.UnitTotal != nil && r.UnitKind != "none" {
		detail = append(detail, plural(int(*r.UnitTotal), r.UnitKind, r.UnitKind+"s"))
	}
	return views.EntryRow{
		Position: deref(r.PositionLabel),
		Title:    r.Title,
		Href:     EntryPath(r.PublicID, r.Slug),
		Role:     deref(r.Role),
		// the format's label reads better here: an anime movie is "film"
		Medium:     r.FormatLabel,
		FormatHref: formatPath(r.Medium, r.Format),
		Detail:     strings.Join(detail, " · "),
	}
}

// plural is "1 entry", "6 entries".
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}
