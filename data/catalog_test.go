package data_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/khinshankhan/listo/data/sqlcgen"
)

func TestCatalogQueries(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()

	run := fmt.Sprint(time.Now().UnixNano())
	medium := "test_medium_" + run

	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := tx.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	node := func(typ, name string) (uuid.UUID, string) {
		t.Helper()
		slug := fmt.Sprintf("test-%s-%s", name, run)
		var id uuid.UUID
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ($1, $2) returning id`, typ, slug).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id, slug
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)

	// titles nothing in the seed resembles, so hits are ours
	entry, entrySlug := node("entry", "entry")
	exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'tv', 'Zorblax Quintessence')`, entry, medium)
	exec(`insert into catalog.title (node_id, kind, title) values
	      ($1, 'english', 'Zorblax Quintessence'), ($1, 'synonym', 'Kwazimodo Frunt'), ($1, 'native', 'Zörblax Quintessénce')`, entry)

	top, topSlug := node("group", "top")
	nested, nestedSlug := node("group", "nested")
	exec(`insert into catalog.grp (id, kind, name) values ($1, 'franchise', 'Zorblax'), ($2, 'series', 'Zorblax Shorts')`, top, nested)
	exec(`insert into catalog.membership (group_id, member_id) values ($1, $2), ($1, $3)`, top, nested, entry)
	// one more entry a level down, which the top group still counts
	deep, _ := node("entry", "deep")
	exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'tv', 'Zorblax Shorts: One')`, deep, medium)
	exec(`insert into catalog.membership (group_id, member_id) values ($1, $2)`, nested, deep)

	q := sqlcgen.New(tx)

	t.Run("search", func(t *testing.T) {
		cases := []struct {
			query, wantMatched string
		}{
			{"zorblax quintessence", "Zorblax Quintessence"},
			{"zorblax quintesence", "Zorblax Quintessence"}, // a typo
			{"kwazimodo frunt", "Kwazimodo Frunt"},          // an alias
			{"quintessence", "Zorblax Quintessence"},        // part of a title
		}
		for _, tc := range cases {
			t.Run(tc.query, func(t *testing.T) {
				hits, err := q.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: tc.query, MaxResults: 50})
				if err != nil {
					t.Fatal(err)
				}
				found := 0
				for _, h := range hits {
					if h.Slug != entrySlug {
						continue
					}
					found++
					// titles fold the same, so any matching one will do
					if h.Title != "Zorblax Quintessence" || h.Type != "entry" || h.MediumLabel == nil || *h.MediumLabel != "test medium" {
						t.Errorf("hit = %+v", h)
					}
					if tc.wantMatched == "Kwazimodo Frunt" && h.Matched != tc.wantMatched {
						t.Errorf("Matched = %q, want %q", h.Matched, tc.wantMatched)
					}
				}
				// three titles match "zorblax quintessence"; the entry still counts once
				if found != 1 {
					t.Errorf("entry found %d times, want 1", found)
				}
			})
		}
	})

	t.Run("ties come back the same way: groups, then oldest first", func(t *testing.T) {
		// its own medium, so the browse counts below don't see these
		plimMedium := medium + "_plim"
		exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
		      values ($1, 'plim medium', 'watching', 'plan to watch', 'rewatching')`, plimMedium)
		exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, plimMedium)
		// added newest first, so insertion order can't pass for the rule
		newer, _ := node("entry", "plim-2011")
		older, _ := node("entry", "plim-2009")
		group, _ := node("group", "plim")
		exec(`insert into catalog.entry (id, medium, format, title, start_date, start_precision) values
		      ($1, $3, 'tv', 'Plimvort', '2011-04-01', 'month'), ($2, $3, 'tv', 'Plimvort', '2009-10-01', 'month')`, newer, older, plimMedium)
		exec(`insert into catalog.grp (id, kind, name) values ($1, 'franchise', 'Plimvort')`, group)
		exec(`insert into catalog.title (node_id, kind, title) values ($1, 'english', 'Plimvort'), ($2, 'english', 'Plimvort'), ($3, 'english', 'Plimvort')`,
			newer, older, group)

		for range 5 {
			hits, err := q.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: "plimvort", MaxResults: 10})
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, h := range hits {
				if h.Title == "Plimvort" {
					got = append(got, strings.TrimSuffix(h.Slug, "-"+run))
				}
			}
			if want := []string{"test-plim", "test-plim-2009", "test-plim-2011"}; strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("order = %v, want %v", got, want)
			}
		}
	})

	t.Run("a medium narrows the search before the limit", func(t *testing.T) {
		// media of its own, so the browse counts below don't see these
		first, other := medium+"_first", medium+"_other"
		for _, m := range []string{first, other} {
			exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
			      values ($1, 'a medium', 'watching', 'plan to watch', 'rewatching')`, m)
			exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, m)
		}
		// tied on score, so A sorts first and fills a limit of one
		a, _ := node("entry", "quorvex-a")
		b, _ := node("entry", "quorvex-b")
		exec(`insert into catalog.entry (id, medium, format, title) values ($1, $3, 'tv', 'Quorvex A'), ($2, $4, 'tv', 'Quorvex B')`, a, b, first, other)
		exec(`insert into catalog.title (node_id, kind, title) values ($1, 'english', 'Quorvex A'), ($2, 'english', 'Quorvex B')`, a, b)

		hits, err := q.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: "quorvex", Medium: &other, MaxResults: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(hits) != 1 || hits[0].Title != "Quorvex B" {
			t.Errorf("got %+v, want Quorvex B", hits)
		}

		counts, err := q.SearchMedia(ctx, "quorvex")
		if err != nil {
			t.Fatal(err)
		}
		got := map[string]int32{}
		for _, c := range counts {
			if c.Medium != nil {
				got[*c.Medium] = c.Hits
			}
		}
		if got[first] != 1 || got[other] != 1 {
			t.Errorf("counts = %v, want one in each medium", got)
		}
	})

	t.Run("search honours the limit", func(t *testing.T) {
		hits, err := q.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: "zorblax", MaxResults: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(hits) != 1 {
			t.Errorf("got %d hits, want 1", len(hits))
		}
	})

	t.Run("top groups leave out nested ones", func(t *testing.T) {
		groups, err := q.TopGroups(ctx)
		if err != nil {
			t.Fatal(err)
		}
		var sawTop, sawNested bool
		for _, g := range groups {
			switch g.Slug {
			case topSlug:
				sawTop = true
				if g.Entries != 2 {
					t.Errorf("top group entries = %d, want 2 (one direct, one nested)", g.Entries)
				}
			case nestedSlug:
				sawNested = true
			}
		}
		if !sawTop || sawNested {
			t.Errorf("saw top = %v, nested = %v; want only the top group", sawTop, sawNested)
		}
	})

	t.Run("browse filters by medium", func(t *testing.T) {
		rows, err := q.BrowseEntries(ctx, sqlcgen.BrowseEntriesParams{Medium: &medium, PageSize: 100})
		if err != nil {
			t.Fatal(err)
		}
		// undated, so by title
		if len(rows) != 2 || rows[0].Slug != entrySlug {
			t.Errorf("got %+v, want the two test entries, by title", rows)
		}

		counts, err := q.MediumCounts(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range counts {
			if c.Medium == medium && c.Entries != 2 {
				t.Errorf("count for test medium = %d, want 2", c.Entries)
			}
		}
	})

	t.Run("tags", func(t *testing.T) {
		exec(`insert into catalog.tag (kind, name) values ('genre', $1), ('theme', $2)`, "zz-genre-"+run, "aa-theme-"+run)
		exec(`insert into catalog.node_tag (node_id, tag_id) select $1, id from catalog.tag where name in ($2, $3)`,
			entry, "zz-genre-"+run, "aa-theme-"+run)

		got, err := q.NodeTags(ctx, entry)
		if err != nil {
			t.Fatal(err)
		}
		// by kind, then name
		if len(got) != 2 || got[0].Kind != "genre" || got[1].Kind != "theme" {
			t.Errorf("got %+v, want the genre then the theme", got)
		}
	})
}

func TestBrowseQueries(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()

	run := fmt.Sprint(time.Now().UnixNano())
	medium := "test_medium_" + run
	genre, unused := "test-genre-"+run, "test-unused-"+run

	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := tx.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	entry := func(name, start, precision string) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ('entry', $1) returning id`,
			fmt.Sprintf("test-%s-%s", name, run)).Scan(&id); err != nil {
			t.Fatal(err)
		}
		var date, prec any
		if start != "" {
			date, prec = start, precision
		}
		exec(`insert into catalog.entry (id, medium, format, title, start_date, start_precision) values ($1, $2, 'tv', $3, $4, $5)`,
			id, medium, name, date, prec)
		return id
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)
	spring := entry("spring", "2011-04-01", "month")
	entry("year-only", "2011-01-01", "year")
	entry("fall", "2018-10-06", "day")
	entry("undated", "", "")
	exec(`insert into catalog.tag (kind, name) values ('genre', $1), ('genre', $2)`, genre, unused)
	exec(`insert into catalog.node_tag (node_id, tag_id) select $1, id from catalog.tag where name = $2`, spring, genre)

	q := sqlcgen.New(tx)
	year := func(y int32) *int32 { return &y }
	season := func(s int32) *int32 { return &s }

	cases := []struct {
		name string
		args sqlcgen.BrowseEntriesParams
		want string
	}{
		{"newest first, undated last", sqlcgen.BrowseEntriesParams{}, "fall spring year-only undated"},
		{"a genre", sqlcgen.BrowseEntriesParams{Genre: &genre}, "spring"},
		{"a year", sqlcgen.BrowseEntriesParams{Year: year(2011)}, "spring year-only"},
		// a year-only date can't be placed in a season
		{"a season", sqlcgen.BrowseEntriesParams{Season: season(1)}, "spring"},
		{"fall", sqlcgen.BrowseEntriesParams{Season: season(3)}, "fall"},
		// stored as January 1st, but that's not a winter release
		{"a year-only date isn't winter", sqlcgen.BrowseEntriesParams{Season: season(0)}, ""},
		{"filters combine", sqlcgen.BrowseEntriesParams{Year: year(2011), Season: season(3)}, ""},
		{"a page", sqlcgen.BrowseEntriesParams{PageSize: 2, Skip: 1}, "spring year-only"},
		{"past the end", sqlcgen.BrowseEntriesParams{PageSize: 2, Skip: 4}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// scoped to the test medium, so the seed's entries stay out of it
			tc.args.Medium = &medium
			if tc.args.PageSize == 0 {
				tc.args.PageSize = 100
			}
			rows, err := q.BrowseEntries(ctx, tc.args)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, r := range rows {
				got = append(got, r.Title)
			}
			if strings.Join(got, " ") != tc.want {
				t.Errorf("got %q, want %q", strings.Join(got, " "), tc.want)
			}
		})
	}

	t.Run("each row of a page knows the whole count", func(t *testing.T) {
		rows, err := q.BrowseEntries(ctx, sqlcgen.BrowseEntriesParams{Medium: &medium, PageSize: 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 2 || rows[0].Total != 4 || rows[1].Total != 4 {
			t.Errorf("got %d rows, totals %v; want 2 rows of 4", len(rows), rows)
		}
	})

	t.Run("shelves hold a medium's newest few", func(t *testing.T) {
		rows, err := q.BrowseShelves(ctx, 2)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, r := range rows {
			if r.Medium == medium {
				got = append(got, r.Title)
			}
		}
		if strings.Join(got, " ") != "fall spring" {
			t.Errorf("got %v, want fall then spring", got)
		}
	})

	t.Run("filters only offer what something has", func(t *testing.T) {
		genres, err := q.BrowseGenres(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(genres, genre) || slices.Contains(genres, unused) {
			t.Errorf("genres = %v, want %s and not %s", genres, genre, unused)
		}
		years, err := q.BrowseYears(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(years, 2018) || !slices.Contains(years, 2011) || !slices.IsSortedFunc(years, func(a, b int32) int { return int(b - a) }) {
			t.Errorf("years = %v, want 2018 and 2011, newest first", years)
		}
	})

	// last, since it adds an entry the counts above don't expect
	t.Run("a format narrows its medium, and only formats with entries are offered", func(t *testing.T) {
		exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'ova', 'OVA'), ($1, 'special', 'special')`, medium)
		var id uuid.UUID
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ('entry', $1) returning id`, "test-ova-"+run).Scan(&id); err != nil {
			t.Fatal(err)
		}
		exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'ova', 'ova')`, id, medium)
		format := "ova"
		rows, err := q.BrowseEntries(ctx, sqlcgen.BrowseEntriesParams{Medium: &medium, Format: &format, PageSize: 100})
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 1 || rows[0].Title != "ova" {
			t.Errorf("got %+v, want just the ova", rows)
		}
		formats, err := q.BrowseFormats(ctx, medium)
		if err != nil {
			t.Fatal(err)
		}
		want := []sqlcgen.BrowseFormatsRow{{Format: "ova", Label: "OVA"}, {Format: "tv", Label: "tv"}}
		if !slices.Equal(formats, want) {
			t.Errorf("formats = %+v, want %+v", formats, want)
		}
	})
}
