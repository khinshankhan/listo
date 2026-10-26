package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data"
	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/platform/config"
)

func TestReadManifest(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    []imageSource
		wantErr string
	}{
		{"every column",
			"slug\tkind\tfile\tsource\tsource_url\n" +
				"mega-man-2\tcover\tmm2.jpg\twikipedia\thttps://en.wikipedia.org/wiki/Mega_Man_2\n",
			[]imageSource{{Slug: "mega-man-2", Kind: "cover", From: "/raw/mm2.jpg", Source: "wikipedia", SourceURL: "https://en.wikipedia.org/wiki/Mega_Man_2"}}, ""},
		{"columns in any order, blanks and comments skipped, crlf too",
			"file\tslug\r\n\n# later\nmm3.png\tmega-man-3\r\n",
			[]imageSource{{Slug: "mega-man-3", Kind: "cover", From: "/raw/mm3.png", Source: "staff"}}, ""},
		{"a url stays a url",
			"slug\tkind\tfile\nmega-man-x\tbanner\thttps://example.test/x.webp\n",
			[]imageSource{{Slug: "mega-man-x", Kind: "banner", From: "https://example.test/x.webp", Source: "staff"}}, ""},
		{"no file column", "slug\tkind\nmega-man-2\tcover\n", nil, "no file column"},
		{"a row with no file", "slug\tfile\nmega-man-2\n", nil, "line 2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readManifest(strings.NewReader(tc.in), "/raw")
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want one with %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v\nwant %+v", got, tc.want)
			}
		})
	}
}

func TestOpenSourceURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cover.png":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte("png bytes"))
		case "/page":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<html>"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	cases := []struct {
		path    string
		wantErr string
	}{
		{"/cover.png", ""},
		{"/page", `"text/html" isn't an image`},
		{"/gone.jpg", "status 404"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			rc, err := openSource(context.Background(), srv.URL+tc.path)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want one with %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer rc.Close()
			if b, _ := io.ReadAll(rc); string(b) != "png bytes" {
				t.Errorf("read %q", b)
			}
		})
	}
}

// pngOf is a flat image of the given size.
func pngOf(t *testing.T, w, h int) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewGray(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
	return &buf
}

func TestAddImageCrownsTheBestCover(t *testing.T) {
	if config.DatabaseURL() == "" && os.Getenv("PGHOST") == "" {
		t.Skip("no database configured: make db-reset, or set LISTO_DATABASE_URL")
	}
	ctx := context.Background()
	pool, err := data.Open(ctx, config.DatabaseURL())
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

	run := fmt.Sprint(time.Now().UnixNano())
	slug := "test-images-" + run
	medium := "test_medium_" + run
	setup := []struct {
		sql  string
		args []any
	}{
		{`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
		  values ($1, 'test medium', 'playing', 'plan to play', 'replaying')`, []any{medium}},
		{`insert into catalog.medium_format (medium, format, label) values ($1, 'full', 'game')`, []any{medium}},
		{`with n as (insert into catalog.node (type, slug) values ('entry', $1) returning id)
		  insert into catalog.entry (id, medium, format, title) select id, $2, 'full', 'Test Game' from n`, []any{slug, medium}},
	}
	for _, s := range setup {
		if _, err := tx.Exec(ctx, s.sql, s.args...); err != nil {
			t.Fatal(err)
		}
	}
	store := media.NewFileStore(t.TempDir(), "/media/")

	// canonKey is the entry's canon cover now
	canonKey := func() string {
		var key string
		err := tx.QueryRow(ctx, `select i.storage_key from catalog.entry_image i join catalog.node n on n.id = i.entry_id
		                         where n.slug = $1 and i.kind = 'cover' and i.is_canon`, slug).Scan(&key)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			t.Fatal(err)
		}
		return key
	}
	add := func(w, h int) added {
		t.Helper()
		res, err := addImage(ctx, tx, store, imageSource{Slug: slug, Kind: media.KindCover, Source: "staff"}, pngOf(t, w, h))
		if err != nil {
			t.Fatal(err)
		}
		return res
	}

	first := add(460, 644)
	if !first.New || !first.Canon || canonKey() != first.Key {
		t.Fatalf("the first cover: %+v, canon %q", first, canonKey())
	}
	bigger := add(920, 1288)
	if !bigger.Canon || canonKey() != bigger.Key {
		t.Errorf("a bigger cover the same shape should take over: %+v, canon %q", bigger, canonKey())
	}
	smaller := add(300, 420)
	if !smaller.New || smaller.Canon || canonKey() != bigger.Key {
		t.Errorf("a smaller cover the same shape shouldn't: %+v, canon %q", smaller, canonKey())
	}
	tiny := add(150, 210)
	if !tiny.New || tiny.Canon || canonKey() != bigger.Key {
		t.Errorf("a tiny cover shouldn't: %+v, canon %q", tiny, canonKey())
	}
	again := add(460, 644)
	if again.New || again.Canon || canonKey() != bigger.Key {
		t.Errorf("the first again should be a no-op: %+v, canon %q", again, canonKey())
	}
	var rows int
	if err := tx.QueryRow(ctx, `select count(*) from catalog.entry_image i join catalog.node n on n.id = i.entry_id where n.slug = $1`, slug).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 4 {
		t.Errorf("%d rows, want 4", rows)
	}

	if _, err := addImage(ctx, tx, store, imageSource{Slug: "no-such-entry-" + run, Kind: media.KindCover}, pngOf(t, 460, 644)); err == nil || !strings.Contains(err.Error(), "no entry") {
		t.Errorf("unknown slug: err = %v", err)
	}
}
