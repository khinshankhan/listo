package handler

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/views"
)

func TestRelative(t *testing.T) {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		at   time.Time
		want string
	}{
		{"days and hours ahead", now.Add(3*24*time.Hour + 4*time.Hour + 30*time.Minute), "in 3d 4h"},
		{"hours and minutes ahead", now.Add(5*time.Hour + 7*time.Minute), "in 5h 7m"},
		{"minutes ahead", now.Add(40 * time.Minute), "in 40m"},
		{"seconds ahead", now.Add(20 * time.Second), "any minute"},
		{"just gone", now.Add(-10 * time.Minute), "just now"},
		{"an hour ago", now.Add(-time.Hour), "1 hour ago"},
		{"hours ago", now.Add(-5 * time.Hour), "5 hours ago"},
		{"yesterday", now.Add(-30 * time.Hour), "yesterday"},
		{"days ago", now.Add(-2*24*time.Hour - time.Hour), "2 days ago"},
		{"a week ago", now.Add(-8 * 24 * time.Hour), "1 week ago"},
		{"weeks ago", now.Add(-20 * 24 * time.Hour), "2 weeks ago"},
		{"a date, past a month", time.Date(2011, 4, 5, 15, 0, 0, 0, time.UTC), "Apr 5, 2011"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := relative(tc.at, now); got != tc.want {
				t.Errorf("relative = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDurationAndUnits(t *testing.T) {
	n := func(v int32) *int32 { return &v }
	available := func(v int64) pgtype.Numeric { return pgtype.Numeric{Int: big.NewInt(v), Valid: true} }
	cases := []struct {
		name, got, want string
	}{
		{"an episode", duration(n(24), nil), "24 min"},
		{"an episode over a whole runtime", duration(n(24), n(600)), "24 min"},
		{"a film", duration(nil, n(90)), "1 hr 30 min"},
		{"on the hour", duration(nil, n(120)), "2 hr"},
		{"nothing known", duration(nil, nil), ""},
		{"finished", units(n(24), pgtype.Numeric{}, "finished"), "24"},
		{"airing, out of a known total", units(n(28), available(19), "releasing"), "19 of 28"},
		{"airing, total unknown", units(nil, available(19), "releasing"), "19 so far"},
		{"a stale count once it's finished", units(n(28), available(19), "finished"), "28"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("got %q, want %q", tc.got, tc.want)
			}
		})
	}
}

func TestClip(t *testing.T) {
	cases := []struct {
		name, text string
		n          int
		want       string
	}{
		{"short enough", "A short one.", 20, "A short one."},
		{"at a word", "One two three four five", 14, "One two three…"},
		{"no trailing comma", "One two, three four", 10, "One two…"},
		{"runes, not bytes", "ふりーれん ふりーれん", 8, "ふりーれん…"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := clip(tc.text, tc.n)
			if got != tc.want || len([]rune(got)) > tc.n {
				t.Errorf("clip = %q, want %q within %d runes", got, tc.want, tc.n)
			}
		})
	}
}

func TestEpisodeList(t *testing.T) {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	at := func(d time.Duration) *time.Time { v := now.Add(d); return &v }
	title := func(s string) *string { return &s }
	day := 24 * time.Hour
	airing := []sqlcgen.EntryUnitsRow{
		{Label: "Ep. 1", Title: title("One"), AirDate: at(-15 * day)},
		{Label: "Ep. 2", Title: title("Two"), AirDate: at(-8 * day)},
		{Label: "Ep. 3", AirDate: at(2 * day)},
		{Label: "Ep. 4", AirDate: at(9 * day)},
	}
	labels := func(eps []views.Episode) string {
		var out []string
		for _, e := range eps {
			out = append(out, e.Number)
		}
		return strings.Join(out, " ")
	}

	t.Run("airing: the next one, then newest first", func(t *testing.T) {
		next, eps := episodeList(airing, "releasing", now)
		if next == nil || next.Label != "Ep. 3" || next.In != "in 2d 0h" {
			t.Errorf("next = %+v, want Ep. 3 in 2d 0h", next)
		}
		if got := labels(eps); got != "Ep. 3 Ep. 2 Ep. 1" {
			t.Errorf("list = %s, want Ep. 3 Ep. 2 Ep. 1", got)
		}
	})
	t.Run("finished: in order", func(t *testing.T) {
		_, eps := episodeList(airing[:2], "finished", now)
		if got := labels(eps); got != "Ep. 1 Ep. 2" {
			t.Errorf("list = %s, want Ep. 1 Ep. 2", got)
		}
	})
	t.Run("bare numbers aren't worth a list", func(t *testing.T) {
		next, eps := episodeList([]sqlcgen.EntryUnitsRow{{Label: "Ep. 1"}, {Label: "Ep. 2"}}, "finished", now)
		if next != nil || eps != nil {
			t.Errorf("got %+v, %+v; want nothing", next, eps)
		}
	})
}
