package data_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

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
