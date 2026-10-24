package handler

import (
	"slices"
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
		tlRelation(7, "crossover_with", 5),
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

	t.Run("links", func(t *testing.T) {
		var got []views.TimelineLink
		for _, l := range tl.Links {
			l.Path, l.Head, l.LabelX, l.LabelY = "", "", 0, 0
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

// A link between rows is labelled between them, not over a name above.
func TestTimelineLabelsBetweenRows(t *testing.T) {
	rows := []sqlcgen.GroupEntriesRow{
		tlEntry(1, "A Certain Magical Index", "anime", 2008, 10, "Index"),
		tlEntry(2, "A Certain Scientific Railgun", "anime", 2009, 10, "Railgun"),
	}
	tl := timeline(rows, []sqlcgen.GroupRelationsRow{tlRelation(2, "spin_off_of", 1)})
	index, railgun := workNamed(t, tl, "A Certain Magical Index"), workNamed(t, tl, "A Certain Scientific Railgun")
	if len(tl.Links) != 1 {
		t.Fatalf("links = %+v, want the spin-off", tl.Links)
	}
	if l := tl.Links[0]; l.Label != "spin-off" || l.LabelY <= index.Y || l.LabelY >= railgun.NameY {
		t.Errorf("label %q at y %d, want between Index's row (%d) and Railgun's name (%d)", l.Label, l.LabelY, index.Y, railgun.NameY)
	}
}
