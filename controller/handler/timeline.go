package handler

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// the timeline's drawing, in SVG units
const (
	tlWidth = 960
	// the first year starts at tlLeft and the last ends at tlRight
	tlLeft, tlRight = 56, 904
	// names keep this far from the edge
	tlEdge = 24
	tlRow  = 64
	// between two works on a row, and between boxes
	tlGap = 24
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
	// a work sits where its first release does
	names := map[string]int{}
	for _, w := range works {
		first := w.entries[0]
		if first.BranchType == "group" {
			w.branch, w.branchName = first.BranchPublicID, deref(first.BranchName)
		}
		w.view.Name = first.Title
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
		// work's other side stories while there's room
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
			default:
				l = &lane{}
				lanes = append(lanes, l)
			}
			l.right = max(l.right, w.right)
			laneOf[w] = l
		}

		top := y
		firstRow := y + 34
		if branch != 0 {
			firstRow = y + 62
		}
		for w, l := range laneOf {
			w.view.Y = firstRow + slices.Index(lanes, l)*tlRow
			w.view.NameY = w.view.Y - 14
		}
		bottom := firstRow + (len(lanes)-1)*tlRow + 34
		if branch != 0 {
			tl.Boxes = append(tl.Boxes, views.TimelineBox{X: 16, Y: top, W: tlWidth - 32, H: bottom - top, Name: name})
		}
		y = bottom + 16
	}
	tl.Height = y + 24

	dotAt := func(id int64) (int, int) {
		w := workOf[id]
		for i, e := range w.entries {
			if e.PublicID == id {
				return w.view.Dots[dotIndex(w, i)].X, w.view.Y
			}
		}
		return 0, 0
	}
	// labels go where they cover no text
	var taken []box
	for _, w := range works {
		name := nameFont.width(w.view.Name)
		left := float64(w.view.NameX)
		if w.view.EndAnchored {
			left -= name
		}
		taken = append(taken, box{left, float64(w.view.NameY - 11), left + name, float64(w.view.NameY + 3)})
		for _, d := range w.view.Dots {
			half := dotFont.width(d.Label) / 2
			taken = append(taken, box{float64(d.X) - half, float64(w.view.Y + 8), float64(d.X) + half, float64(w.view.Y + 23)},
				box{float64(d.X - 7), float64(w.view.Y - 7), float64(d.X + 7), float64(w.view.Y + 7)})
		}
	}
	for _, b := range tl.Boxes {
		taken = append(taken, box{float64(b.X + 18), float64(b.Y + 12), float64(b.X+18) + float64(utf8.RuneCountInString(b.Name))*7.8, float64(b.Y + 30)})
	}
	// a label sits on its line, or just beside it; with nowhere clear to go
	// it's left off, and the words still say it
	label := func(link *views.TimelineLink, c bezier) {
		half := float64(len(link.Label)*6+8) / 2
		for _, beside := range []float64{0, -half - 8, half + 8} {
			for _, t := range []float64{0.5, 0.4, 0.6, 0.3, 0.7, 0.2, 0.8, 0.1, 0.9} {
				px, py := c.at(t)
				px += beside
				at := box{px - half, py - 8, px + half, py + 8}
				if at.x1 > 0 && at.x2 < tlWidth && !slices.ContainsFunc(taken, at.overlaps) {
					taken = append(taken, at)
					link.LabelX, link.LabelY = int(math.Round(px)), int(math.Round(py+4))
					return
				}
			}
		}
		link.Label = ""
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

		link := views.TimelineLink{Sequel: rel.Family == "sequence", Spoiler: rel.Spoiler}
		line := fmt.Sprintf("%s, %s %s", called(byID[rel.FromID]), rel.Label, called(byID[rel.ToID]))
		if link.Sequel {
			// from the release it continues to the start of the sequel
			x1, y1 := dotAt(rel.ToID)
			x2, y2 := from.view.Dots[0].X-9, from.view.Y
			c := connect(x1, y1, x2, y2)
			link.Path = c.path()
			link.Head = fmt.Sprintf("M%d %d L%d %d L%d %dz", x2-1, y2-3, x2+6, y2, x2-1, y2+3)
			if link.Label = name; name != "" {
				label(&link, c)
			}
		} else {
			when := max(year(byID[rel.FromID]), year(byID[rel.ToID]))
			if rel.EstablishedOn != nil {
				when = float64(rel.EstablishedOn.Year())
			}
			link.Label = fmt.Sprintf("%s, %d", name, int(when))
			line += fmt.Sprintf(", %d", int(when))
			// drawn left to right, whichever way the relation reads
			x1, y1 := dotAt(rel.FromID)
			x2, y2 := dotAt(rel.ToID)
			if x2 < x1 {
				x1, y1, x2, y2 = x2, y2, x1, y1
			}
			c := connect(x1, y1, x2-7, y2)
			if y1 == y2 {
				// over the row, clear of its bars
				c = bezier{float64(x1), float64(y1 - 7), float64(x1), float64(y1 - 34), float64(x2), float64(y1 - 34), float64(x2), float64(y2 - 7)}
			}
			link.Path = c.path()
			label(&link, c)
		}
		tl.Spoilers = tl.Spoilers || link.Spoiler
		tl.Links = append(tl.Links, link)
		tl.Lines = append(tl.Lines, views.TimelineLine{Text: line, Spoiler: rel.Spoiler})
	}

	var lines []views.TimelineLine
	for _, w := range works {
		var releases []string
		for _, e := range w.entries {
			releases = append(releases, fmt.Sprintf("%s %d", short(e.FormatLabel), e.StartDate.Year()))
		}
		lines = append(lines, views.TimelineLine{Text: called(w.entries[0]) + ": " + strings.Join(releases, ", ")})
		tl.Works = append(tl.Works, w.view)
	}
	tl.Lines = append(lines, tl.Lines...)
	return tl
}

// lane is a row being filled left to right, and the row it hangs under.
type lane struct {
	right float64
	under *lane
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

// box is a rectangle of the drawing, corner to corner.
type box struct{ x1, y1, x2, y2 float64 }

func (b box) overlaps(o box) bool {
	return b.x1 < o.x2 && o.x1 < b.x2 && b.y1 < o.y2 && o.y1 < b.y2
}

// bezier is a cubic curve: its start, two control points and its end.
type bezier struct{ x0, y0, x1, y1, x2, y2, x3, y3 float64 }

// connect runs from the dot at x1, y1 to x2 on the row at y2, arriving level
// so it misses the names above the dots and the labels below. Along a row
// it's straight; a long way across it's an S, leaving the dot level too; a
// short way, to a row below, it drops from the dot and turns in, rather than
// loop out and back.
func connect(x1, y1, x2, y2 int) bezier {
	switch {
	case y1 == y2:
		return straight(x1+7, y1, x2, y2)
	case x2-x1 >= 80:
		pull := float64((x2 - x1 - 7) / 2)
		return bezier{float64(x1 + 7), float64(y1), float64(x1+7) + pull, float64(y1), float64(x2) - pull, float64(y2), float64(x2), float64(y2)}
	default:
		down := float64(y1 + 7*cmp.Compare(y2, y1))
		return bezier{float64(x1), down, float64(x1), float64(y2), float64(x1), float64(y2), float64(x2), float64(y2)}
	}
}

// straight is a straight run as a curve, so labels can sit along it the same way.
func straight(x1, y1, x2, y2 int) bezier {
	dx, dy := float64(x2-x1)/3, float64(y2-y1)/3
	return bezier{float64(x1), float64(y1), float64(x1) + dx, float64(y1) + dy, float64(x2) - dx, float64(y2) - dy, float64(x2), float64(y2)}
}

// at is the point t of the way along the curve.
func (c bezier) at(t float64) (float64, float64) {
	u := 1 - t
	a, b, d, e := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
	return a*c.x0 + b*c.x1 + d*c.x2 + e*c.x3, a*c.y0 + b*c.y1 + d*c.y2 + e*c.y3
}

func (c bezier) path() string {
	r := func(f float64) int { return int(math.Round(f)) }
	return fmt.Sprintf("M%d %d C%d %d %d %d %d %d", r(c.x0), r(c.y0), r(c.x1), r(c.y1), r(c.x2), r(c.y2), r(c.x3), r(c.y3))
}

// place sets a work's dots and name, and the stretch of its row they take.
func place(w *tlWork, x func(float64) float64) {
	for i, e := range w.entries {
		label := short(e.FormatLabel)
		if i > 0 && e.StartDate.Year() == w.entries[i-1].StartDate.Year() {
			d := &w.view.Dots[len(w.view.Dots)-1]
			if !strings.Contains(d.Label, label) {
				d.Label += ", " + label
			}
			continue
		}
		w.view.Dots = append(w.view.Dots, views.TimelineDot{X: int(math.Round(x(year(e)))), Label: label})
	}

	firstX, lastX := float64(w.view.Dots[0].X), float64(w.view.Dots[len(w.view.Dots)-1].X)
	name := nameFont.width(w.view.Name)
	w.view.NameX = int(firstX) - 5
	w.left, w.right = firstX-5, firstX-5+name
	// too long to fit, it ends at its dot instead
	if w.right > tlWidth-tlEdge {
		w.view.EndAnchored, w.view.NameX = true, int(firstX)+6
		w.left, w.right = firstX+6-name, firstX+6
	}
	for _, d := range w.view.Dots {
		half := dotFont.width(d.Label) / 2
		w.left, w.right = min(w.left, float64(d.X)-half), max(w.right, float64(d.X)+half, lastX)
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
