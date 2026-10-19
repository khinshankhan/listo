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
	exec(`insert into catalog.ordering_item (ordering_id, position, entry_id, necessity, note) values ($1, 2, $2, 'optional', 'if you liked it')`,
		ordering, later)
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

	t.Run("members: subgroups first, then entries by release", func(t *testing.T) {
		got, err := q.GroupMembers(ctx, franchise)
		if err != nil {
			t.Fatal(err)
		}
		var titles []string
		for _, m := range got {
			titles = append(titles, m.Title)
		}
		if want := "Zz Series|Bb Earlier|Aa Later"; strings.Join(titles, "|") != want {
			t.Errorf("titles = %v, want %s", titles, want)
		}
		if got[1].PositionLabel == nil || *got[1].PositionLabel != "#1" {
			t.Errorf("earlier entry's position = %v, want #1", got[1].PositionLabel)
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
	})
}
