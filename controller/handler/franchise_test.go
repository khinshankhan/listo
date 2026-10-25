package handler

import (
	"reflect"
	"testing"
	"time"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

func TestFranchiseGrouping(t *testing.T) {
	ptr := func(s string) *string { return &s }
	date := func(y int, m time.Month) *time.Time { d := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC); return &d }
	month, eps := ptr("month"), int32(24)
	branch := func(r sqlcgen.GroupEntriesRow, id int64, slug, name, kind, standing string) sqlcgen.GroupEntriesRow {
		r.BranchPublicID, r.BranchSlug, r.BranchType, r.BranchName, r.BranchKind = id, slug, "group", ptr(name), ptr(kind)
		if standing != "" {
			r.BranchStanding = ptr(standing)
		}
		return r
	}
	vn := sqlcgen.GroupEntriesRow{PublicID: 6, Slug: "vn", Title: "Steins;Gate", Medium: "visual_novel", MediumLabel: "visual novel",
		Format: "full", FormatLabel: "visual novel", StartDate: date(2009, 10), StartPrecision: ptr("year"), UnitKind: "none", Role: ptr("source material")}
	anime := sqlcgen.GroupEntriesRow{PublicID: 7, Slug: "anime", Title: "Steins;Gate", Medium: "anime", MediumLabel: "anime",
		Format: "tv", FormatLabel: "anime", StartDate: date(2011, 4), StartPrecision: month, UnitKind: "episode", UnitTotal: &eps, Role: ptr("adaptation of the VN")}
	zero := sqlcgen.GroupEntriesRow{PublicID: 11, Slug: "zero", Title: "Steins;Gate 0", Medium: "anime", MediumLabel: "anime",
		Format: "tv", FormatLabel: "anime", StartDate: date(2018, 4), StartPrecision: month, UnitKind: "none"}
	film := sqlcgen.GroupEntriesRow{PublicID: 8, Slug: "film", Title: "The Movie", Medium: "anime", MediumLabel: "anime",
		Format: "movie", FormatLabel: "film", StartDate: date(2013, 4), StartPrecision: ptr("year"), UnitKind: "none"}
	loose := sqlcgen.GroupEntriesRow{PublicID: 20, Slug: "loose", Title: "Loose", Medium: "anime", MediumLabel: "anime",
		Format: "tv", FormatLabel: "anime", UnitKind: "none", BranchPublicID: 20, BranchSlug: "loose", BranchType: "entry"}

	// release order, as the query returns them; anime also sits in a series
	rows := []sqlcgen.GroupEntriesRow{
		branch(vn, 3, "main", "Main timeline", "continuity", "primary"),
		branch(anime, 3, "main", "Main timeline", "continuity", "primary"),
		branch(anime, 9, "extras", "Extras", "series", ""),
		branch(film, 3, "main", "Main timeline", "continuity", "primary"),
		branch(zero, 4, "zero", "Divergence (0)", "continuity", "concluded"),
		loose,
	}

	vnRow := views.EntryRow{Title: "Steins;Gate", Href: "/entries/6/vn", Role: "source material", Medium: "visual novel",
		FormatHref: "/entries?format=full&medium=visual_novel", Detail: "2009"}
	animeRow := views.EntryRow{Title: "Steins;Gate", Href: "/entries/7/anime", Role: "adaptation of the VN", Medium: "anime",
		FormatHref: "/entries?format=tv&medium=anime", Detail: "spring 2011 · 24 episodes"}
	filmRow := views.EntryRow{Title: "The Movie", Href: "/entries/8/film", Medium: "film", FormatHref: "/entries?format=movie&medium=anime", Detail: "2013"}
	zeroRow := views.EntryRow{Title: "Steins;Gate 0", Href: "/entries/11/zero", Medium: "anime", FormatHref: "/entries?format=tv&medium=anime", Detail: "spring 2018"}
	looseRow := views.EntryRow{Title: "Loose", Href: "/entries/20/loose", Medium: "anime", FormatHref: "/entries?format=tv&medium=anime"}

	cases := []struct {
		by   string
		want []views.EntrySection
	}{
		// continuities by standing, then other groups, then loose entries
		{"continuity", []views.EntrySection{
			{Label: "Main timeline", Href: "/groups/3/main", Standing: "primary", Rows: []views.EntryRow{vnRow, animeRow, filmRow}},
			{Label: "Divergence (0)", Href: "/groups/4/zero", Standing: "concluded", Rows: []views.EntryRow{zeroRow}},
			{Label: "Extras", Href: "/groups/9/extras", Rows: []views.EntryRow{animeRow}},
			{Label: "directly in Steins;Gate", Rows: []views.EntryRow{looseRow}},
		}},
		// each entry once per medium, mediums in order of first release
		{"medium", []views.EntrySection{
			{Label: "visual novel", Rows: []views.EntryRow{vnRow}},
			{Label: "anime", Rows: []views.EntryRow{animeRow, filmRow, zeroRow, looseRow}},
		}},
		{"release", []views.EntrySection{
			{Rows: []views.EntryRow{vnRow, animeRow, filmRow, zeroRow, looseRow}},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.by, func(t *testing.T) {
			if got := entrySections(rows, tc.by, "Steins;Gate"); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got\n%+v\nwant\n%+v", got, tc.want)
			}
		})
	}

	t.Run("continuity cards", func(t *testing.T) {
		want := []views.ContinuityCard{
			{Name: "Main timeline", Href: "/groups/3/main", Standing: "primary", Entries: 3},
			{Name: "Divergence (0)", Href: "/groups/4/zero", Standing: "concluded", Entries: 1},
		}
		if got := continuityCards(rows); !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
