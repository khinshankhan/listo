package handler

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// the timeline's drawing, in SVG units
const (
	tlWidth = 960
	// the first year starts at tlLeft and the last ends at tlRight
	tlLeft, tlRight = 56, 904
	// names keep this far from the edge, and labels under dots this far apart
	tlEdge, tlLabelGap = 24, 6
	tlRow              = 88
	// between two works on a row
	tlGap = 40
	// a work's card: this much either side of its stretch of row, from this
	// far above its dots, this tall
	tlPad, tlAbove, tlCard = 10, 32, 58
	// links to other rows hang off this far into a card, where the cards after
	// it in time can't cover them
	tlTrunk = 8
	// a step line leaves and meets a card this straight, turns this round, and
	// runs across this far short of a card it meets from above or below
	tlOffset, tlRadius, tlTurn = 14, 10, 10
	// past this many works linked to nothing in a box, they share one card
	tlLoneMax = 12
)

// mergeKinds make two entries one work: the same story in another medium or
// edition.
var mergeKinds = map[string]bool{"adaptation_of": true, "remaster_of": true, "recap_of": true, "alternative_version_of": true}

// linkNames are what a link's line says, for every kind drawn but a sequel.
var linkNames = map[string]string{
	"side_story_of": "side story", "spin_off_of": "spin-off", "remake_of": "remake", "reboot_of": "reboot",
	"crossover_with": "crossover", "same_setting_as": "same setting", "references": "reference",
}

// shortFormats fit under a dot.
var shortFormats = map[string]string{"visual novel": "VN", "light novel": "LN"}

// tlWork is a work being laid out: its releases in order, the part of the
// group it sits in, and where it went.
type tlWork struct {
	entries []sqlcgen.GroupEntriesRow
	// branch is the group it's boxed in, or 0 when it sits straight in this one
	branch     int64
	branchName string
	// attach is the work it hangs off, if any; follows when it's a sequel,
	// which carries on that work's row
	attach      *tlWork
	follows     bool
	left, right float64
	view        views.TimelineWork
	// lone is how many works linked to nothing this card stands for, if any
	lone int
}

// timeline draws a group's entries by year, or is nil when there are fewer
// than two works to connect.
func timeline(rows []sqlcgen.GroupEntriesRow, rels []sqlcgen.GroupRelationsRow) *views.Timeline {
	tl := &views.Timeline{Width: tlWidth}

	// an entry under two branches comes once per branch; the first will do
	byID := map[int64]sqlcgen.GroupEntriesRow{}
	var dated []sqlcgen.GroupEntriesRow
	for _, r := range rows {
		if _, seen := byID[r.PublicID]; seen {
			continue
		}
		byID[r.PublicID] = r
		if r.StartDate == nil {
			tl.Undated++
			continue
		}
		dated = append(dated, r)
	}

	// adaptations and editions join their source
	parent := map[int64]int64{}
	var find func(int64) int64
	find = func(id int64) int64 {
		if p, ok := parent[id]; ok && p != id {
			root := find(p)
			parent[id] = root
			return root
		}
		return id
	}
	for _, rel := range rels {
		_, from := byID[rel.FromID]
		_, to := byID[rel.ToID]
		if mergeKinds[rel.Kind] && from && to {
			parent[find(rel.FromID)] = find(rel.ToID)
		}
	}
	var works []*tlWork
	workOf := map[int64]*tlWork{}
	for _, r := range sortedByDate(dated) {
		root := find(r.PublicID)
		w := workOf[root]
		if w == nil {
			w = &tlWork{}
			workOf[root] = w
			works = append(works, w)
		}
		w.entries = append(w.entries, r)
		workOf[r.PublicID] = w
	}
	if len(works) < 2 {
		return nil
	}
	linked := map[*tlWork]bool{}
	for _, rel := range rels {
		if from, to := workOf[rel.FromID], workOf[rel.ToID]; from != nil && to != nil && from != to {
			linked[from], linked[to] = true, true
		}
	}
	works = foldLone(works, linked)
	// a work sits where its first release does
	names := map[string]int{}
	for _, w := range works {
		first := w.entries[0]
		if first.BranchType == "group" {
			w.branch, w.branchName = first.BranchPublicID, deref(first.BranchName)
		}
		w.view.ID, w.view.Name, w.view.Href = first.PublicID, first.Title, EntryPath(first.PublicID, first.Slug)
		if w.lone > 0 {
			// not an entry's id, so no link can mean it
			w.view.ID, w.view.Name, w.view.Href = -w.branch-1, fmt.Sprintf("%d more, linked to nothing", w.lone), "#entries"
			continue
		}
		names[first.Title]++
	}
	// in words, works that share a name are told apart by year
	called := func(e sqlcgen.GroupEntriesRow) string {
		if names[e.Title] > 1 {
			return fmt.Sprintf("%s (%d)", e.Title, e.StartDate.Year())
		}
		return e.Title
	}

	// a work hangs off the one it continues or sits beside, a sequel first
	for _, rel := range rels {
		from, to := workOf[rel.FromID], workOf[rel.ToID]
		if rel.Family != "sequence" || from == nil || to == nil || from == to {
			continue
		}
		if from.attach == nil || (rel.Kind == "sequel_of" && !from.follows) {
			from.attach, from.follows = to, rel.Kind == "sequel_of"
		}
	}

	first, last := year(works[0].entries[0]), 0.0
	for _, w := range works {
		last = max(last, year(w.entries[len(w.entries)-1]))
	}
	minYear := math.Floor(first)
	span := max(math.Floor(last)+1-minYear, 4)
	x := func(at float64) float64 { return tlLeft + (at-minYear)*(tlRight-tlLeft)/span }
	step := int(math.Ceil(span / 8))
	for y := int(minYear); float64(y) < minYear+span; y += step {
		tl.Ticks = append(tl.Ticks, views.Tick{X: int(math.Round(x(float64(y)))), Label: strconv.Itoa(y)})
	}

	for _, w := range works {
		place(w, x)
	}

	// the loose works first, then each box in the order it starts
	var order []int64
	for _, w := range works {
		if !slices.Contains(order, w.branch) {
			order = append(order, w.branch)
		}
	}
	slices.SortStableFunc(order, func(a, b int64) int {
		return cmp.Compare(min(a, 1), min(b, 1))
	})

	y := 20
	for _, branch := range order {
		// a sequel carries on its prequel's row when there's room; anything
		// else hanging off a work gets a row right under it, shared with the
		// work's other side stories while there's room; works linked to
		// nothing share rows with each other the same way
		var lanes []*lane
		laneOf := map[*tlWork]*lane{}
		var name string
		for _, w := range works {
			if w.branch != branch {
				continue
			}
			name = w.branchName
			on := laneOf[w.attach]
			var l *lane
			switch {
			case on != nil && w.follows && on.right+tlGap <= w.left:
				l = on
			case on != nil:
				for _, c := range lanes {
					if c.under == on && c.right+tlGap <= w.left {
						l = c
						break
					}
				}
				if l == nil {
					l = &lane{under: on}
					at := slices.Index(lanes, on) + 1
					for at < len(lanes) && lanes[at].within(on) {
						at++
					}
					lanes = slices.Insert(lanes, at, l)
				}
			case !linked[w]:
				for _, c := range lanes {
					if c.lone && c.right+tlGap <= w.left {
						l = c
						break
					}
				}
				if l == nil {
					l = &lane{lone: true}
					lanes = append(lanes, l)
				}
			default:
				l = &lane{}
				lanes = append(lanes, l)
			}
			l.right = max(l.right, w.right)
			laneOf[w] = l
		}

		top := y
		firstRow := y + 44
		if branch != 0 {
			firstRow = y + 84
		}
		for w, l := range laneOf {
			w.view.Y = firstRow + slices.Index(lanes, l)*tlRow
			w.view.Card = views.TimelineCard{X: int(math.Floor(w.left)) - tlPad, Y: w.view.Y - tlAbove, H: tlCard}
			w.view.Card.W = int(math.Ceil(w.right)) + tlPad - w.view.Card.X
			w.view.NameX, w.view.NameY = w.view.Card.X+tlPad, w.view.Card.Y+20
		}
		bottom := firstRow + (len(lanes)-1)*tlRow + 44
		if branch != 0 {
			tl.Boxes = append(tl.Boxes, views.TimelineBox{X: 16, Y: top, W: tlWidth - 32, H: bottom - top, Name: name})
		}
		y = bottom + 16
	}
	tl.Height = y + 24

	dotX := func(id int64) float64 {
		w := workOf[id]
		for i, e := range w.entries {
			if e.PublicID == id {
				return float64(w.view.Dots[dotIndex(w, i)].X)
			}
		}
		return 0
	}

	type drawn struct {
		from, to *tlWork
		name     string
	}
	var seen []drawn
	for _, rel := range rels {
		from, to := workOf[rel.FromID], workOf[rel.ToID]
		name, other := linkNames[rel.Kind]
		sequel := rel.Kind == "sequel_of"
		if from == nil || to == nil || from == to || (!sequel && !other) {
			continue
		}
		if slices.Contains(seen, drawn{from, to, name}) {
			continue
		}
		seen = append(seen, drawn{from, to, name})

		link := views.TimelineLink{Sequel: rel.Family == "sequence", Spoiler: rel.Spoiler, Label: name}
		line := fmt.Sprintf("%s, %s %s", called(byID[rel.FromID]), rel.Label, called(byID[rel.ToID]))
		// drawn left to right: a sequel from the work it continues, anything
		// else from whichever came first
		source, sourceX, target, targetX := to, dotX(rel.ToID), from, dotX(rel.FromID)
		if !link.Sequel {
			when := max(year(byID[rel.FromID]), year(byID[rel.ToID]))
			if rel.EstablishedOn != nil {
				when = float64(rel.EstablishedOn.Year())
			}
			link.Label = fmt.Sprintf("%s, %d", name, int(when))
			line += fmt.Sprintf(", %d", int(when))
			if dotX(rel.FromID) < dotX(rel.ToID) {
				source, sourceX, target, targetX = from, dotX(rel.FromID), to, dotX(rel.ToID)
			}
		}
		link.Source, link.Target = source.view.ID, target.view.ID
		// a run in along the target's row would pass under any card there
		// between the trunk and the target
		clear := true
		for _, w := range works {
			c := w.view.Card
			if w != target && w.view.Y == target.view.Y && c.X < target.view.Card.X && c.X+c.W > source.view.Card.X+tlTrunk {
				clear = false
			}
		}
		route(&link, source.view.Card, sourceX, target.view.Card, targetX, clear)
		tl.Spoilers = tl.Spoilers || link.Spoiler
		tl.Links = append(tl.Links, link)
		tl.Lines = append(tl.Lines, views.TimelineLine{Text: line, Spoiler: rel.Spoiler})
	}

	var lines []views.TimelineLine
	for _, w := range works {
		tl.Works = append(tl.Works, w.view)
		if w.lone > 0 {
			lines = append(lines, views.TimelineLine{Text: w.view.Name + ", listed under entries"})
			continue
		}
		var releases []string
		for _, e := range w.entries {
			releases = append(releases, fmt.Sprintf("%s %d", short(e.FormatLabel), e.StartDate.Year()))
		}
		lines = append(lines, views.TimelineLine{Text: called(w.entries[0]) + ": " + strings.Join(releases, ", ")})
	}
	tl.Lines = append(lines, tl.Lines...)
	return tl
}

// foldLone puts a box's works linked to nothing on one card when there are
// more than tlLoneMax of them: one by one they'd bury the works that lead
// somewhere, and the list under the timeline has them all anyway.
func foldLone(works []*tlWork, linked map[*tlWork]bool) []*tlWork {
	// by box, 0 for none
	box := func(w *tlWork) int64 {
		if w.entries[0].BranchType == "group" {
			return w.entries[0].BranchPublicID
		}
		return 0
	}
	lone := map[int64][]*tlWork{}
	for _, w := range works {
		if !linked[w] {
			lone[box(w)] = append(lone[box(w)], w)
		}
	}
	var out []*tlWork
	for _, w := range works {
		ws := lone[box(w)]
		if linked[w] || len(ws) <= tlLoneMax {
			out = append(out, w)
			continue
		}
		// the box's first lone work stands in for the rest
		if ws[0] != w {
			continue
		}
		all := &tlWork{lone: len(ws)}
		for _, o := range ws {
			all.entries = append(all.entries, o.entries...)
		}
		all.entries = sortedByDate(all.entries)
		out = append(out, all)
	}
	return out
}

// lane is a row being filled left to right, and the row it hangs under.
// A lone lane holds only works linked to nothing.
type lane struct {
	right float64
	under *lane
	lone  bool
}

// within reports whether l is top or hangs somewhere beneath it.
func (l *lane) within(top *lane) bool {
	for ; l != nil; l = l.under {
		if l == top {
			return true
		}
	}
	return false
}

// route draws a link from the source card to the target, leaving the dot at
// sourceX or meeting the one at targetX where it needs to. Along a row it
// runs card to card, arcing over the row unless it's a sequel. To another
// row a sequel, side story or spin-off hangs off its card's trunk and runs
// in when nothing's in the way, labelled on that run or beside the trunk; one
// whose target came out first leaves its card's left edge and goes round the
// left of both, and so of every card between, which came out later still;
// anything else, or one whose target starts too near the trunk, has cards
// before it on its row or no room for its label, drops and meets the target
// from above or below, labelled on the run across.
func route(l *views.TimelineLink, from views.TimelineCard, sourceX float64, to views.TimelineCard, targetX float64, clear bool) {
	type pt struct{ x, y float64 }
	top := func(c views.TimelineCard) float64 { return float64(c.Y) }
	bottom := func(c views.TimelineCard) float64 { return float64(c.Y + c.H) }
	row := func(c views.TimelineCard) float64 { return float64(c.Y + tlAbove) }
	trunk, left := float64(from.X+tlTrunk), float64(to.X)
	edge, meet, turn := top, bottom, float64(tlTurn)
	if from.Y < to.Y {
		edge, meet, turn = bottom, top, -tlTurn
	}
	labelW := math.Ceil(dotFont.width(l.Label)) + 12
	// right-aligned on the run in, or left of the line down when it won't
	// fit, or past the line down when that's off the edge
	onRun := func(run, down float64) bool { return run-14-labelW > down+8 }
	offEdge := func(down float64) bool { return down-6-labelW < 0 }
	beside := func(run, y, down float64) {
		l.LabelEnd, l.LabelX, l.LabelY = true, int(down-6), int(y)
		switch {
		case onRun(run, down):
			l.LabelX = int(run - 14)
		case offEdge(down):
			l.LabelEnd, l.LabelX = false, int(math.Ceil(down+8+labelW/2))
		}
	}

	var pts []pt
	switch {
	case from.Y == to.Y && l.Sequel:
		y := row(from)
		pts = []pt{{float64(from.X + from.W), y}, {left, y}}
		l.LabelX, l.LabelY = int((float64(from.X+from.W)+left)/2), int(y)
	case from.Y == to.Y:
		over := top(from) - tlOffset
		pts = []pt{{sourceX, top(from)}, {sourceX, over}, {targetX, over}, {targetX, top(to)}}
		l.LabelX, l.LabelY = int((sourceX+targetX)/2), int(over)
	case l.Sequel && clear && left > trunk+tlOffset+tlRadius && (onRun(left, trunk) || !offEdge(trunk)):
		y := row(to)
		pts = []pt{{trunk, edge(from)}, {trunk, y}, {left, y}}
		beside(left, y, trunk)
	case l.Sequel && targetX < trunk:
		gutter := max(float64(min(from.X, to.X))-tlOffset, 4)
		y := row(to)
		pts = []pt{{float64(from.X), row(from)}, {gutter, row(from)}, {gutter, y}, {left, y}}
		beside(left, y, gutter)
	default:
		x := sourceX
		if l.Sequel {
			x = trunk
		}
		across := meet(to) + turn
		pts = []pt{{x, edge(from)}, {x, across}, {targetX, across}, {targetX, meet(to)}}
		if l.Sequel {
			beside(targetX, across, trunk)
		} else {
			l.LabelX, l.LabelY = int((x+targetX)/2), int(across)
		}
	}
	if l.Label != "" {
		l.LabelW = int(labelW)
	}

	// straight runs, rounded where they turn
	var b strings.Builder
	f := func(v float64) string { return strconv.FormatFloat(math.Round(v*10)/10, 'f', -1, 64) }
	fmt.Fprintf(&b, "M%s %s", f(pts[0].x), f(pts[0].y))
	for i := 1; i < len(pts)-1; i++ {
		prev, at, next := pts[i-1], pts[i], pts[i+1]
		in, out := math.Hypot(at.x-prev.x, at.y-prev.y), math.Hypot(next.x-at.x, next.y-at.y)
		r := min(tlRadius, in/2, out/2)
		if r == 0 {
			continue
		}
		fmt.Fprintf(&b, "L%s %s", f(at.x-(at.x-prev.x)/in*r), f(at.y-(at.y-prev.y)/in*r))
		fmt.Fprintf(&b, "Q%s %s %s %s", f(at.x), f(at.y), f(at.x+(next.x-at.x)/out*r), f(at.y+(next.y-at.y)/out*r))
	}
	end := pts[len(pts)-1]
	fmt.Fprintf(&b, "L%s %s", f(end.x), f(end.y))
	l.Path = b.String()
}

// place sets a work's dots, and the stretch of its row they and its name take.
func place(w *tlWork, x func(float64) float64) {
	// a dot a year, and the formats out that year
	var formats [][]string
	for i, e := range w.entries {
		if i > 0 && e.StartDate.Year() == w.entries[i-1].StartDate.Year() {
			formats[len(formats)-1] = append(formats[len(formats)-1], short(e.FormatLabel))
			continue
		}
		w.view.Dots = append(w.view.Dots, views.TimelineDot{X: int(math.Round(x(year(e))))})
		formats = append(formats, []string{short(e.FormatLabel)})
	}

	// dots too close to label one by one share a label under them, so a run
	// of yearly games reads "game ×12", not a pile of "game"s
	type run struct{ from, to int }
	label := func(r run) string {
		var order []string
		count := map[string]int{}
		for _, fs := range formats[r.from:r.to] {
			for _, f := range fs {
				if count[f] == 0 {
					order = append(order, f)
				}
				count[f]++
			}
		}
		for i, f := range order {
			if count[f] > 1 {
				order[i] = fmt.Sprintf("%s \u00d7%d", f, count[f])
			}
		}
		return strings.Join(order, ", ")
	}
	// centred under its dots, but not past the edge
	centre := func(r run) float64 {
		half := dotFont.width(label(r)) / 2
		mid := float64(w.view.Dots[r.from].X+w.view.Dots[r.to-1].X) / 2
		return min(max(mid, tlEdge+half), tlWidth-tlEdge-half)
	}
	var runs []run
	for i := range w.view.Dots {
		runs = append(runs, run{i, i + 1})
	}
	for i := 1; i < len(runs); {
		a, b := runs[i-1], runs[i]
		if centre(a)+dotFont.width(label(a))/2+tlLabelGap > centre(b)-dotFont.width(label(b))/2 {
			runs[i-1].to = b.to
			runs = slices.Delete(runs, i, i+1)
			i = max(i-1, 1)
			continue
		}
		i++
	}
	for _, r := range runs {
		d := &w.view.Dots[r.from]
		d.Label, d.LabelX = label(r), int(math.Round(centre(r)))
	}

	firstX, lastX := float64(w.view.Dots[0].X), float64(w.view.Dots[len(w.view.Dots)-1].X)
	name := nameFont.width(w.view.Name)
	w.left = firstX - 5
	// too long to fit, it ends at its dot instead, or as near as the left
	// edge lets it
	if w.left+name > tlWidth-tlEdge {
		w.left = max(firstX+6-name, tlEdge+tlPad)
	}
	w.right = max(w.left+name, lastX)
	for _, d := range w.view.Dots {
		if d.Label != "" {
			half := dotFont.width(d.Label) / 2
			w.left, w.right = min(w.left, float64(d.LabelX)-half), max(w.right, float64(d.LabelX)+half)
		}
	}
}

// dotIndex is which of a work's dots its i'th release is drawn as.
func dotIndex(w *tlWork, i int) int {
	dot := 0
	for j := 1; j <= i; j++ {
		if w.entries[j].StartDate.Year() != w.entries[j-1].StartDate.Year() {
			dot++
		}
	}
	return dot
}

// year is when an entry came out as a year with the month as a fraction.
func year(e sqlcgen.GroupEntriesRow) float64 {
	y := float64(e.StartDate.Year())
	if deref(e.StartPrecision) != "year" {
		y += float64(e.StartDate.Month()-1) / 12
	}
	return y
}

func short(format string) string {
	if s, ok := shortFormats[format]; ok {
		return s
	}
	return format
}

// sortedByDate is rows oldest first, ties by title.
func sortedByDate(rows []sqlcgen.GroupEntriesRow) []sqlcgen.GroupEntriesRow {
	out := slices.Clone(rows)
	slices.SortStableFunc(out, func(a, b sqlcgen.GroupEntriesRow) int {
		return cmp.Or(cmp.Compare(year(a), year(b)), cmp.Compare(a.Title, b.Title))
	})
	return out
}
