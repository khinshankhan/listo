package data_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data"
	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/config"
)

// testTx is a dev database transaction rolled back after the test.
func testTx(t *testing.T) pgx.Tx {
	t.Helper()

	url := config.DatabaseURL()
	if url == "" && os.Getenv("PGHOST") == "" {
		t.Skip("no database configured: make db-reset, or set LISTO_DATABASE_URL")
	}

	ctx := context.Background()
	pool, err := data.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})
	return tx
}

func TestEntryByPublicID(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()

	// per-run names, so a seed row can never be what the test finds
	run := time.Now().UnixNano()
	medium := fmt.Sprintf("test_medium_%d", run)
	slug := fmt.Sprintf("test-entry-%d", run)

	setup := []string{
		`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
		 values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`,
		`insert into catalog.medium_format (medium, format, label) values ($1, 'movie', 'film')`,
	}
	for _, q := range setup {
		if _, err := tx.Exec(ctx, q, medium); err != nil {
			t.Fatal(err)
		}
	}
	var id int64
	if err := tx.QueryRow(ctx, `
		with n as (insert into catalog.node (type, slug) values ('entry', $1) returning id, public_id),
		     e as (insert into catalog.entry (id, medium, format, country, title, start_date, start_precision, unit_kind)
		           select id, $2, 'movie', 'JP', 'Test Film', '2013-04-01', 'month', 'none' from n)
		select public_id from n`,
		slug, medium).Scan(&id); err != nil {
		t.Fatal(err)
	}

	q := sqlcgen.New(tx)

	t.Run("found", func(t *testing.T) {
		got, err := q.EntryByPublicID(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if got.PublicID != id || got.Slug != slug || got.Medium != medium || got.Format != "movie" || got.Title != "Test Film" {
			t.Errorf("got %+v", got)
		}
		// the format's label comes from the vocabulary row, not the format key
		if got.FormatLabel != "film" {
			t.Errorf("FormatLabel = %q, want %q", got.FormatLabel, "film")
		}
		if got.Country == nil || *got.Country != "JP" {
			t.Errorf("Country = %v, want JP", got.Country)
		}
		if want := time.Date(2013, 4, 1, 0, 0, 0, 0, time.UTC); got.StartDate == nil || !got.StartDate.Equal(want) {
			t.Errorf("StartDate = %v, want %v", got.StartDate, want)
		}
		if got.UnitTotal != nil {
			t.Errorf("UnitTotal = %v, want nil", *got.UnitTotal)
		}
	})

	t.Run("unknown id", func(t *testing.T) {
		// far past the id this test just took, and identity values only go up
		_, err := q.EntryByPublicID(ctx, id+1_000_000)
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Errorf("err = %v, want pgx.ErrNoRows", err)
		}
	})
}

func TestEntryPageQueries(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()

	// relation kinds are letters and underscores only
	run := strings.Map(func(r rune) rune { return 'a' + (r - '0') }, fmt.Sprint(time.Now().UnixNano()))
	medium, kind, source := "test_medium_"+run, "test_kind_"+run, "test_source_"+run
	slug := func(name string) string { return fmt.Sprintf("test-%s-%s", name, run) }

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
			typ, slug(name)).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)
	exec(`insert into catalog.relation_kind (kind, label, inverse_label, is_symmetric, family)
	      values ($1, 'adaptation of', 'adapted as', false, 'derivative')`, kind)
	exec(`insert into catalog.id_source (source, site, url_template) values ($1, 'Test Site', 'https://example.test/{id}/')`, source)

	entry := func(name, title string) uuid.UUID {
		t.Helper()
		id := node("entry", name)
		exec(`insert into catalog.entry (id, medium, format, title) values ($1, $2, 'tv', $3)`, id, medium, title)
		return id
	}
	source1, adaptation := entry("source", "The Source"), entry("adaptation", "The Adaptation")
	exec(`insert into catalog.relation (from_node, kind, to_node) values ($1, $2, $3)`, adaptation, kind, source1)

	franchise, films, specials := node("group", "franchise"), node("group", "films"), node("group", "specials")
	exec(`insert into catalog.grp (id, kind, name) values
	      ($1, 'franchise', 'Test Franchise'), ($2, 'series', 'Test Series'), ($3, 'series', 'Test Specials')`,
		franchise, films, specials)
	exec(`insert into catalog.membership (group_id, member_id, position_label) values ($1, $2, 'Films'), ($1, $3, 'Specials')`,
		franchise, films, specials)
	exec(`insert into catalog.membership (group_id, member_id, position_label) values ($1, $2, '#3'), ($3, $2, '#1')`,
		films, adaptation, specials)

	exec(`insert into catalog.title (node_id, kind, title) values
	      ($1, 'synonym', 'TA'), ($1, 'native', 'ジ・アダプテーション'), ($1, 'english', 'The Adaptation')`, adaptation)
	exec(`insert into catalog.external_id (source, external_id, node_id) values ($1, '42', $2)`, source, adaptation)

	q := sqlcgen.New(tx)

	t.Run("relations read from each side", func(t *testing.T) {
		cases := []struct {
			name      string
			from      uuid.UUID
			wantRead  string
			wantTitle string
		}{
			{"the adaptation sees its source", adaptation, "adaptation of", "The Source"},
			{"the source sees the inverse", source1, "adapted as", "The Adaptation"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := q.NodeRelations(ctx, tc.from)
				if err != nil {
					t.Fatal(err)
				}
				if len(got) != 1 || got[0].Reading != tc.wantRead || got[0].Title != tc.wantTitle {
					t.Errorf("got %+v, want one %q %q", got, tc.wantRead, tc.wantTitle)
				}
			})
		}
	})

	t.Run("groups walk up, nearest first", func(t *testing.T) {
		got, err := q.NodeGroups(ctx, adaptation)
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, g := range got {
			names = append(names, g.Name)
		}
		// the franchise is reached through both series, and listed once
		if want := "Test Series|Test Specials|Test Franchise"; strings.Join(names, "|") != want {
			t.Fatalf("groups = %v, want %s", names, want)
		}
		if got[0].Depth != 1 || got[0].PositionLabel == nil || *got[0].PositionLabel != "#3" {
			t.Errorf("nearest = %+v, want depth 1, #3", got[0])
		}
		// a series' position in the franchise is the series', not the entry's
		if got[2].Depth != 2 || got[2].PositionLabel != nil {
			t.Errorf("franchise = %+v, want depth 2, no position", got[2])
		}
	})

	t.Run("titles, likeliest first", func(t *testing.T) {
		got, err := q.NodeTitles(ctx, adaptation)
		if err != nil {
			t.Fatal(err)
		}
		var kinds []string
		for _, r := range got {
			kinds = append(kinds, r.Kind)
		}
		if want := "english native synonym"; strings.Join(kinds, " ") != want {
			t.Errorf("kinds = %v, want %s", kinds, want)
		}
	})

	t.Run("links fill the url template", func(t *testing.T) {
		got, err := q.NodeLinks(ctx, adaptation)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].Url != "https://example.test/42/" || got[0].Site != "Test Site" {
			t.Errorf("got %+v", got)
		}
	})
}
