package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
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
		{"home", up, "/", http.StatusOK, []string{"the 1 list for all your content", "/assets/islands-abc.js"}},
		{"the demo page is gone", up, "/there", http.StatusNotFound, []string{"back home"}},
		{"every page has the top bar", up, "/nope", http.StatusNotFound, []string{`role="search"`, `id="site-search"`}},
		{"unknown path is a real 404", up, "/nope", http.StatusNotFound, []string{"back home"}},
		{"there has no subpaths", up, "/there/x", http.StatusNotFound, []string{"back home"}},
		{"hashed asset", up, "/assets/app-abc.css", http.StatusOK, []string{"body{}"}},
		{"missing asset", up, "/assets/gone.js", http.StatusNotFound, nil},
		{"assets aren't a directory listing", up, "/assets/", http.StatusNotFound, nil},
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

			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := plainLinks(string(raw))
			if res.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tc.status)
			}
			for _, w := range tc.want {
				if !strings.Contains(body, w) {
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

func TestCatalogPages(t *testing.T) {
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
		dontWant []string
		once     []string
	}{
		{name: "canonical", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{"Steins;Gate</h1>",
				`<link rel="canonical" href="/entries/7/steins-gate-anime">`,
				`<meta name="description" content="Steins;Gate, an anime from 2011.">`}},
		{name: "facts, with the start named for the medium", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">format</dt>", ">episodes</dt>", ">24</dd>", ">aired</dt>", ">Apr 2011</dd>", ">country</dt>", ">JP</dd>"}},
		{name: "breadcrumb browses the medium", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`aria-label="breadcrumb"`, `href="/entries?medium=anime">anime</a>`}},
		{name: "a tv anime's eyebrow has its season, and the format only when it adds something", path: "/entries/7/steins-gate-anime",
			status: http.StatusOK, want: []string{`uppercase">spring 2011</span>`},
			once: []string{`text-foreground-subtle uppercase">anime</span>`}},
		{name: "tags", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`aria-label="tags"`, ">sci-fi</li>", ">thriller</li>"}},
		{name: "a film names its format, is released, and shows no units", path: "/entries/8/steins-gate-movie", status: http.StatusOK,
			want: []string{"Load Region", `uppercase">film</span>`, `uppercase">2013</span>`, ">released</dt>", ">Apr 20, 2013</dd>",
				`content="Steins;Gate: The Movie - Load Region of Déjà Vu, an anime film from 2013."`},
			dontWant: []string{">episodes</dt>", `aria-label="tags"`}},
		{name: "a link out is never a script", path: "/entries/8/steins-gate-movie", status: http.StatusOK,
			want: []string{">Bad</a>"}, dontWant: []string{`href="javascript:`}},
		{name: "native title is tagged with its language", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`lang="ja">シュタインズ・ゲート</dd>`}},
		{name: "relations link to their entries", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`href="/entries/6/steins-gate-vn">Steins;Gate</a>`, "visual novel · 2009",
				`href="/entries/11/steins-gate-0-anime">Steins;Gate 0</a>`, "branches off episode 23"}},
		{name: "groups, nearest first, linked", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{`#1</span> in <a href="/groups/3/steins-gate-main">Main timeline</a>`,
				`part of</span> <a href="/groups/2/steins-gate">Steins;Gate</a>`}},
		{name: "a group's id under /entries moves to its page", path: "/entries/2/steins-gate", status: http.StatusMovedPermanently,
			location: "/groups/2/steins-gate"},
		{name: "group", path: "/groups/2/steins-gate", status: http.StatusOK,
			dontWant: []string{"w-12 shrink-0"},
			want: []string{`uppercase">franchise`, "Steins;Gate</h1>", `href="/groups/3/steins-gate-main">Main timeline</a>`, ">continuity</span>",
				`<link rel="canonical" href="/groups/2/steins-gate">`, `content="Steins;Gate, a franchise on listo."`}},
		{name: "a universe, not an universe", path: "/groups/1/science-adventure", status: http.StatusOK,
			want: []string{`content="Science Adventure, a universe on listo."`}},
		{name: "a continuity shows its standing", path: "/groups/3/steins-gate-main", status: http.StatusOK,
			want: []string{"w-12 shrink-0", "continuity &middot; primary", `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, "anime · 2011"}},
		{name: "an ordering, position by position", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{">recommended order</h2>", "first-timer order", "Watch the series, then the branch",
				">01</span>", "Ep. 1 to Ep. 22", ">02</span>", ">23β</span>", ">optional</span>", "epilogue"}},
		{name: "a fork keeps each choice whole, in order", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{">pick one</p>", ">the branch</li>",
				// second choice: the film, then the sequel
				`href="/entries/8/steins-gate-movie">Steins;Gate: The Movie</a>`,
				`href="/entries/11/steins-gate-0-anime">Steins;Gate 0</a>`},
			dontWant: []string{">or</span>"}},
		{name: "a relation says where it attaches", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">here at Ep. 23</span>"}},
		{name: "group wrong slug redirects", path: "/groups/2/sg?x=1", status: http.StatusMovedPermanently,
			location: "/groups/2/steins-gate?x=1"},
		{name: "group missing slug redirects", path: "/groups/3", status: http.StatusMovedPermanently,
			location: "/groups/3/steins-gate-main"},
		{name: "an entry's id under /groups moves to its page", path: "/groups/7", status: http.StatusMovedPermanently,
			location: "/entries/7/steins-gate-anime"},
		{name: "unknown group", path: "/groups/999/x", status: http.StatusNotFound, want: []string{"back home"}},
		{name: "group database error is a 500", path: "/groups/500/x", status: http.StatusInternalServerError},
		{name: "every title and the links out", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">synonym</dt>", ">S;G</dd>", `href="https://anilist.co/anime/9253"`, ">9253</span>"},
			// the heading already shows the display title
			dontWant: []string{">english</dt>"}},
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

			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := plainLinks(string(raw))
			if res.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tc.status)
			}
			if got := res.Header.Get("Location"); got != tc.location {
				t.Errorf("Location = %q, want %q", got, tc.location)
			}
			for _, w := range tc.want {
				if !strings.Contains(body, w) {
					t.Errorf("body missing %q", w)
				}
			}
			for _, w := range tc.dontWant {
				if strings.Contains(body, w) {
					t.Errorf("body has %q", w)
				}
			}
			for _, w := range tc.once {
				if n := strings.Count(body, w); n != 1 {
					t.Errorf("body has %q %d times, want once", w, n)
				}
			}
		})
	}
}

func TestFrontPages(t *testing.T) {
	bundle := &web.Bundle{Files: fstest.MapFS{}}
	srv := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(srv.Close)

	cases := []struct {
		name     string
		path     string
		status   int
		want     []string
		dontWant []string
	}{
		{name: "home searches and starts browsing", path: "/", status: http.StatusOK,
			want: []string{`<form action="/search" method="get" role="search"`,
				`href="/groups/2/steins-gate">Steins;Gate</a>`, "franchise &middot; 14 entries",
				`href="/entries?medium=anime"`, `href="/entries"`},
			// home is top bar only
			dontWant: []string{"data-island", "noindex", `aria-label="catalog"`}},
		{name: "empty search points at the box", path: "/search", status: http.StatusOK,
			want: []string{`id="site-search" name="q"`, `<meta name="robots" content="noindex">`, `aria-label="catalog"`},
			// no medium is current while searching
			dontWant: []string{">results</h2>", `aria-current`}},
		{name: "search finds groups and entries", path: "/search?q=steins", status: http.StatusOK,
			want: []string{`value="steins"`, ">results</h2>", `href="/groups/2/steins-gate">Steins;Gate</a>`,
				`href="/entries/7/steins-gate-anime">Steins;Gate</a>`, "anime · 2011", ">franchise</span>"}},
		{name: "a hit on an alias says so", path: "/search?q=sg", status: http.StatusOK,
			want: []string{"as S;G</span>"}},
		{name: "no hits", path: "/search?q=nothing", status: http.StatusOK,
			want: []string{"nothing matches &ldquo;nothing&rdquo; yet."}},
		{name: "the query is escaped", path: "/search?q=%3Cscript%3E", status: http.StatusOK,
			want: []string{"&lt;script&gt;"}, dontWant: []string{`value="<script>`, "&ldquo;<script>"}},
		{name: "search error is a 500", path: "/search?q=boom", status: http.StatusInternalServerError},
		{name: "browse everything", path: "/entries", status: http.StatusOK,
			want: []string{`href="/entries" aria-current="page"`, ">everything</span>", `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, `href="/entries/6/steins-gate-vn">Steins;Gate</a>`,
				`<link rel="canonical" href="/entries">`}},
		{name: "browse one medium", path: "/entries?medium=anime", status: http.StatusOK,
			want: []string{`href="/entries/7/steins-gate-anime">Steins;Gate</a>`,
				`<link rel="canonical" href="/entries?medium=anime">`, `href="/entries?medium=anime" aria-current="page"`},
			dontWant: []string{"/entries/6/steins-gate-vn"}},
		{name: "unknown medium", path: "/entries?medium=pottery", status: http.StatusNotFound, want: []string{"back home"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(srv.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := plainLinks(string(raw))
			if res.StatusCode != tc.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tc.status)
			}
			for _, w := range tc.want {
				if !strings.Contains(body, w) {
					t.Errorf("body missing %q", w)
				}
			}
			for _, w := range tc.dontWant {
				if strings.Contains(body, w) {
					t.Errorf("body has %q", w)
				}
			}
		})
	}
}

// fakeCatalog serves a few seed nodes by public id. 500 fails outright, 501
// fails on its relations.
type fakeCatalog struct{}

var (
	saUniverse  = uuid.MustParse("01a151f2-0000-7000-8000-000000000001")
	sgFranchise = uuid.MustParse("01a151f2-0000-7000-8000-000000000002")
	sgMain      = uuid.MustParse("01a151f2-0000-7000-8000-000000000003")
	sgOrder     = uuid.MustParse("01a151f2-0000-7000-8000-00000000f001")
	sgAnime     = uuid.MustParse("01a151f2-0000-7000-8000-000000000007")
	sgMovie     = uuid.MustParse("01a151f2-0000-7000-8000-000000000008")
	broken      = uuid.MustParse("01a151f2-0000-7000-8000-000000000501")
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

func (fakeCatalog) NodeByPublicID(_ context.Context, id int64) (sqlcgen.NodeByPublicIDRow, error) {
	nodes := map[int64]sqlcgen.NodeByPublicIDRow{
		2: {PublicID: 2, Slug: "steins-gate", Type: "group"},
		3: {PublicID: 3, Slug: "steins-gate-main", Type: "group"},
		7: {PublicID: 7, Slug: "steins-gate-anime", Type: "entry"},
		8: {PublicID: 8, Slug: "steins-gate-movie", Type: "entry"},
	}
	if n, ok := nodes[id]; ok {
		return n, nil
	}
	return sqlcgen.NodeByPublicIDRow{}, pgx.ErrNoRows
}

func (fakeCatalog) GroupByPublicID(_ context.Context, id int64) (sqlcgen.GroupByPublicIDRow, error) {
	primary := "primary"
	switch id {
	case 1:
		return sqlcgen.GroupByPublicIDRow{ID: saUniverse, PublicID: 1, Slug: "science-adventure", Kind: "universe",
			Name: "Science Adventure"}, nil
	case 2:
		return sqlcgen.GroupByPublicIDRow{ID: sgFranchise, PublicID: 2, Slug: "steins-gate", Kind: "franchise", Name: "Steins;Gate"}, nil
	case 3:
		return sqlcgen.GroupByPublicIDRow{ID: sgMain, PublicID: 3, Slug: "steins-gate-main", Kind: "continuity",
			Name: "Main timeline", Standing: &primary}, nil
	case 500:
		return sqlcgen.GroupByPublicIDRow{}, errors.New("connection reset")
	}
	return sqlcgen.GroupByPublicIDRow{}, pgx.ErrNoRows
}

func (fakeCatalog) GroupMembers(_ context.Context, id uuid.UUID) ([]sqlcgen.GroupMembersRow, error) {
	continuity, anime, first := "continuity", "anime", "#1"
	start := time.Date(2011, 4, 6, 0, 0, 0, 0, time.UTC)
	switch id {
	case sgFranchise:
		return []sqlcgen.GroupMembersRow{
			{PublicID: 3, Slug: "steins-gate-main", Type: "group", Title: "Main timeline", GroupKind: &continuity},
		}, nil
	case sgMain:
		return []sqlcgen.GroupMembersRow{
			{PublicID: 7, Slug: "steins-gate-anime", Type: "entry", Title: "Steins;Gate", MediumLabel: &anime, StartDate: &start,
				PositionLabel: &first},
		}, nil
	}
	return nil, nil
}

func (fakeCatalog) GroupOrderings(_ context.Context, id uuid.UUID) ([]sqlcgen.GroupOrderingsRow, error) {
	if id != sgFranchise {
		return nil, nil
	}
	ep1, ep22, beta := "Ep. 1", "Ep. 22", "23β"
	rationale, note := "Watch the series, then the branch the sequel follows.", "epilogue"
	step := func(pos, branch int32, publicID int64, slug, title string) sqlcgen.GroupOrderingsRow {
		return sqlcgen.GroupOrderingsRow{OrderingID: sgOrder, Kind: "recommended", Name: "first-timer order",
			Rationale: &rationale, Position: pos, Branch: branch, Step: 1, PublicID: publicID, Slug: slug,
			Title: title, Necessity: "required"}
	}
	section := "the branch"
	first := step(1, 1, 7, "steins-gate-anime", "Steins;Gate")
	first.UnitFrom, first.UnitTo = &ep1, &ep22
	// position 2 forks: the 23β episode alone, or the film then the sequel
	second := step(2, 1, 7, "steins-gate-anime", "Steins;Gate")
	second.UnitFrom, second.UnitTo, second.Section = &beta, &beta, &section
	film := step(2, 2, 8, "steins-gate-movie", "Steins;Gate: The Movie")
	film.Necessity, film.Note, film.Section = "optional", &note, &section
	sequel := step(2, 2, 11, "steins-gate-0-anime", "Steins;Gate 0")
	sequel.Step, sequel.Section = 2, &section
	return []sqlcgen.GroupOrderingsRow{first, second, film, sequel}, nil
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
	switch id {
	case sgAnime:
		return []sqlcgen.NodeLinksRow{
			{Site: "AniList", Source: "anilist_anime", ExternalID: "9253", Url: "https://anilist.co/anime/9253"},
		}, nil
	case sgMovie:
		// an id_source url_template gone wrong
		return []sqlcgen.NodeLinksRow{{Site: "Bad", Source: "bad", ExternalID: "1", Url: "javascript:alert(1)"}}, nil
	}
	return nil, nil
}

func (fakeCatalog) NodeRelations(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeRelationsRow, error) {
	switch id {
	case sgAnime:
		vn, year := "visual novel", time.Date(2009, 10, 15, 0, 0, 0, 0, time.UTC)
		note, ep23 := "branches off episode 23", "Ep. 23"
		return []sqlcgen.NodeRelationsRow{
			{Reading: "adaptation of", PublicID: 6, Slug: "steins-gate-vn", Type: "entry",
				Title: "Steins;Gate", MediumLabel: &vn, StartDate: &year},
			{Reading: "prequel of", PublicID: 11, Slug: "steins-gate-0-anime", Type: "entry",
				Title: "Steins;Gate 0", Note: &note, MyUnit: &ep23},
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

func (fakeCatalog) SearchNodes(_ context.Context, arg sqlcgen.SearchNodesParams) ([]sqlcgen.SearchNodesRow, error) {
	franchise, anime := "franchise", "anime"
	start := time.Date(2011, 4, 6, 0, 0, 0, 0, time.UTC)
	switch arg.Q {
	case "steins":
		return []sqlcgen.SearchNodesRow{
			{PublicID: 2, Slug: "steins-gate", Type: "group", Title: "Steins;Gate", Matched: "Steins;Gate", GroupKind: &franchise},
			{PublicID: 7, Slug: "steins-gate-anime", Type: "entry", Title: "Steins;Gate", Matched: "Steins;Gate",
				MediumLabel: &anime, StartDate: &start},
		}, nil
	case "sg":
		return []sqlcgen.SearchNodesRow{
			{PublicID: 7, Slug: "steins-gate-anime", Type: "entry", Title: "Steins;Gate", Matched: "S;G", MediumLabel: &anime},
		}, nil
	case "boom":
		return nil, errors.New("connection reset")
	}
	return nil, nil
}

func (fakeCatalog) TopGroups(context.Context) ([]sqlcgen.TopGroupsRow, error) {
	return []sqlcgen.TopGroupsRow{{PublicID: 2, Slug: "steins-gate", Kind: "franchise", Name: "Steins;Gate", Entries: 14}}, nil
}

func (fakeCatalog) MediumCounts(context.Context) ([]sqlcgen.MediumCountsRow, error) {
	return []sqlcgen.MediumCountsRow{
		{Medium: "anime", Label: "anime", Entries: 1},
		{Medium: "visual_novel", Label: "visual novel", Entries: 1},
	}, nil
}

func (fakeCatalog) BrowseEntries(_ context.Context, medium *string) ([]sqlcgen.BrowseEntriesRow, error) {
	all := []sqlcgen.BrowseEntriesRow{
		{PublicID: 7, Slug: "steins-gate-anime", Title: "Steins;Gate", MediumLabel: "anime"},
		{PublicID: 6, Slug: "steins-gate-vn", Title: "Steins;Gate", MediumLabel: "visual novel"},
	}
	if medium == nil {
		return all, nil
	}
	if *medium == "anime" {
		return all[:1], nil
	}
	return nil, nil
}

func (fakeCatalog) NodeTags(_ context.Context, id uuid.UUID) ([]sqlcgen.NodeTagsRow, error) {
	if id != sgAnime {
		return nil, nil
	}
	return []sqlcgen.NodeTagsRow{{Kind: "genre", Name: "sci-fi"}, {Kind: "genre", Name: "thriller"}}, nil
}

var (
	anchorTag = regexp.MustCompile(`<a href="([^"]*)"([^>]*)>`)
	linkIcon  = regexp.MustCompile(`<svg class="mb-0[^"]*"[^>]*>.*?</svg>\s*<span class="sr-only"> \((external|file) link\)</span>`)
)

// plainLinks reduces anchors to their href (and aria-current) and drops link
// icons, so page tests check where links go rather than their classes.
func plainLinks(body string) string {
	body = linkIcon.ReplaceAllString(body, "")
	return anchorTag.ReplaceAllStringFunc(body, func(tag string) string {
		m := anchorTag.FindStringSubmatch(tag)
		if strings.Contains(m[2], `aria-current="page"`) {
			return fmt.Sprintf(`<a href="%s" aria-current="page">`, m[1])
		}
		return fmt.Sprintf(`<a href="%s">`, m[1])
	})
}
