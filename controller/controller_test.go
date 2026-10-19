package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/web"
)

func TestRoutes(t *testing.T) {
	bundle := &web.Bundle{
		Files: fstest.MapFS{"assets/app-abc.css": {Data: []byte("body{}")}},
		Tags:  web.Tags{CSS: []string{"/assets/app-abc.css"}, Script: "/assets/islands-abc.js"},
	}
	up := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}, fakeEntries{}).srv.Handler)
	t.Cleanup(up.Close)
	down := httptest.NewServer(newHTTPServer("", bundle, fakeDB{err: errors.New("connection refused")}, fakeEntries{}).srv.Handler)
	t.Cleanup(down.Close)

	cases := []struct {
		name   string
		srv    *httptest.Server
		path   string
		status int
		want   []string
	}{
		{"home", up, "/", http.StatusOK, []string{"hello world", `data-island="wave"`, `{&#34;name&#34;:&#34;world&#34;}`, "/assets/islands-abc.js"}},
		{"there", up, "/there", http.StatusOK, []string{"hello there", `data-island="wave"`}},
		{"unknown path is a real 404", up, "/nope", http.StatusNotFound, []string{"back home"}},
		{"there has no subpaths", up, "/there/x", http.StatusNotFound, []string{"back home"}},
		{"hashed asset", up, "/assets/app-abc.css", http.StatusOK, []string{"body{}"}},
		{"missing asset", up, "/assets/gone.js", http.StatusNotFound, nil},
		{"meta", up, "/api/meta", http.StatusOK, []string{`"module":`, `"version":`, `"database":"ok"`, `"systemTime":"`}},
		{"meta still answers with the database down", down, "/api/meta", http.StatusOK, []string{`"database":"unreachable"`}},
		{"unrouted api is the blackhole, not a page", up, "/api/nope", http.StatusNotFound, []string{"DATA NOT FOUND"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(tc.srv.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tc.status)
			}
			for _, w := range tc.want {
				if !strings.Contains(string(body), w) {
					t.Errorf("body missing %q", w)
				}
			}
		})
	}
}

// fakeDB stands in for the pool; the API only pings it so far.
type fakeDB struct{ err error }

func (f fakeDB) Ping(context.Context) error {
	return f.err
}

func TestEntryPage(t *testing.T) {
	bundle := &web.Bundle{Files: fstest.MapFS{}}
	srv := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}, fakeEntries{}).srv.Handler)
	t.Cleanup(srv.Close)

	// redirects are the point here, so don't follow them
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	cases := []struct {
		name     string
		path     string
		status   int
		location string
		want     []string
	}{
		{name: "canonical", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{"Steins;Gate", "anime &middot; JP", "Apr 2011 &middot; 24 episodes",
				`<link rel="canonical" href="/entries/7/steins-gate-anime">`}},
		{name: "a film shows no units", path: "/entries/8/steins-gate-movie", status: http.StatusOK,
			want: []string{"Load Region", "Apr 20, 2013"}},
		{name: "wrong slug redirects", path: "/entries/7/steins-gate", status: http.StatusMovedPermanently,
			location: "/entries/7/steins-gate-anime"},
		{name: "missing slug redirects", path: "/entries/7", status: http.StatusMovedPermanently,
			location: "/entries/7/steins-gate-anime"},
		{name: "redirect keeps the query", path: "/entries/7/x?tab=relations", status: http.StatusMovedPermanently,
			location: "/entries/7/steins-gate-anime?tab=relations"},
		{name: "unknown id", path: "/entries/999/steins-gate-anime", status: http.StatusNotFound, want: []string{"back home"}},
		{name: "id is not a number", path: "/entries/steins-gate-anime", status: http.StatusNotFound, want: []string{"back home"}},
		{name: "id zero", path: "/entries/0/x", status: http.StatusNotFound, want: []string{"back home"}},
		{name: "database error is a 500, not a 404", path: "/entries/500/x", status: http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Get(srv.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tc.status)
			}
			if got := res.Header.Get("Location"); got != tc.location {
				t.Errorf("Location = %q, want %q", got, tc.location)
			}
			for _, w := range tc.want {
				if !strings.Contains(string(body), w) {
					t.Errorf("body missing %q", w)
				}
			}
		})
	}
}

// fakeEntries serves two seed entries by public id; 500 fails.
type fakeEntries struct{}

func (fakeEntries) EntryByPublicID(_ context.Context, id int64) (sqlcgen.EntryByPublicIDRow, error) {
	jp, month, day, total := "JP", "month", "day", int32(24)
	switch id {
	case 7:
		start := time.Date(2011, 4, 1, 0, 0, 0, 0, time.UTC)
		return sqlcgen.EntryByPublicIDRow{
			PublicID: 7, Slug: "steins-gate-anime", Medium: "anime", Format: "tv", FormatLabel: "anime",
			Country: &jp, Title: "Steins;Gate", StartDate: &start, StartPrecision: &month,
			UnitKind: "episode", UnitTotal: &total,
		}, nil
	case 8:
		start := time.Date(2013, 4, 20, 0, 0, 0, 0, time.UTC)
		return sqlcgen.EntryByPublicIDRow{
			PublicID: 8, Slug: "steins-gate-movie", Medium: "anime", Format: "movie", FormatLabel: "film",
			Country: &jp, Title: "Steins;Gate: The Movie - Load Region of Déjà Vu", StartDate: &start,
			StartPrecision: &day, UnitKind: "none",
		}, nil
	case 500:
		return sqlcgen.EntryByPublicIDRow{}, errors.New("connection reset")
	default:
		return sqlcgen.EntryByPublicIDRow{}, pgx.ErrNoRows
	}
}
