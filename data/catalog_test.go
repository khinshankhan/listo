package data_test

import (
	"context"
	"fmt"
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
		rows, err := q.BrowseEntries(ctx, &medium)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 2 || rows[0].Slug != entrySlug {
			t.Errorf("got %+v, want the two test entries, alphabetically", rows)
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
