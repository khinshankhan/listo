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

func TestGroupPageQueries(t *testing.T) {
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
	node := func(typ, name string) (uuid.UUID, int64) {
		t.Helper()
		var id uuid.UUID
		var publicID int64
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ($1, $2) returning id, public_id`,
			typ, fmt.Sprintf("test-%s-%s", name, run)).Scan(&id, &publicID); err != nil {
			t.Fatal(err)
		}
		return id, publicID
	}
	entry := func(name, title, start string) uuid.UUID {
		t.Helper()
		id, _ := node("entry", name)
		exec(`insert into catalog.entry (id, medium, format, title, start_date, start_precision, unit_kind)
		      values ($1, $2, 'tv', $3, $4, 'day', 'episode')`, id, medium, title, start)
		return id
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)

	franchise, franchiseID := node("group", "franchise")
	sub, _ := node("group", "sub")
	exec(`insert into catalog.grp (id, kind, name) values ($1, 'franchise', 'Test Franchise'), ($2, 'series', 'Zz Series')`,
		franchise, sub)
	// added out of release order, to prove the query sorts them
	later := entry("later", "Aa Later", "2020-01-01")
	earlier := entry("earlier", "Bb Earlier", "2010-01-01")
	exec(`insert into catalog.membership (group_id, member_id, position_label) values
	      ($1, $2, '#2'), ($1, $3, '#1'), ($1, $4, null)`, franchise, later, earlier, sub)

	exec(`insert into catalog.unit (entry_id, ordinal, kind, label) values ($1, 1, 'episode', 'Ep. 1'), ($1, 2, 'episode', 'Ep. 2')`, earlier)
	var ordering uuid.UUID
	if err := tx.QueryRow(ctx, `insert into catalog.ordering (group_id, kind, name) values ($1, 'recommended', 'start here') returning id`,
		franchise).Scan(&ordering); err != nil {
		t.Fatal(err)
	}
	exec(`insert into catalog.ordering_item (ordering_id, position, entry_id, necessity, note, section)
	      values ($1, 2, $2, 'optional', 'if you liked it', 'the later arc')`, ordering, later)
	exec(`insert into catalog.ordering_item (ordering_id, position, entry_id, unit_from, unit_to)
	      select $1, 1, $2, (select id from catalog.unit where entry_id = $2 and ordinal = 1),
	             (select id from catalog.unit where entry_id = $2 and ordinal = 2)`, ordering, earlier)

	q := sqlcgen.New(tx)

	t.Run("node kind by public id", func(t *testing.T) {
		got, err := q.NodeByPublicID(ctx, franchiseID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Type != "group" || !strings.HasPrefix(got.Slug, "test-franchise-") {
			t.Errorf("got %+v, want the franchise group", got)
		}
	})

	t.Run("group by public id", func(t *testing.T) {
		got, err := q.GroupByPublicID(ctx, franchiseID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != franchise || got.Kind != "franchise" || got.Name != "Test Franchise" || got.Standing != nil {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("entries by release, each with its position", func(t *testing.T) {
		got, err := q.GroupEntries(ctx, franchise)
		if err != nil {
			t.Fatal(err)
		}
		var rows []string
		for _, e := range got {
			rows = append(rows, e.Title+" "+deref(e.PositionLabel))
		}
		// the empty series has no entries to bring
		if want := "Bb Earlier #1|Aa Later #2"; strings.Join(rows, "|") != want {
			t.Errorf("rows = %v, want %s", rows, want)
		}
	})

	t.Run("orderings by position, with unit labels", func(t *testing.T) {
		got, err := q.GroupOrderings(ctx, franchise)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("got %d steps, want 2: %+v", len(got), got)
		}
		first, second := got[0], got[1]
		if first.Title != "Bb Earlier" || first.UnitFrom == nil || *first.UnitFrom != "Ep. 1" || first.UnitTo == nil || *first.UnitTo != "Ep. 2" {
			t.Errorf("first step = %+v, want Bb Earlier, Ep. 1 to Ep. 2", first)
		}
		if second.Title != "Aa Later" || second.Necessity != "optional" || second.Note == nil || *second.Note != "if you liked it" {
			t.Errorf("second step = %+v, want Aa Later, optional, with its note", second)
		}
		if first.Section != nil || second.Section == nil || *second.Section != "the later arc" {
			t.Errorf("sections = %v, %v, want none, then the later arc", first.Section, second.Section)
		}
	})
}

func TestFranchiseQueries(t *testing.T) {
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
	node := func(typ, name string) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ($1, $2) returning id`,
			typ, fmt.Sprintf("test-%s-%s", name, run)).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	group := func(name, kind string, standing any) uuid.UUID {
		t.Helper()
		id := node("group", name)
		exec(`insert into catalog.grp (id, kind, name, standing) values ($1, $2, $3, $4)`, id, kind, name, standing)
		return id
	}
	entry := func(name string) uuid.UUID {
		t.Helper()
		id := node("entry", name)
		exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'tv', $3)`, id, medium, name)
		return id
	}
	member := func(g, m uuid.UUID, role any) {
		t.Helper()
		exec(`insert into catalog.membership (group_id, member_id, role) values ($1, $2, $3)`, g, m, role)
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)

	// universe > franchise > {main > series > nested, main > shared, alt > shared, loose}
	universe := group("universe", "universe", nil)
	franchise := group("franchise", "franchise", nil)
	main := group("main", "continuity", "primary")
	alt := group("alt", "continuity", "concluded")
	series := group("series", "series", nil)
	nested, shared, loose := entry("nested"), entry("shared"), entry("loose")
	member(universe, franchise, nil)
	member(franchise, main, nil)
	member(franchise, alt, nil)
	member(franchise, loose, nil)
	member(main, series, nil)
	member(series, nested, nil)
	member(main, shared, "source material")
	member(alt, shared, "retold")
	// a sibling franchise with two entries, one nested
	sibling := group("sibling", "franchise", nil)
	sibSeries := group("sibling-series", "series", nil)
	member(universe, sibling, nil)
	member(sibling, entry("sibling-a"), nil)
	member(sibling, sibSeries, nil)
	member(sibSeries, entry("sibling-b"), nil)

	q := sqlcgen.New(tx)

	t.Run("entries by branch", func(t *testing.T) {
		rows, err := q.GroupEntries(ctx, franchise)
		if err != nil {
			t.Fatal(err)
		}
		got := map[string]string{}
		for _, r := range rows {
			key := strings.TrimSuffix(strings.TrimPrefix(r.Slug, "test-"), "-"+run) + "@" + strings.TrimSuffix(strings.TrimPrefix(r.BranchSlug, "test-"), "-"+run)
			got[key] = deref(r.Role)
		}
		want := map[string]string{
			"nested@main": "",                // a series inside a continuity counts under the continuity
			"shared@main": "source material", // an entry in two continuities shows in both,
			"shared@alt":  "retold",          // with the role it has in each
			"loose@loose": "",                // directly in the franchise: its own branch
		}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for k, v := range want {
			if r, ok := got[k]; !ok || r != v {
				t.Errorf("%s: role %q (present %v), want %q", k, r, ok, v)
			}
		}
	})

	t.Run("siblings count every entry beneath", func(t *testing.T) {
		sibs, err := q.GroupSiblings(ctx, franchise)
		if err != nil {
			t.Fatal(err)
		}
		if len(sibs) != 1 || sibs[0].Name != "sibling" || sibs[0].Entries != 2 {
			t.Errorf("got %+v, want only the sibling, with 2 entries", sibs)
		}
	})
}

// A group's relations are the ones between its own entries, however deep, and
// none that reach outside it or land on a group.
func TestGroupRelations(t *testing.T) {
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
	node := func(typ, name string) (uuid.UUID, int64) {
		t.Helper()
		var id uuid.UUID
		var publicID int64
		if err := tx.QueryRow(ctx, `insert into catalog.node (type, slug) values ($1, $2) returning id, public_id`,
			typ, fmt.Sprintf("test-%s-%s", name, run)).Scan(&id, &publicID); err != nil {
			t.Fatal(err)
		}
		return id, publicID
	}
	entry := func(name string) (uuid.UUID, int64) {
		t.Helper()
		id, publicID := node("entry", name)
		exec(`insert into catalog.entry (id, medium, format, title, unit_kind) values ($1, $2, 'tv', $3, 'none')`,
			id, medium, name)
		return id, publicID
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)

	franchise, _ := node("group", "franchise")
	sub, _ := node("group", "sub")
	exec(`insert into catalog.grp (id, kind, name) values ($1, 'franchise', 'Test Franchise'), ($2, 'series', 'Test Series')`,
		franchise, sub)
	top, topID := entry("top")
	deep, deepID := entry("deep")
	outside, _ := entry("outside")
	exec(`insert into catalog.membership (group_id, member_id) values ($1, $2), ($1, $3), ($3, $4)`, franchise, top, sub, deep)
	exec(`insert into catalog.relation (from_node, kind, to_node, spoiler) values
	      ($1, 'sequel_of', $2, true), ($2, 'references', $3, false)`, deep, top, outside)
	exec(`insert into catalog.relation (from_node, kind, to_node, is_symmetric) values (least($1::uuid, $2::uuid), 'crossover_with', greatest($1::uuid, $2::uuid), true)`,
		deep, sub)

	got, err := sqlcgen.New(tx).GroupRelations(ctx, franchise)
	if err != nil {
		t.Fatal(err)
	}
	want := []sqlcgen.GroupRelationsRow{{FromID: deepID, Kind: "sequel_of", Label: "sequel of", Family: "sequence", ToID: topID, Spoiler: true}}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
