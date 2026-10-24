package handler

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

func TestOrderingsKeepTheirShape(t *testing.T) {
	first, second := uuid.New(), uuid.New()
	arc := "the movie arc"
	row := func(o uuid.UUID, name string, pos, branch, step int32, id int64, title string, section *string) sqlcgen.GroupOrderingsRow {
		return sqlcgen.GroupOrderingsRow{OrderingID: o, Kind: "recommended", Name: name, Position: pos, Branch: branch,
			Step: step, Section: section, PublicID: id, Slug: "s", Title: title, Necessity: "required"}
	}

	got := orderings([]sqlcgen.GroupOrderingsRow{
		row(first, "start here", 1, 1, 1, 1, "Series", nil),
		// position 2 forks: one film, or two recap films in order
		row(first, "start here", 2, 1, 1, 2, "Film", &arc),
		row(first, "start here", 2, 2, 1, 3, "Recap I", &arc),
		row(first, "start here", 2, 2, 2, 4, "Recap II", &arc),
		// same section continues, so no second heading
		row(first, "start here", 3, 1, 1, 5, "Sequel", &arc),
		// a new ordering starts its sections over
		row(second, "release order", 1, 1, 1, 2, "Film", &arc),
	}, "/groups/9/g")

	step := func(id, title string) views.Step {
		return views.Step{Title: title, Href: "/entries/" + id + "/s"}
	}
	want := []views.Ordering{
		{ID: first, Kind: "recommended", Name: "start here", Href: "/groups/9/g/orderings/" + first.String() + "/start-here", Positions: []views.Position{
			{Ordinal: 1, Number: "01", Choices: []views.Choice{{Steps: []views.Step{step("1", "Series")}}}},
			{Ordinal: 2, Number: "02", Section: "the movie arc", Choices: []views.Choice{
				{Steps: []views.Step{step("2", "Film")}},
				{Steps: []views.Step{step("3", "Recap I"), step("4", "Recap II")}},
			}},
			{Ordinal: 3, Number: "03", Choices: []views.Choice{{Steps: []views.Step{step("5", "Sequel")}}}},
		}},
		{ID: second, Kind: "recommended", Name: "release order", Href: "/groups/9/g/orderings/" + second.String() + "/release-order", Positions: []views.Position{
			{Ordinal: 1, Number: "01", Section: "the movie arc", Choices: []views.Choice{{Steps: []views.Step{step("2", "Film")}}}},
		}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestSlugify(t *testing.T) {
	cases := []struct{ name, want string }{
		{"first-timer order", "first-timer-order"},
		{"chronological (shape test)", "chronological-shape-test"},
		{"  Anime-only, no VN!  ", "anime-only-no-vn"},
		{"23β cut", "23-cut"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := slugify(tc.name); got != tc.want {
				t.Errorf("slugify = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAnchor(t *testing.T) {
	ep := func(s string) *string { return &s }
	cases := []struct {
		name         string
		mine, theirs *string
		want         string
	}{
		{"none", nil, nil, ""},
		{"on the other work", nil, ep("Ep. 23"), "at Ep. 23"},
		{"on this one", ep("Ep. 23"), nil, "here at Ep. 23"},
		{"both", ep("Ep. 3"), ep("Ch. 12"), "here at Ep. 3, there at Ch. 12"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := anchor(tc.mine, tc.theirs); got != tc.want {
				t.Errorf("anchor = %q, want %q", got, tc.want)
			}
		})
	}
}
