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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/web"
)

func TestRoutes(t *testing.T) {
	bundle := &web.Bundle{
		Files: fstest.MapFS{"assets/app-abc.css": {Data: []byte("body{}")}},
		Tags:  web.Tags{CSS: []string{"/assets/app-abc.css"}, Script: "/assets/islands-abc.js"},
	}
	up := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(up.Close)
	down := httptest.NewServer(newHTTPServer("", bundle, fakeDB{err: errors.New("connection refused")}, fakeCatalog{}).srv.Handler)
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
	srv := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}, fakeCatalog{}).srv.Handler)
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
			want: []string{"Steins;Gate", "Apr 2011 &middot; 24 episodes",
				`<link rel="canonical" href="/entries/7/steins-gate-anime">`,
				`<meta name="description" content="Steins;Gate, an anime from 2011.">`}},
		{name: "format that only repeats the medium is left out", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">anime &middot; JP</p>"}},
		{name: "a film names its format and shows no units", path: "/entries/8/steins-gate-movie", status: http.StatusOK,
			want: []string{"Load Region", "anime &middot; film &middot; JP", "Apr 20, 2013",
				`content="Steins;Gate: The Movie - Load Region of Déjà Vu, an anime film from 2013."`}},
		{name: "native title is tagged with its language", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`lang="ja">シュタインズ・ゲート</p>`}},
		{name: "relations link to their entries", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`href="/entries/6/steins-gate-vn">Steins;Gate</a>`, "visual novel · 2009",
				`href="/entries/11/steins-gate-0-anime">Steins;Gate 0</a>`, "branches off episode 23"}},
		{name: "groups, nearest first", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{"#1</span> in Main timeline", `part of</span> Steins;Gate`}},
		{name: "every title and the links out", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">synonym</span>", "S;G", `href="https://anilist.co/anime/9253"`, ">9253</span>"}},
		{name: "a failed secondary read is a 500, not a half page", path: "/entries/501/broken",
			status: http.StatusInternalServerError},
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

// fakeCatalog serves two seed entries by public id. 500 fails outright, 501
// fails on its relations.
type fakeCatalog struct{}

var (
	sgAnime = uuid.MustParse("01a151f2-0000-7000-8000-000000000007")
	sgMovie = uuid.MustParse("01a151f2-0000-7000-8000-000000000008")
	broken  = uuid.MustParse("01a151f2-0000-7000-8000-000000000501")
)

func (fakeCatalog) EntryByPublicID(_ context.Context, id int64) (sqlcgen.EntryByPublicIDRow, error) {
	jp, month, day, total := "JP", "month", "day", int32(24)
	switch id {
	case 7:
		start := time.Date(2011, 4, 1, 0, 0, 0, 0, time.UTC)
		return sqlcgen.EntryByPublicIDRow{
			ID: sgAnime, PublicID: 7, Slug: "steins-gate-anime", Medium: "anime", MediumLabel: "anime",
			Format: "tv", FormatLabel: "anime", Country: &jp, Title: "Steins;Gate", StartDate: &start,
			StartPrecision: &month, UnitKind: "episode", UnitTotal: &total,
		}, nil
	case 8:
		start := time.Date(2013, 4, 20, 0, 0, 0, 0, time.UTC)
		return sqlcgen.EntryByPublicIDRow{
			ID: sgMovie, PublicID: 8, Slug: "steins-gate-movie", Medium: "anime", MediumLabel: "anime",
			Format: "movie", FormatLabel: "film", Country: &jp,
			Title: "Steins;Gate: The Movie - Load Region of Déjà Vu", StartDate: &start,
			StartPrecision: &day, UnitKind: "none",
		}, nil
	case 500:
		return sqlcgen.EntryByPublicIDRow{}, errors.New("connection reset")
	case 501:
		return sqlcgen.EntryByPublicIDRow{ID: broken, PublicID: 501, Slug: "broken", Medium: "anime",
			MediumLabel: "anime", Format: "tv", FormatLabel: "anime", Title: "Broken", UnitKind: "none"}, nil
	default:
		return sqlcgen.EntryByPublicIDRow{}, pgx.ErrNoRows
	}
}

func (fakeCatalog) NodeTitles(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeTitlesRow, error) {
	if id != sgAnime {
		return nil, nil
	}
	return []sqlcgen.NodeTitlesRow{
		{Kind: "english", Title: "Steins;Gate"},
		{Kind: "native", Title: "シュタインズ・ゲート"},
		{Kind: "synonym", Title: "S;G"},
	}, nil
}

func (fakeCatalog) NodeLinks(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeLinksRow, error) {
	if id != sgAnime {
		return nil, nil
	}
	return []sqlcgen.NodeLinksRow{
		{Site: "AniList", Source: "anilist_anime", ExternalID: "9253", Url: "https://anilist.co/anime/9253"},
	}, nil
}

func (fakeCatalog) NodeRelations(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeRelationsRow, error) {
	switch id {
	case sgAnime:
		vn, year := "visual novel", time.Date(2009, 10, 15, 0, 0, 0, 0, time.UTC)
		note := "branches off episode 23"
		return []sqlcgen.NodeRelationsRow{
			{Reading: "adaptation of", PublicID: 6, Slug: "steins-gate-vn", Type: "entry",
				Title: "Steins;Gate", MediumLabel: &vn, StartDate: &year},
			{Reading: "prequel of", PublicID: 11, Slug: "steins-gate-0-anime", Type: "entry",
				Title: "Steins;Gate 0", Note: &note},
		}, nil
	case broken:
		return nil, errors.New("connection reset")
	}
	return nil, nil
}

func (fakeCatalog) NodeGroups(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeGroupsRow, error) {
	if id != sgAnime {
		return nil, nil
	}
	first := "#1"
	return []sqlcgen.NodeGroupsRow{
		{PublicID: 3, Slug: "steins-gate-main", Kind: "continuity", Name: "Main timeline", PositionLabel: &first, Depth: 1},
		{PublicID: 2, Slug: "steins-gate", Kind: "franchise", Name: "Steins;Gate", Depth: 2},
	}, nil
}
