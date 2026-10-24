package handler

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

// branches are the test groups' public ids, by name.
var branches = map[string]int64{"Chaos;Head": 900, "Steins;Gate": 901, "Classic": 902, "2016 reboot": 903, "Index": 904, "Railgun": 905}

// tlEntry is an entry released in year and month (0 for just a year), under
// branch, or straight in the group when branch is "".
func tlEntry(id int64, title, format string, year, month int, branch string) sqlcgen.GroupEntriesRow {
	row := sqlcgen.GroupEntriesRow{PublicID: id, Title: title, FormatLabel: format, BranchPublicID: id, BranchType: "entry"}
	if year != 0 {
		precision := "year"
		if month != 0 {
			precision = "month"
		} else {
			month = 1
		}
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		row.StartDate, row.StartPrecision = &start, &precision
	}
	if branch != "" {
		kind := "franchise"
		row.BranchPublicID, row.BranchType, row.BranchName, row.BranchKind = branches[branch], "group", &branch, &kind
	}
	return row
}

func tlRelation(from int64, kind string, to int64) sqlcgen.GroupRelationsRow {
	labels := map[string][2]string{
		"adaptation_of":  {"adaptation of", "derivative"},
		"sequel_of":      {"sequel of", "sequence"},
		"crossover_with": {"crossover with", "setting"},
		"reboot_of":      {"reboot of", "derivative"},
		"references":     {"references", "setting"},
		"side_story_of":  {"side story of", "sequence"},
		"spin_off_of":    {"spin-off of", "sequence"},
	}
	return sqlcgen.GroupRelationsRow{FromID: from, Kind: kind, Label: labels[kind][0], Family: labels[kind][1], ToID: to}
}

func workNamed(t *testing.T, tl *views.Timeline, name string) views.TimelineWork {
	t.Helper()
	for _, w := range tl.Works {
		if w.Name == name {
			return w
		}
	}
	t.Fatalf("no work %q among %v", name, tl.Works)
	return views.TimelineWork{}
}

func dotLabels(w views.TimelineWork) []string {
	var out []string
	for _, d := range w.Dots {
		out = append(out, d.Label)
	}
	return out
}

func TestTimelineNeedsTwoWorks(t *testing.T) {
	// an anime of a VN is the same work, so this is one
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(1, "Steins;Gate", "visual novel", 2009, 0, ""),
		tlEntry(2, "Steins;Gate", "anime", 2011, 4, ""),
	}
	if tl := timeline(rows, []sqlcgen.GroupRelationsRow{tlRelation(2, "adaptation_of", 1)}); tl != nil {
		t.Errorf("timeline = %+v, want none for one work", tl)
	}
}

func TestTimelineWorks(t *testing.T) {
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(1, "Chaos;Head", "visual novel", 2008, 4, "Chaos;Head"),
		tlEntry(2, "Chaos;Head", "anime", 2008, 10, "Chaos;Head"),
		tlEntry(3, "Chaos;Child", "visual novel", 2014, 12, "Chaos;Head"),
		tlEntry(4, "Chaos;Child", "anime", 2017, 1, "Chaos;Head"),
		tlEntry(5, "Steins;Gate", "visual novel", 2009, 0, "Steins;Gate"),
		tlEntry(6, "Steins;Gate", "anime", 2011, 4, "Steins;Gate"),
		tlEntry(7, "Robotics;Notes DaSH", "visual novel", 2019, 0, ""),
		tlEntry(8, "Lost Episode", "visual novel", 0, 0, ""),
	}
	rels := []sqlcgen.GroupRelationsRow{
		tlRelation(2, "adaptation_of", 1),
		tlRelation(3, "sequel_of", 1),
		tlRelation(4, "adaptation_of", 3),
		tlRelation(6, "adaptation_of", 5),
		tlRelation(7, "crossover_with", 6),
		// made by Chaos;Child, long after either end
		tlRelation(5, "references", 1),
	}
	rels[4].Spoiler = true
	crossover := time.Date(2019, 1, 31, 0, 0, 0, 0, time.UTC)
	reference := time.Date(2014, 12, 18, 0, 0, 0, 0, time.UTC)
	rels[4].EstablishedOn, rels[5].EstablishedOn = &crossover, &reference

	tl := timeline(rows, rels)
	if tl == nil {
		t.Fatal("no timeline")
	}

	t.Run("adaptations join their source, a year to a dot", func(t *testing.T) {
		tests := []struct {
			work string
			dots []string
		}{
			{"Chaos;Head", []string{"VN, anime"}},
			{"Chaos;Child", []string{"VN", "anime"}},
			{"Steins;Gate", []string{"VN", "anime"}},
			{"Robotics;Notes DaSH", []string{"VN"}},
		}
		for _, tt := range tests {
			if got := dotLabels(workNamed(t, tl, tt.work)); !slices.Equal(got, tt.dots) {
				t.Errorf("%s dots = %q, want %q", tt.work, got, tt.dots)
			}
		}
		if len(tl.Works) != 4 {
			t.Errorf("%d works, want 4", len(tl.Works))
		}
	})

	t.Run("a sequel with room follows on the same row", func(t *testing.T) {
		if head, child := workNamed(t, tl, "Chaos;Head"), workNamed(t, tl, "Chaos;Child"); head.Y != child.Y {
			t.Errorf("Chaos;Child on y %d, Chaos;Head on %d", child.Y, head.Y)
		}
		if head, sg := workNamed(t, tl, "Chaos;Head"), workNamed(t, tl, "Steins;Gate"); head.Y == sg.Y {
			t.Error("unrelated works share a row")
		}
	})

	t.Run("each group gets a box, loose entries none", func(t *testing.T) {
		var names []string
		for _, b := range tl.Boxes {
			names = append(names, b.Name)
		}
		if want := []string{"Chaos;Head", "Steins;Gate"}; !slices.Equal(names, want) {
			t.Errorf("boxes = %q, want %q", names, want)
		}
		dash := workNamed(t, tl, "Robotics;Notes DaSH")
		for _, b := range tl.Boxes {
			if dash.Y > b.Y && dash.Y < b.Y+b.H {
				t.Errorf("DaSH (y %d) drawn inside %s", dash.Y, b.Name)
			}
		}
	})

	// the island lights a work's line by these
	t.Run("links run left to right", func(t *testing.T) {
		type ends struct{ source, target int64 }
		want := []ends{
			{1, 3},
			// drawn from Steins;Gate, though DaSH is the one crossing over
			{5, 7},
			// back in time, so flipped
			{1, 5},
		}
		var got []ends
		for _, l := range tl.Links {
			got = append(got, ends{l.Source, l.Target})
		}
		if !slices.Equal(got, want) {
			t.Errorf("link ends = %v, want %v", got, want)
		}
	})

	// DaSH, in no box, sits above them all
	t.Run("a crossover leaves the release it crosses over with", func(t *testing.T) {
		sg := workNamed(t, tl, "Steins;Gate")
		if start, _ := pathEnds(t, tl.Links[1].Path); start != fmt.Sprintf("%d %d", sg.Dots[1].X, sg.Card.Y) {
			t.Errorf("crossover starts at %s, want the anime's dot (x %d) on the card's top (%d)", start, sg.Dots[1].X, sg.Card.Y)
		}
	})

	t.Run("links", func(t *testing.T) {
		var got []views.TimelineLink
		for _, l := range tl.Links {
			l.Path, l.LabelX, l.LabelY, l.LabelW, l.LabelEnd = "", 0, 0, 0, false
			l.Source, l.Target = 0, 0
			got = append(got, l)
		}
		want := []views.TimelineLink{
			{Sequel: true},
			{Label: "crossover, 2019", Spoiler: true},
			{Label: "reference, 2014"},
		}
		if !slices.Equal(got, want) {
			t.Errorf("links = %+v, want %+v", got, want)
		}
		if !tl.Spoilers {
			t.Error("Spoilers = false with a spoiler link")
		}
	})

	t.Run("in words", func(t *testing.T) {
		want := []views.TimelineLine{
			{Text: "Chaos;Head: VN 2008, anime 2008"},
			{Text: "Steins;Gate: VN 2009, anime 2011"},
			{Text: "Chaos;Child: VN 2014, anime 2017"},
			{Text: "Robotics;Notes DaSH: VN 2019"},
			{Text: "Chaos;Child, sequel of Chaos;Head"},
			{Text: "Robotics;Notes DaSH, crossover with Steins;Gate, 2019", Spoiler: true},
			{Text: "Steins;Gate, references Chaos;Head, 2014"},
		}
		if !slices.Equal(tl.Lines, want) {
			t.Errorf("lines = %+v\nwant %+v", tl.Lines, want)
		}
	})

	if tl.Undated != 1 {
		t.Errorf("Undated = %d, want 1", tl.Undated)
	}
}

// A long name leaves no room for its sequel, which drops to a row of its own;
// a reboot's link is dated by the reboot.
func TestTimelineCrowdedSequels(t *testing.T) {
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(68, "Ben 10", "TV", 2005, 0, "Classic"),
		tlEntry(70, "Ben 10: Alien Force", "TV", 2008, 0, "Classic"),
		tlEntry(71, "Ben 10: Ultimate Alien", "TV", 2010, 0, "Classic"),
		tlEntry(73, "Ben 10", "TV", 2016, 0, "2016 reboot"),
	}
	rels := []sqlcgen.GroupRelationsRow{
		tlRelation(70, "sequel_of", 68),
		tlRelation(71, "sequel_of", 70),
		tlRelation(73, "reboot_of", 68),
	}
	tl := timeline(rows, rels)

	ben, alien, ultimate := workNamed(t, tl, "Ben 10"), workNamed(t, tl, "Ben 10: Alien Force"), workNamed(t, tl, "Ben 10: Ultimate Alien")
	if alien.Y != ben.Y {
		t.Errorf("Alien Force on y %d, want Ben 10's %d", alien.Y, ben.Y)
	}
	if ultimate.Y <= alien.Y {
		t.Errorf("Ultimate Alien on y %d, want below Alien Force's %d", ultimate.Y, alien.Y)
	}
	labels := []string{}
	for _, l := range tl.Links {
		labels = append(labels, l.Label)
	}
	if want := []string{"", "", "reboot, 2016"}; !slices.Equal(labels, want) {
		t.Errorf("link labels = %q, want %q", labels, want)
	}
	// two works called Ben 10 need their years to tell them apart in words
	if !slices.Contains(tl.Lines, views.TimelineLine{Text: "Ben 10 (2016), reboot of Ben 10 (2005), 2016"}) {
		t.Errorf("lines = %+v, want the reboot told apart by year", tl.Lines)
	}
	if !slices.Contains(tl.Lines, views.TimelineLine{Text: "Ben 10 (2005): TV 2005"}) {
		t.Errorf("lines = %+v, want the first Ben 10 told apart by year", tl.Lines)
	}
}

// Rows keep the shape of the story: a side story sits right under the work it
// hangs off, ahead of anything later, and side stories with room share a row.
func TestTimelineKeepsStructure(t *testing.T) {
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(1, "Index", "anime", 2008, 10, "Index"),
		tlEntry(2, "Unrelated", "anime", 2009, 1, "Index"),
		tlEntry(3, "Side A", "OVA", 2010, 1, "Index"),
		tlEntry(4, "Index II", "anime", 2010, 10, "Index"),
		tlEntry(5, "Side B", "OVA", 2014, 1, "Index"),
	}
	rels := []sqlcgen.GroupRelationsRow{
		tlRelation(3, "side_story_of", 1),
		tlRelation(4, "sequel_of", 1),
		tlRelation(5, "side_story_of", 1),
	}
	tl := timeline(rows, rels)
	index, ii, a, b, other := workNamed(t, tl, "Index"), workNamed(t, tl, "Index II"), workNamed(t, tl, "Side A"),
		workNamed(t, tl, "Side B"), workNamed(t, tl, "Unrelated")
	if ii.Y != index.Y {
		t.Errorf("Index II on y %d, want Index's %d", ii.Y, index.Y)
	}
	if !(index.Y < a.Y && a.Y < other.Y) {
		t.Errorf("side story on y %d, want between Index (%d) and the later root (%d)", a.Y, index.Y, other.Y)
	}
	if b.Y != a.Y {
		t.Errorf("Side B on y %d, want Side A's %d", b.Y, a.Y)
	}
}

// Where a link leaves and meets its cards, and where its label goes.
func TestTimelineRoutes(t *testing.T) {
	index := tlEntry(1, "A Certain Magical Index", "anime", 2008, 10, "Index")
	route := func(t *testing.T, child sqlcgen.GroupEntriesRow, kind string) (views.TimelineWork, views.TimelineWork, views.TimelineLink) {
		t.Helper()
		tl := timeline([]sqlcgen.GroupEntriesRow{index, child}, []sqlcgen.GroupRelationsRow{tlRelation(2, kind, 1)})
		if len(tl.Links) != 1 {
			t.Fatalf("links = %+v, want one", tl.Links)
		}
		return workNamed(t, tl, index.Title), workNamed(t, tl, child.Title), tl.Links[0]
	}
	at := func(x, y int) string { return fmt.Sprintf("%d %d", x, y) }

	t.Run("a sequel with room carries on along the row, card to card", func(t *testing.T) {
		a, b, l := route(t, tlEntry(2, "Index II", "anime", 2012, 10, "Index"), "sequel_of")
		if want := fmt.Sprintf("M%d %dL%d %d", a.Card.X+a.Card.W, a.Y, b.Card.X, b.Y); l.Path != want {
			t.Errorf("path = %s, want %s", l.Path, want)
		}
	})

	t.Run("a side story hangs off the trunk and runs in", func(t *testing.T) {
		a, b, l := route(t, tlEntry(2, "Railgun", "anime", 2010, 10, "Index"), "side_story_of")
		start, end := pathEnds(t, l.Path)
		if start != at(a.Card.X+8, a.Card.Y+a.Card.H) || end != at(b.Card.X, b.Y) {
			t.Errorf("runs %s to %s, want %s to %s", start, end, at(a.Card.X+8, a.Card.Y+a.Card.H), at(b.Card.X, b.Y))
		}
		if !l.LabelEnd || l.LabelX != b.Card.X-14 || l.LabelY != b.Y {
			t.Errorf("label ends at %d, %d (end %t), want 14 short of the card on its row: %d, %d", l.LabelX, l.LabelY, l.LabelEnd, b.Card.X-14, b.Y)
		}
	})

	t.Run("one starting too near the trunk meets the card from above, labelled beside the trunk", func(t *testing.T) {
		a, b, l := route(t, tlEntry(2, "OVA", "OVA", 2008, 11, "Index"), "side_story_of")
		if _, end := pathEnds(t, l.Path); end != at(b.Dots[0].X, b.Card.Y) {
			t.Errorf("ends at %s, want its dot on the card's top: %s", end, at(b.Dots[0].X, b.Card.Y))
		}
		if !l.LabelEnd || l.LabelX != a.Card.X+8-6 || l.LabelY != b.Card.Y-10 {
			t.Errorf("label ends at %d, %d, want left of the trunk on the run across: %d, %d", l.LabelX, l.LabelY, a.Card.X+2, b.Card.Y-10)
		}
	})

	t.Run("anything else drops from its dot onto the card", func(t *testing.T) {
		a, b, l := route(t, tlEntry(2, "Railgun", "anime", 2010, 10, "Railgun"), "references")
		start, end := pathEnds(t, l.Path)
		if start != at(a.Dots[0].X, a.Card.Y+a.Card.H) || end != at(b.Dots[0].X, b.Card.Y) {
			t.Errorf("runs %s to %s, want dot to dot: %s to %s", start, end, at(a.Dots[0].X, a.Card.Y+a.Card.H), at(b.Dots[0].X, b.Card.Y))
		}
		if l.LabelEnd || l.LabelY != b.Card.Y-10 {
			t.Errorf("label at y %d (end %t), want centred on the run across, %d", l.LabelY, l.LabelEnd, b.Card.Y-10)
		}
	})
}

// A sequel that came out first, as Bakemonogatari did before the
// Kizumonogatari films it follows, goes round the left of everything between
// them rather than up through it.
func TestTimelineBackwardSequel(t *testing.T) {
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(1, "Bakemonogatari", "anime", 2009, 7, "Index"),
		tlEntry(2, "Nisemonogatari", "anime", 2012, 1, "Index"),
		tlEntry(3, "Hanamonogatari", "anime", 2014, 8, "Index"),
		tlEntry(4, "Kizumonogatari III", "film", 2017, 1, "Index"),
	}
	rels := []sqlcgen.GroupRelationsRow{tlRelation(2, "sequel_of", 1), tlRelation(1, "sequel_of", 4)}
	tl := timeline(rows, rels)
	bake, hana, kizu := workNamed(t, tl, "Bakemonogatari"), workNamed(t, tl, "Hanamonogatari"), workNamed(t, tl, "Kizumonogatari III")
	var back views.TimelineLink
	for _, l := range tl.Links {
		if l.Source == 4 {
			back = l
		}
	}
	start, end := pathEnds(t, back.Path)
	if want := fmt.Sprintf("%d %d", kizu.Card.X, kizu.Y); start != want {
		t.Errorf("starts at %s, want Kizumonogatari's left edge on its row, %s", start, want)
	}
	if want := fmt.Sprintf("%d %d", bake.Card.X, bake.Y); end != want {
		t.Errorf("ends at %s, want Bakemonogatari's left edge on its row, %s", end, want)
	}
	// so round Hanamonogatari, between them, too
	if left := slices.Min(pathXs(t, back.Path)); left >= float64(min(bake.Card.X, hana.Card.X, kizu.Card.X)) {
		t.Errorf("goes no further left than x %.0f, want left of the cards (from x %d)", left, min(bake.Card.X, hana.Card.X, kizu.Card.X))
	}
}

// pathXs is every x a path visits.
func pathXs(t *testing.T, path string) []float64 {
	t.Helper()
	var xs []float64
	for _, seg := range strings.FieldsFunc(path, func(r rune) bool { return r == 'M' || r == 'L' || r == 'Q' }) {
		nums := strings.Fields(seg)
		for i := 0; i+1 < len(nums); i += 2 {
			x, err := strconv.ParseFloat(nums[i], 64)
			if err != nil {
				t.Fatalf("path %q: %v", path, err)
			}
			xs = append(xs, x)
		}
	}
	return xs
}

// pathEnds is where a path starts and ends, as "x y".
func pathEnds(t *testing.T, path string) (string, string) {
	t.Helper()
	first := strings.IndexAny(path[1:], "LQ") + 1
	last := strings.LastIndex(path, "L")
	if !strings.HasPrefix(path, "M") || first < 1 || last < 0 {
		t.Fatalf("path %q isn't M..L", path)
	}
	return path[1:first], path[last+1:]
}
