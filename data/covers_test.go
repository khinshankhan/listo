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

// Every query that lists entries carries one cover each: the canon one, else
// the widest, else none.
func TestCoverQueries(t *testing.T) {
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
	entry := func(name string) (uuid.UUID, int64) {
		t.Helper()
		var id uuid.UUID
		var public int64
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ('entry', $1) returning id, public_id`,
			fmt.Sprintf("test-%s-%s", name, run)).Scan(&id, &public); err != nil {
			t.Fatal(err)
		}
		exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'tv', $3)`, id, medium, name)
		exec(`insert into catalog.title (node_id, kind, title) values ($1, 'english', $2)`, id, name)
		return id, public
	}
	image := func(entry uuid.UUID, kind, key string, width int, canon bool) {
		t.Helper()
		exec(`insert into catalog.entry_image (entry_id, kind, storage_key, width, height, fit, is_canon)
		      values ($1, $2, $3, $4, $4 * 7 / 5, 'cover', $5)`, entry, kind, key, width, canon)
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)

	_, barePub := entry("bare")
	canon, canonPub := entry("canon")
	image(canon, "cover", "cover/aa/aa/wide", 900, false)
	image(canon, "cover", "cover/bb/bb/canon", 300, true)
	widest, widestPub := entry("widest")
	image(widest, "cover", "cover/cc/cc/narrow", 200, false)
	image(widest, "cover", "cover/dd/dd/wide", 460, false)
	// a banner is never a cover, canon or not
	banner, bannerPub := entry("banner")
	image(banner, "banner", "banner/ee/ee/canon", 1920, true)

	var group uuid.UUID
	if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ('group', $1) returning id`, "test-group-"+run).Scan(&group); err != nil {
		t.Fatal(err)
	}
	exec(`insert into catalog.grp (id, kind, name) values ($1, 'series', 'Test Series')`, group)
	var ordering uuid.UUID
	if err := tx.QueryRow(ctx, `insert into catalog.ordering (group_id, kind, name) values ($1, 'release', 'release') returning id`, group).Scan(&ordering); err != nil {
		t.Fatal(err)
	}
	exec(`insert into catalog.membership (group_id, member_id) select $1, e.id from catalog.entry e where e.medium = $2`, group, medium)
	exec(`insert into catalog.ordering_item (ordering_id, position, entry_id)
	      select $1, row_number() over (order by e.title), e.id from catalog.entry e where e.medium = $2`, ordering, medium)

	want := map[int64]string{barePub: "", canonPub: "cover/bb/bb/canon", widestPub: "cover/dd/dd/wide", bannerPub: ""}
	key := func(k *string) string {
		if k == nil {
			return ""
		}
		return *k
	}
	check := func(t *testing.T, public int64, got *string) {
		t.Helper()
		w, ok := want[public]
		if !ok {
			return
		}
		if key(got) != w {
			t.Errorf("entry %d: cover %q, want %q", public, key(got), w)
		}
	}

	q := sqlcgen.New(tx)
	t.Run("entry page", func(t *testing.T) {
		for public := range want {
			e, err := q.EntryByPublicID(ctx, public)
			if err != nil {
				t.Fatal(err)
			}
			check(t, public, e.CoverKey)
		}
	})
	t.Run("browse", func(t *testing.T) {
		rows, err := q.BrowseEntries(ctx, sqlcgen.BrowseEntriesParams{Medium: &medium, PageSize: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != len(want) {
			t.Fatalf("%d rows, want %d", len(rows), len(want))
		}
		for _, r := range rows {
			check(t, r.PublicID, r.CoverKey)
		}
	})
	t.Run("shelves", func(t *testing.T) {
		rows, err := q.BrowseShelves(ctx, 10)
		if err != nil {
			t.Fatal(err)
		}
		seen := 0
		for _, r := range rows {
			if r.Medium == medium {
				seen++
				check(t, r.PublicID, r.CoverKey)
			}
		}
		if seen != len(want) {
			t.Errorf("%d rows on the test shelf, want %d", seen, len(want))
		}
	})
	t.Run("group", func(t *testing.T) {
		rows, err := q.GroupEntries(ctx, group)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != len(want) {
			t.Fatalf("%d rows, want %d", len(rows), len(want))
		}
		for _, r := range rows {
			check(t, r.PublicID, r.CoverKey)
		}
	})
	t.Run("ordering", func(t *testing.T) {
		rows, err := q.GroupOrderings(ctx, group)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != len(want) {
			t.Fatalf("%d rows, want %d", len(rows), len(want))
		}
		for _, r := range rows {
			check(t, r.PublicID, r.CoverKey)
		}
	})
	t.Run("relations", func(t *testing.T) {
		// relation kinds allow letters and underscores only
		kind := "test_sequel_" + strings.Map(func(r rune) rune { return 'a' + (r - '0') }, run)
		exec(`insert into catalog.relation_kind (kind, label, inverse_label, is_symmetric, family)
		      values ($1, 'sequel of', 'prequel of', false, 'sequence')`, kind)
		for _, other := range []uuid.UUID{canon, widest, banner} {
			exec(`insert into catalog.relation (from_node, kind, to_node) values ($1, $2, $3)`, other, kind, group)
		}
		rows, err := q.NodeRelations(ctx, group)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 3 {
			t.Fatalf("%d rows, want 3", len(rows))
		}
		for _, r := range rows {
			check(t, r.PublicID, r.CoverKey)
		}
	})
	t.Run("search", func(t *testing.T) {
		rows, err := q.SearchNodes(ctx, sqlcgen.SearchNodesParams{Q: "canon", Medium: &medium, MaxResults: 10})
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) == 0 {
			t.Fatal("no hits")
		}
		for _, r := range rows {
			check(t, r.PublicID, r.CoverKey)
		}
	})
}
