package handler

import (
	"testing"

	"github.com/khinshankhan/listo/data/sqlcgen"
)

func TestEpisodes(t *testing.T) {
	n := func(v int32) *int32 { return &v }
	whole := func(total *int32) sqlcgen.GroupOrderingsRow {
		return sqlcgen.GroupOrderingsRow{Name: "o", Branch: 1, UnitKind: "episode", UnitTotal: total}
	}
	ranged := func(from, to *int32) sqlcgen.GroupOrderingsRow {
		r := whole(n(24))
		r.FromOrdinal, r.ToOrdinal = from, to
		return r
	}
	film := sqlcgen.GroupOrderingsRow{Name: "o", Branch: 1, UnitKind: "none"}
	fork := whole(n(12))
	fork.Branch = 2
	filmFork := film
	filmFork.Branch = 2
	other := whole(n(50))
	other.Name = "another"

	cases := []struct {
		name  string
		rows  []sqlcgen.GroupOrderingsRow
		want  int
		rough bool
	}{
		{"whole entries add up", []sqlcgen.GroupOrderingsRow{whole(n(24)), whole(n(13))}, 37, false},
		{"a range counts both ends", []sqlcgen.GroupOrderingsRow{ranged(n(1), n(22))}, 22, false},
		{"one unit", []sqlcgen.GroupOrderingsRow{ranged(n(23), n(23))}, 1, false},
		{"from a unit to the end", []sqlcgen.GroupOrderingsRow{ranged(n(9), nil)}, 16, false},
		{"films don't count", []sqlcgen.GroupOrderingsRow{whole(n(24)), film}, 24, false},
		{"a fork takes the first way", []sqlcgen.GroupOrderingsRow{whole(n(24)), fork}, 24, true},
		{"a fork of films is still a fork", []sqlcgen.GroupOrderingsRow{whole(n(12)), filmFork}, 12, true},
		{"an unknown total is a guess", []sqlcgen.GroupOrderingsRow{whole(n(24)), whole(nil)}, 24, true},
		{"other orderings don't count", []sqlcgen.GroupOrderingsRow{whole(n(24)), other}, 24, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, rough := episodes(tc.rows, "o")
			if got != tc.want || rough != tc.rough {
				t.Errorf("episodes = %d, %v; want %d, %v", got, rough, tc.want, tc.rough)
			}
		})
	}
}
