package controller

import (
	"context"
	"encoding/json"
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
		{"home", up, "/", http.StatusOK, []string{"keep one list for everything you watch, read and play.", "/assets/islands-abc.js"}},
		{"the demo page is gone", up, "/there", http.StatusNotFound, []string{"back home"}},
		{"every page has the top bar", up, "/nope", http.StatusNotFound, []string{`role="search"`, `id="site-search"`}},
		// the first stop for a keyboard, past the top bar and sidebar
		{"every page skips to its content", up, "/nope", http.StatusNotFound, []string{`<a href="#content">skip to content</a>`, `<main id="content"`}},
		// web/src/search.ts finds these by id and drives the aria attributes
		{"the search box is a combobox over a hidden panel", up, "/nope", http.StatusNotFound, []string{
			`role="combobox" aria-autocomplete="list" aria-expanded="false" aria-controls="site-search-results"`,
			`<div id="site-search-panel" hidden`, `<ul id="site-search-results" role="listbox"`, `id="site-search-all"`,
			`<p id="site-search-status" role="status" class="sr-only"></p>`, `<p id="site-search-note" hidden`}},
		{"unknown path is a real 404", up, "/nope", http.StatusNotFound, []string{"back home"}},
		{"there has no subpaths", up, "/there/x", http.StatusNotFound, []string{"back home"}},
		{"hashed asset", up, "/assets/app-abc.css", http.StatusOK, []string{"body{}"}},
		{"every page names its icon", up, "/nope", http.StatusNotFound, []string{`<link rel="icon" href="/favicon.svg" type="image/svg+xml">`}},
		{"the icon", up, "/favicon.svg", http.StatusOK, []string{`<svg xmlns="http://www.w3.org/2000/svg"`}},
		// followed to the svg, so the default request isn't a 404
		{"favicon.ico finds the icon", up, "/favicon.ico", http.StatusOK, []string{`<svg xmlns="http://www.w3.org/2000/svg"`}},
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
			once: []string{`leading-[14px] tracking-wide text-foreground-muted uppercase">anime</span>`}},
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
		{name: "a franchise gets its own page", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{"Steins;Gate</h1>", `aria-current="page">Steins;Gate</span>`,
				`uppercase">universe</span> <a href="/groups/1/science-adventure">Science Adventure</a>`, ">1 entry</span>", ">1 medium</span>",
				">1 continuity</span>", `<link rel="canonical" href="/groups/2/steins-gate">`,
				`content="Steins;Gate, a franchise on listo: 1 entry, 1 medium, 1 continuity."`},
			dontWant: []string{">members</h2>"}},
		{name: "a continuity card says how it stands", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{`id="continuities"`, `href="/groups/3/steins-gate-main">Main timeline</a>`, ">primary</span>",
				">the one to start with</span>"}},
		{name: "entries by continuity, with their part and when", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{`id="entries"`, `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, ">adaptation of the VN</span>",
				">spring 2011 · 24 episodes</span>", `href="/groups/2/steins-gate" aria-current="page">by continuity</a>`,
				`href="/groups/2/steins-gate?by=release">by release</a>`}},
		{name: "a franchise can group by medium", path: "/groups/2/steins-gate?by=medium", status: http.StatusOK,
			want: []string{`href="/groups/2/steins-gate?by=medium" aria-current="page">by medium</a>`}},
		{name: "an unknown grouping is the default", path: "/groups/2/steins-gate?by=vibes", status: http.StatusOK,
			want: []string{`href="/groups/2/steins-gate" aria-current="page">by continuity</a>`}},
		{name: "start here jumps to the orderings", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{`href="#start-here">`, `id="start-here"`, ">recommended order</h2>"}},
		{name: "what sits beside it", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{`id="related"`, `href="/groups/5/robotics-notes">Robotics;Notes</a>`, "franchise · 2 entries", ">same universe</span>"}},
		{name: "the sidebar knows the page's sections", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{"on this page</p>", `href="#continuities">`, `href="#entries">`, `href="#related">`}},
		{name: "a universe, not an universe", path: "/groups/1/science-adventure", status: http.StatusOK,
			want: []string{`content="Science Adventure, a universe on listo."`}},
		// nothing seeded under the universe itself in the fakes
		{name: "an empty group says so, and has no stats to show", path: "/groups/1/science-adventure", status: http.StatusOK,
			want:     []string{"Nothing in Science Adventure yet."},
			dontWant: []string{"0 entries"}},
		{name: "every kind of group gets the same page", path: "/groups/3/steins-gate-main", status: http.StatusOK,
			want: []string{`aria-current="page">Main timeline</span>`, ">continuity</span>", "Main timeline</h1>", ">primary</span>",
				`content="Main timeline, a continuity on listo: 1 entry, 1 medium."`}},
		{name: "a series' entries keep their positions", path: "/groups/3/steins-gate-main", status: http.StatusOK,
			want: []string{">#1</span>", `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, ">spring 2011 · 24 episodes</span>"}},
		// entries all straight in the group: nothing to group them by
		{name: "a group without sub-groups lists by release", path: "/groups/3/steins-gate-main", status: http.StatusOK,
			want:     []string{`href="/groups/3/steins-gate-main" aria-current="page">by release</a>`, `href="/groups/3/steins-gate-main?by=medium">by medium</a>`},
			dontWant: []string{"by continuity", "by group", "directly in"}},
		{name: "an old grouping falls back to release", path: "/groups/3/steins-gate-main?by=continuity", status: http.StatusOK,
			want: []string{`href="/groups/3/steins-gate-main" aria-current="page">by release</a>`}},
		{name: "an ordering, position by position", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{">recommended order</h2>", "first-timer order", "Watch the series, then the branch",
				">01</span>", "Ep. 1 to Ep. 22", ">02</span>", ">23β</span>", ">optional</span>", "epilogue"}},
		{name: "a fork keeps each choice whole, in order", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{">pick one</p>", ">the branch</li>",
				// second choice: the film, then the sequel
				`href="/entries/8/steins-gate-movie">Steins;Gate: The Movie</a>`,
				`href="/entries/11/steins-gate-0-anime">Steins;Gate 0</a>`},
			dontWant: []string{">or</span>"}},
		{name: "an ordering links to its own page", path: "/groups/2/steins-gate", status: http.StatusOK,
			want: []string{`href="/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order">first-timer order</a>`}},
		{name: "an ordering's page", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order", status: http.StatusOK,
			want: []string{"first-timer order</h1>", `<link rel="canonical" href="/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order">`,
				`<a href="/groups/2/steins-gate">Steins;Gate</a>`, `aria-current="page">ordering</span>`,
				">2 steps</span>", ">2 required</span>", ">maintained by staff</span>",
				">why this order</p>", "Watch the series, then the branch"}},
		{name: "an ordering's steps, forks and notes", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order", status: http.StatusOK,
			want: []string{">pick one</p>", ">optional</span>", "<span>epilogue</span>", ">the branch</span>",
				`href="/entries/11/steins-gate-0-anime">Steins;Gate 0</a>`, ">required</span>"},
			dontWant: []string{"more steps"}},
		{name: "a first-timer's ordering points rewatchers elsewhere", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order", status: http.StatusOK,
			want: []string{"rewatching, or already know the story?", "release order goes in the order they came out.",
				`href="/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f002/release-order">view release order</a>`}},
		{name: "and the others point first-timers back", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f002/release-order", status: http.StatusOK,
			want: []string{"watching for the first time?", `href="/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order">view first-timer order</a>`,
				">auto-seeded</span>"}},
		{name: "the sidebar lists the group's orderings", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f002/release-order", status: http.StatusOK,
			want: []string{"orderings</p>", ">recommended · staff</span>", `href="/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f002/release-order" aria-current="page"`,
				`href="/groups/2/steins-gate" aria-current="page"`}},
		{name: "an ordering follows its group's slug", path: "/groups/2/wrong/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order", status: http.StatusMovedPermanently,
			location: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order"},
		{name: "an ordering the group doesn't have", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000ffff/vibes", status: http.StatusNotFound},
		{name: "an ordering's id isn't its name", path: "/groups/2/steins-gate/orderings/first-timer-order", status: http.StatusNotFound},
		{name: "an ordering's name is decoration", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/wrong", status: http.StatusMovedPermanently,
			location: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order"},
		{name: "or can be left off", path: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001", status: http.StatusMovedPermanently,
			location: "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order"},
		// names that slugify alike, or to nothing, each keep a page
		{name: "two names, one slug, two pages", path: "/groups/1/science-adventure/orderings/01a151f2-0000-7000-8000-00000000f201/release-order",
			status: http.StatusOK, want: []string{"Release order</h1>"}},
		{name: "the other of the two", path: "/groups/1/science-adventure/orderings/01a151f2-0000-7000-8000-00000000f202/release-order",
			status: http.StatusOK, want: []string{"Release-order</h1>"}},
		{name: "a name with no slug", path: "/groups/1/science-adventure/orderings/01a151f2-0000-7000-8000-00000000f203",
			status: http.StatusOK, want: []string{"初見順</h1>", `<link rel="canonical" href="/groups/1/science-adventure/orderings/01a151f2-0000-7000-8000-00000000f203">`}},
		{name: "an ordering of no group", path: "/groups/999/x/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order", status: http.StatusNotFound},
		// the form sends "any" as empty; one address per view
		{name: "browse drops empty filters", path: "/entries?medium=&genre=sci-fi&year=&season=", status: http.StatusFound,
			location: "/entries?genre=sci-fi"},
		{name: "all empty is plain browse", path: "/entries?medium=&genre=&year=&season=", status: http.StatusFound,
			location: "/entries"},
		{name: "page one is the plain address", path: "/entries?genre=many&page=1", status: http.StatusFound,
			location: "/entries?genre=many"},
		{name: "an entry tells its own story", path: "/entries/10/airing", status: http.StatusOK,
			want: []string{"A show still on the air, with a description long enough", `<meta name="description" content="A show still on the air`,
				`at a word and say so.</p>`},
			// the head's copy is clipped at a word, the page's is whole
			dontWant: []string{`say so."`}},
		{name: "what airs next, and how soon", path: "/entries/10/airing", status: http.StatusOK,
			want: []string{">next to air</span>", ">Ep. 3 in 3d 4h</span>"}},
		{name: "who made it, how long, and where it's at", path: "/entries/10/airing", status: http.StatusOK,
			want: []string{">studio</dt>", ">Studio A, Studio B</dd>", ">publisher</dt>", ">Pub</dd>", ">duration</dt>", ">24 min</dd>",
				">status</dt>", ">releasing</dd>", ">source</dt>", ">light novel</dd>"}},
		{name: "episodes, newest first while airing", path: "/entries/10/airing", status: http.StatusOK,
			want: []string{`id="episodes"`, "Ep. 3</span>", "not yet aired", ">The Middle</span>", ">2 days ago</span>",
				">The Start</span>", ">1 week ago</span>"},
			// only the next one ahead, not every scheduled one
			dontWant: []string{"Ep. 4"}},
		{name: "a studio, and nothing made up", path: "/entries/7/steins-gate-anime", status: http.StatusOK,
			want: []string{">White Fox</dd>"}, dontWant: []string{`id="episodes"`, "next to air", ">status</dt>", ">source</dt>"}},
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
		{name: "home says what listo is", path: "/", status: http.StatusOK,
			want: []string{"keep one list for everything you watch, read and play.</h1>",
				`<form action="/search" method="get" role="search" class="flex min-w-0 items-center rounded-lg border border-border focus-within:outline-2 focus-within:outline-stark-contrast relative h-[52px]`},
			// top bar only, and the one search is the hero's
			dontWant: []string{"data-island", "noindex", `aria-label="catalog"`, "max-w-[420px]"}},
		{name: "home's media: browsable, or coming soon", path: "/", status: http.StatusOK,
			want: []string{`href="/entries?medium=anime">anime</a>`, `href="/entries?medium=visual_novel">visual novel</a>`,
				`book<span class="sr-only">, coming soon</span>`}},
		{name: "home previews a franchise guide", path: "/", status: http.StatusOK,
			want: []string{`<a href="/groups/2/steins-gate">Steins;Gate</a></h3>`, "<span>1 entry</span>", "alternate episode 23β",
				">Main timeline</span>", `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, ">anime · 2011</span>"}},
		{name: "home previews an ordering guide", path: "/", status: http.StatusOK,
			want: []string{`<a href="/groups/103/monogatari/orderings/01a151f2-0000-7000-8000-00000000f103/first-timer-order">Monogatari</a></h3>`,
				"<span>6 steps</span>", "<span>68 episodes</span>", `href="/entries/104/bakemonogatari">Bakemonogatari</a>`,
				"a prequel, but not first</span>", ">recommended</span>", "preview · 5 of 6 steps"},
			// the sixth step stays on the ordering's page
			dontWant: []string{"Owarimonogatari"}},
		{name: "a medium with no entries isn't browsable yet", path: "/entries", status: http.StatusOK,
			dontWant: []string{"medium=book"}},
		{name: "empty search points at the box", path: "/search", status: http.StatusOK,
			want: []string{`id="site-search" name="q"`, `<meta name="robots" content="noindex">`, `aria-label="catalog"`},
			// no medium is current while searching
			dontWant: []string{">results</h2>", `aria-current`}},
		{name: "search finds groups and entries", path: "/search?q=steins", status: http.StatusOK,
			want: []string{`id="q" name="q" type="search" value="steins"`, ">3 results</span>",
				`href="/groups/2/steins-gate">Steins;Gate</a>`, `href="/entries/7/steins-gate-anime">Steins;Gate</a>`,
				">franchise</span>", ">anime</span>", ">spring 2011 · 24 episodes</span>",
				`<a href="/search">clear</a>`},
			// the page has the box, so the top bar's stays empty
			dontWant: []string{`id="site-search" name="q" type="search" value="steins"`, "matched"}},
		{name: "a hit on an alias says so", path: "/search?q=sg", status: http.StatusOK,
			want: []string{"matched alias: S;G</p>"}},
		{name: "results filter by medium", path: "/search?q=steins", status: http.StatusOK,
			want: []string{`href="/search?q=steins" aria-current="page">all</a>`, `href="/search?medium=anime&amp;q=steins">anime</a>`,
				`href="/search?medium=visual_novel&amp;q=steins">visual novel</a>`}},
		{name: "a medium keeps only its hits", path: "/search?q=steins&medium=anime", status: http.StatusOK,
			want:     []string{">1 result</span>", `href="/search?medium=anime&amp;q=steins" aria-current="page">anime</a>`},
			dontWant: []string{">franchise</span>", `href="/entries/6/steins-gate-vn">`}},
		// the manga is 31st overall: past the limit unless the filter comes first
		{name: "a medium past the limit can still be reached", path: "/search?q=crowded&medium=manga", status: http.StatusOK,
			want:     []string{">1 result</span>", `href="/entries/600/crowd-manga">Crowd Manga</a>`, `aria-current="page">manga</a>`},
			dontWant: []string{">Crowd 00</a>"}},
		{name: "the count is every hit, and says when it shows fewer", path: "/search?q=crowded", status: http.StatusOK,
			want:     []string{">31 results</span>", "the closest 30 shown", `href="/search?medium=manga&amp;q=crowded">manga</a>`},
			dontWant: []string{"Crowd Manga"}},
		{name: "a medium not among the hits shows them all", path: "/search?q=steins&medium=manga", status: http.StatusOK,
			want: []string{">3 results</span>", `href="/search?q=steins" aria-current="page">all</a>`}},
		{name: "one kind of hit has nothing to filter", path: "/search?q=just+one", status: http.StatusOK,
			want: []string{">1 result</span>"}, dontWant: []string{`aria-label="filter by medium"`}},
		{name: "results that miss point to browsing", path: "/search?q=steins", status: http.StatusOK,
			want:     []string{"not finding it?", `<a href="/entries">browse the catalog</a>`},
			dontWant: []string{"anilist.co/search"}},
		{name: "no hits", path: "/search?q=nothing", status: http.StatusOK,
			want: []string{">nothing in the catalog yet</h2>", "No entry matches &ldquo;nothing&rdquo;.",
				`href="/entries">browse the catalog</a>`},
			dontWant: []string{"or try", "results</span>", "anilist.co/search"}},
		{name: "no hits, but titles nearby", path: "/search?q=toaru+index+s3", status: http.StatusOK,
			want: []string{">or try</p>", `href="/search?q=Toaru+Majutsu+no+Index">Toaru Majutsu no Index</a>`,
				`href="/search?q=A+Certain+Magical+Index">A Certain Magical Index</a>`,
			},
			// every medium is in the catalog, so no one site gets the miss
			dontWant: []string{"anilist.co/search"}},
		{name: "a failed suggestion is a 500", path: "/search?q=near+boom", status: http.StatusInternalServerError},
		{name: "the query is escaped", path: "/search?q=%3Cscript%3E", status: http.StatusOK,
			want: []string{"&lt;script&gt;"}, dontWant: []string{`value="<script>`, "&ldquo;<script>"}},
		{name: "search error is a 500", path: "/search?q=boom", status: http.StatusInternalServerError},
		{name: "browse everything", path: "/entries", status: http.StatusOK,
			want: []string{`href="/entries" aria-current="page"`, ">everything</span>", `href="/entries/7/steins-gate-anime">Steins;Gate</a>`, `href="/entries/6/steins-gate-vn">Steins;Gate</a>`,
				`<link rel="canonical" href="/entries">`}},
		{name: "unfiltered, a shelf per medium", path: "/entries", status: http.StatusOK,
			want: []string{`id="shelf-anime"`, ">anime</h2>", ">visual novel</h2>",
				`<a href="/entries?medium=anime">view all<span class="sr-only"> anime</span></a>`,
				">anime · spring 2011 · 24 episodes</span>", ">visual novel · 2009</span>"},
			// shelves, not the filtered grid
			dontWant: []string{`id="matches"`, "clear filters"}},
		{name: "browse has the one search, up top", path: "/entries", status: http.StatusOK,
			want: []string{`relative h-[52px] w-full`}, dontWant: []string{"max-w-[420px]"}},
		{name: "the filters, each with any", path: "/entries", status: http.StatusOK,
			want: []string{`<label for="filter-medium"`, `<option value="" selected>any</option>`, `<option value="visual_novel">visual novel</option>`,
				`<option value="sci-fi">sci-fi</option>`, `<option value="2009">2009</option>`, `<option value="fall">fall</option>`,
				`<button type="submit"`}},
		{name: "a filter shows what it chose", path: "/entries?genre=sci-fi", status: http.StatusOK,
			want: []string{`<option value="sci-fi" selected>sci-fi</option>`, ">1 entry</h2>", "clear filters",
				`href="/entries/7/steins-gate-anime">Steins;Gate</a>`, `<meta name="robots" content="noindex">`},
			dontWant: []string{"/entries/6/steins-gate-vn", `id="shelf-anime"`}},
		{name: "filters that match nothing say so", path: "/entries?medium=visual_novel&year=2011", status: http.StatusOK,
			want: []string{">0 entries</h2>", "Nothing matches these filters yet. Try loosening one.", `<a href="/entries">clear filters</a>`}},
		{name: "a genre nobody has", path: "/entries?genre=pottery", status: http.StatusNotFound},
		{name: "results come a page at a time", path: "/entries?genre=many", status: http.StatusOK,
			want: []string{">100 entries</h2>", ">Many 000</a>", ">Many 047</a>", "page 1 of 3",
				`<a href="/entries?genre=many&amp;page=2">next`},
			dontWant: []string{">Many 048</a>", "previous"}},
		{name: "a middle page links both ways", path: "/entries?genre=many&page=2", status: http.StatusOK,
			want: []string{">Many 048</a>", "page 2 of 3", `<a href="/entries?genre=many">`, `<a href="/entries?genre=many&amp;page=3">next`,
				"<title>browse, page 2 | listo</title>"},
			dontWant: []string{">Many 047</a>", ">Many 096</a>"}},
		{name: "the last page has no next", path: "/entries?genre=many&page=3", status: http.StatusOK,
			want: []string{">Many 099</a>", "page 3 of 3", "previous"}, dontWant: []string{"next"}},
		{name: "one page needs no pager", path: "/entries?genre=sci-fi", status: http.StatusOK,
			dontWant: []string{`aria-label="pages"`}},
		{name: "a page past the end", path: "/entries?genre=many&page=4", status: http.StatusNotFound},
		{name: "a page that isn't a number", path: "/entries?genre=many&page=two", status: http.StatusNotFound},
		// (page-1)*48 past int32: these wrapped to offset 32 and to a negative
		{name: "a page so far it would wrap the offset", path: "/entries?genre=many&page=89478487", status: http.StatusNotFound},
		{name: "a page so far it would go negative", path: "/entries?genre=many&page=44739244", status: http.StatusNotFound},
		{name: "the shelves are one page", path: "/entries?page=2", status: http.StatusNotFound},
		{name: "cards by default, and the switch says so", path: "/entries", status: http.StatusOK,
			want: []string{`aria-label="view"`, `<a href="/entries?view=cards" aria-current="page">cards</a>`,
				`<a href="/entries?view=list">list</a>`, `<a href="/entries?view=text">text only</a>`, "aspect-[172/250]"}},
		{name: "a list keeps small covers", path: "/entries?genre=sci-fi&view=list", status: http.StatusOK,
			want:     []string{"h-[46px] w-[33px]", `<input type="hidden" name="view" value="list">`, `aria-current="page">list</a>`},
			dontWant: []string{"aspect-[172/250]"}},
		{name: "text only has no covers at all", path: "/entries?view=text", status: http.StatusOK,
			want:     []string{`aria-current="page">text only</a>`, `<a href="/entries?medium=anime&amp;view=text">view all`},
			dontWant: []string{"aspect-[172/250]", "h-[46px] w-[33px]"}},
		{name: "a view that isn't one", path: "/entries?view=gallery", status: http.StatusNotFound},
		{name: "a year nobody started in", path: "/entries?year=1890", status: http.StatusNotFound},
		{name: "a season that isn't one", path: "/entries?season=monsoon", status: http.StatusNotFound},
		{name: "browse one medium", path: "/entries?medium=anime", status: http.StatusOK,
			want: []string{`href="/entries/7/steins-gate-anime">Steins;Gate</a>`, `<option value="anime" selected>anime</option>`,
				`<link rel="canonical" href="/entries?medium=anime">`, `href="/entries?medium=anime" aria-current="page"`},
			dontWant: []string{"/entries/6/steins-gate-vn", "noindex"}},
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
	sgRelease   = uuid.MustParse("01a151f2-0000-7000-8000-00000000f002")
	monoSeries  = uuid.MustParse("01a151f2-0000-7000-8000-000000000103")
	bare        = uuid.MustParse("01a151f2-0000-7000-8000-000000000009")
	airing      = uuid.MustParse("01a151f2-0000-7000-8000-000000000010")
	monoOrder   = uuid.MustParse("01a151f2-0000-7000-8000-00000000f103")
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
	case 10:
		// airing: what's out and what's next come from its units
		start, month, releasing, desc, source := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), "month", "releasing", airingDescription, "light_novel"
		mins := int32(24)
		return sqlcgen.EntryByPublicIDRow{ID: airing, PublicID: 10, Slug: "airing", Medium: "anime", MediumLabel: "anime",
			Format: "tv", FormatLabel: "anime", Title: "Airing", Description: &desc, Status: &releasing, StartDate: &start,
			StartPrecision: &month, UnitKind: "episode", UnitTotal: &total, UnitMinutes: &mins, SourceType: &source}, nil
	case 9:
		// nothing but what every entry has: no date, units, titles, tags,
		// relations, groups or links
		return sqlcgen.EntryByPublicIDRow{ID: bare, PublicID: 9, Slug: "bare", Medium: "anime", MediumLabel: "anime",
			Format: "tv", FormatLabel: "anime", Title: "Bare", UnitKind: "none"}, nil
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
	case 103:
		return sqlcgen.GroupByPublicIDRow{ID: monoSeries, PublicID: 103, Slug: "monogatari", Kind: "series", Name: "Monogatari"}, nil
	case 500:
		return sqlcgen.GroupByPublicIDRow{}, errors.New("connection reset")
	}
	return sqlcgen.GroupByPublicIDRow{}, pgx.ErrNoRows
}

func (fakeCatalog) GroupOrderings(_ context.Context, id uuid.UUID) ([]sqlcgen.GroupOrderingsRow, error) {
	if id == monoSeries {
		return monogatariOrder(), nil
	}
	if id == saUniverse {
		// names that slugify alike, and one with no slug at all
		one := func(o, name string) sqlcgen.GroupOrderingsRow {
			return sqlcgen.GroupOrderingsRow{OrderingID: uuid.MustParse(o), Kind: "release", Name: name, Origin: "staff",
				Position: 1, Branch: 1, Step: 1, PublicID: 7, Slug: "steins-gate-anime", Title: "Steins;Gate", Necessity: "required", UnitKind: "none"}
		}
		return []sqlcgen.GroupOrderingsRow{
			one("01a151f2-0000-7000-8000-00000000f201", "Release order"),
			one("01a151f2-0000-7000-8000-00000000f202", "Release-order"),
			one("01a151f2-0000-7000-8000-00000000f203", "初見順"),
		}, nil
	}
	if id != sgFranchise {
		return nil, nil
	}
	ep1, ep22, beta := "Ep. 1", "Ep. 22", "23β"
	rationale, note := "Watch the series, then the branch the sequel follows.", "epilogue"
	step := func(pos, branch int32, publicID int64, slug, title string) sqlcgen.GroupOrderingsRow {
		return sqlcgen.GroupOrderingsRow{OrderingID: sgOrder, Kind: "recommended", Name: "first-timer order",
			Rationale: &rationale, Origin: "staff", Position: pos, Branch: branch, Step: 1, PublicID: publicID, Slug: slug,
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
	release := step(1, 1, 7, "steins-gate-anime", "Steins;Gate")
	release.OrderingID, release.Kind, release.Name, release.Rationale, release.Origin = sgRelease, "release", "release order", nil, "auto"
	return []sqlcgen.GroupOrderingsRow{first, second, film, sequel, release}, nil
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
	if id == sgFranchise {
		return []sqlcgen.NodeGroupsRow{{PublicID: 1, Slug: "science-adventure", Kind: "universe", Name: "Science Adventure", Depth: 1}}, nil
	}
	if id != sgAnime {
		return nil, nil
	}
	first := "#1"
	return []sqlcgen.NodeGroupsRow{
		{PublicID: 3, Slug: "steins-gate-main", Kind: "continuity", Name: "Main timeline", PositionLabel: &first, Depth: 1},
		{PublicID: 2, Slug: "steins-gate", Kind: "franchise", Name: "Steins;Gate", Depth: 2},
	}, nil
}

// SearchNodes narrows to a medium, then takes the limit, as the SQL does.
func (fakeCatalog) SearchNodes(_ context.Context, arg sqlcgen.SearchNodesParams) ([]sqlcgen.SearchNodesRow, error) {
	rows, err := searchRows(arg)
	if err != nil {
		return nil, err
	}
	var out []sqlcgen.SearchNodesRow
	for _, r := range rows {
		if arg.Medium != nil && (r.Medium == nil || *r.Medium != *arg.Medium) {
			continue
		}
		if len(out) == int(arg.MaxResults) {
			break
		}
		out = append(out, r)
	}
	return out, nil
}

// SearchMedia counts every row the fake search has, by medium.
func (fakeCatalog) SearchMedia(_ context.Context, q string) ([]sqlcgen.SearchMediaRow, error) {
	rows, err := searchRows(sqlcgen.SearchNodesParams{Q: q, MaxResults: 1000})
	if err != nil {
		return nil, err
	}
	var out []sqlcgen.SearchMediaRow
	at := map[string]int{}
	for _, r := range rows {
		key := deref(r.Medium)
		i, ok := at[key]
		if !ok {
			i = len(out)
			at[key] = i
			out = append(out, sqlcgen.SearchMediaRow{Medium: r.Medium, MediumLabel: r.MediumLabel})
		}
		out[i].Hits++
	}
	return out, nil
}

// searchRows is every hit the fake search has for a query, best first.
func searchRows(arg sqlcgen.SearchNodesParams) ([]sqlcgen.SearchNodesRow, error) {
	franchise, anime, tv, vn, full, month, episode := "franchise", "anime", "tv", "visual_novel", "full", "month", "episode"
	vnLabel, total := "visual novel", int32(24)
	start := time.Date(2011, 4, 6, 0, 0, 0, 0, time.UTC)
	sgAnimeHit := sqlcgen.SearchNodesRow{PublicID: 7, Slug: "steins-gate-anime", Type: "entry", Title: "Steins;Gate",
		Matched: "Steins;Gate", MatchedKind: "english", Medium: &anime, MediumLabel: &anime, Format: &tv, FormatLabel: &anime,
		StartDate: &start, StartPrecision: &month, UnitKind: &episode, UnitTotal: &total, Score: 0.976}
	switch arg.Q {
	case "steins":
		return []sqlcgen.SearchNodesRow{
			{PublicID: 2, Slug: "steins-gate", Type: "group", Title: "Steins;Gate", Matched: "Steins;Gate",
				MatchedKind: "english", GroupKind: &franchise, Score: 1},
			sgAnimeHit,
			{PublicID: 6, Slug: "steins-gate-vn", Type: "entry", Title: "Steins;Gate", Matched: "Steins;Gate",
				MatchedKind: "english", Medium: &vn, MediumLabel: &vnLabel, Format: &full, FormatLabel: &vnLabel, Score: 0.9},
		}, nil
	case "sg":
		alias := sgAnimeHit
		alias.Matched, alias.MatchedKind = "S;G", "synonym"
		return []sqlcgen.SearchNodesRow{alias}, nil
	case "just one":
		return []sqlcgen.SearchNodesRow{sgAnimeHit}, nil
	case "crowded":
		// thirty anime fill the limit; the manga sorts after them all
		manga := "manga"
		var rows []sqlcgen.SearchNodesRow
		for i := range 30 {
			rows = append(rows, sqlcgen.SearchNodesRow{PublicID: int64(500 + i), Slug: "crowd", Type: "entry", Title: fmt.Sprintf("Crowd %02d", i),
				Matched: "Crowd", Medium: &anime, MediumLabel: &anime, FormatLabel: &anime})
		}
		return append(rows, sqlcgen.SearchNodesRow{PublicID: 600, Slug: "crowd-manga", Type: "entry", Title: "Crowd Manga",
			Matched: "Crowd Manga", Medium: &manga, MediumLabel: &manga, FormatLabel: &manga}), nil
	case "boom":
		return nil, errors.New("connection reset")
	case "many":
		// as many hits as asked for, to check the limit
		hits := make([]sqlcgen.SearchNodesRow, arg.MaxResults)
		for i := range hits {
			hits[i] = sqlcgen.SearchNodesRow{PublicID: int64(100 + i), Slug: "x", Type: "entry", Title: "Many", Matched: "Many", MediumLabel: &anime}
		}
		return hits, nil
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
		{Medium: "book", Label: "book", Entries: 0},
	}, nil
}

// browseRows is the fake catalog as browse sees it: the anime dated to the
// month and tagged sci-fi, the VN only to the year.
func browseRows() []sqlcgen.BrowseShelvesRow {
	month, year, total := "month", "year", int32(24)
	anime := time.Date(2011, 4, 1, 0, 0, 0, 0, time.UTC)
	vn := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	return []sqlcgen.BrowseShelvesRow{
		{PublicID: 7, Slug: "steins-gate-anime", Title: "Steins;Gate", Medium: "anime", MediumLabel: "anime", Format: "tv",
			FormatLabel: "anime", StartDate: &anime, StartPrecision: &month, UnitKind: "episode", UnitTotal: &total},
		{PublicID: 6, Slug: "steins-gate-vn", Title: "Steins;Gate", Medium: "visual_novel", MediumLabel: "visual novel",
			Format: "full", FormatLabel: "visual novel", StartDate: &vn, StartPrecision: &year, UnitKind: "none"},
	}
}

// BrowseEntries filters by medium, year and the sci-fi genre, and pages like
// limit and offset; the SQL itself is the data package's to test. The "many"
// genre is 100 entries, to have pages to turn.
func (fakeCatalog) BrowseEntries(_ context.Context, arg sqlcgen.BrowseEntriesParams) ([]sqlcgen.BrowseEntriesRow, error) {
	var matches []sqlcgen.BrowseShelvesRow
	if arg.Genre != nil && *arg.Genre == "many" {
		for i := range 100 {
			matches = append(matches, sqlcgen.BrowseShelvesRow{PublicID: int64(1000 + i), Slug: "many", Title: fmt.Sprintf("Many %03d", i),
				Medium: "anime", MediumLabel: "anime", Format: "tv", FormatLabel: "anime", UnitKind: "none"})
		}
	}
	for _, r := range browseRows() {
		if arg.Genre != nil && (*arg.Genre != "sci-fi" || r.Medium != "anime") {
			continue
		}
		if arg.Medium != nil && *arg.Medium != r.Medium {
			continue
		}
		if arg.Year != nil && int(*arg.Year) != r.StartDate.Year() {
			continue
		}
		matches = append(matches, r)
	}
	var out []sqlcgen.BrowseEntriesRow
	for i, r := range matches {
		if i < int(arg.Skip) || i >= int(arg.Skip+arg.PageSize) {
			continue
		}
		out = append(out, sqlcgen.BrowseEntriesRow{PublicID: r.PublicID, Slug: r.Slug, Title: r.Title, Medium: r.Medium,
			MediumLabel: r.MediumLabel, Format: r.Format, FormatLabel: r.FormatLabel, StartDate: r.StartDate,
			StartPrecision: r.StartPrecision, UnitKind: r.UnitKind, UnitTotal: r.UnitTotal, Total: int32(len(matches))})
	}
	return out, nil
}

func (fakeCatalog) BrowseShelves(context.Context, int32) ([]sqlcgen.BrowseShelvesRow, error) {
	return browseRows(), nil
}

func (fakeCatalog) BrowseGenres(context.Context) ([]string, error) {
	return []string{"many", "sci-fi"}, nil
}

func (fakeCatalog) BrowseYears(context.Context) ([]int32, error) { return []int32{2011, 2009}, nil }

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

// A page with nothing in it says so or leaves the part out: no empty list for
// a screen reader to announce, and no controls for content that isn't there.
// The search suggestions list is empty until the script fills it.
func TestEmptyPagesLeaveNoHusks(t *testing.T) {
	srv := httptest.NewServer(newHTTPServer("", &web.Bundle{Files: fstest.MapFS{}}, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(srv.Close)
	emptyList := regexp.MustCompile(`<(ul|ol|dl)([^>]*)>\s*</(ul|ol|dl)>`)
	cases := []struct {
		path     string
		dontWant []string
	}{
		{"/entries/9/bare", nil},
		{"/groups/1/science-adventure", []string{`aria-label="group entries"`}},
		{"/entries?medium=visual_novel&year=2011", nil},
		{"/search?q=nothing", nil},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			res, err := http.Get(srv.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := string(raw)
			if res.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", res.StatusCode)
			}
			for _, m := range emptyList.FindAllStringSubmatch(body, -1) {
				if !strings.Contains(m[2], `id="site-search-results"`) {
					t.Errorf("empty list: %s", m[0])
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

// A view picked once sticks, from a cookie, so browse pages can't be shared
// out of a cache.
func TestBrowseRemembersTheView(t *testing.T) {
	srv := httptest.NewServer(newHTTPServer("", &web.Bundle{Files: fstest.MapFS{}}, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(srv.Close)
	get := func(path string, cookie *http.Cookie) (*http.Response, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, srv.URL+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		return res, plainLinks(string(body))
	}

	res, _ := get("/entries?view=text", nil)
	var picked *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "listo_view" {
			picked = c
		}
	}
	if picked == nil || picked.Value != "text" || picked.Path != "/entries" || !picked.HttpOnly {
		t.Fatalf("cookie = %+v, want listo_view=text on /entries, http only", picked)
	}

	res, body := get("/entries", picked)
	if !strings.Contains(body, `aria-current="page">text only</a>`) || strings.Contains(body, "aspect-[172/250]") {
		t.Errorf("a plain visit didn't remember text only")
	}
	if cc, vary := res.Header.Get("Cache-Control"), res.Header.Get("Vary"); !strings.HasPrefix(cc, "private") || !strings.Contains(vary, "Cookie") {
		t.Errorf("Cache-Control = %q, Vary = %q; want private, varying by cookie", cc, vary)
	}

	if _, body := get("/entries", &http.Cookie{Name: "listo_view", Value: "gallery"}); !strings.Contains(body, `aria-current="page">cards</a>`) {
		t.Errorf("a cookie for no view should fall back to cards")
	}
	// other pages stay shareable
	if res, _ := get("/groups/2/steins-gate", picked); res.Header.Get("Cache-Control") != "public, max-age=60" {
		t.Errorf("group page Cache-Control = %q, want public", res.Header.Get("Cache-Control"))
	}
}

// A landmark inside main reads as a page region it isn't, so the sidebar is
// the only aside on any page.
func TestOneAsidePerPage(t *testing.T) {
	srv := httptest.NewServer(newHTTPServer("", &web.Bundle{Files: fstest.MapFS{}}, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(srv.Close)
	for _, path := range []string{"/", "/search?q=steins", "/search?q=nothing", "/entries", "/entries/7/steins-gate-anime",
		"/groups/2/steins-gate", "/groups/3/steins-gate-main", "/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f001/first-timer-order",
		"/groups/2/steins-gate/orderings/01a151f2-0000-7000-8000-00000000f002/release-order"} {
		t.Run(path, func(t *testing.T) {
			res, err := http.Get(srv.URL + path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if n := strings.Count(string(body), "<aside"); n > 1 {
				t.Errorf("%d asides, want at most 1", n)
			}
		})
	}
}

func TestSearchAPI(t *testing.T) {
	srv := httptest.NewServer(newHTTPServer("", &web.Bundle{Files: fstest.MapFS{}}, fakeDB{}, fakeCatalog{}).srv.Handler)
	t.Cleanup(srv.Close)

	type hit struct {
		Title, Href, Meta, Matched string
	}
	type response struct {
		Query   string
		Results []hit
		All     string
	}

	cases := []struct {
		name     string
		path     string
		status   int
		hits     int
		first    hit
		all      string
		resultsJ string // how the results array is encoded, when it matters
	}{
		{name: "hits shaped like the search page", path: "/api/search?q=steins", status: http.StatusOK, hits: 3,
			first: hit{Title: "Steins;Gate", Href: "/groups/2/steins-gate", Meta: "franchise"}, all: "/search?q=steins"},
		{name: "an alias match says so", path: "/api/search?q=sg", status: http.StatusOK, hits: 1,
			first: hit{Title: "Steins;Gate", Href: "/entries/7/steins-gate-anime", Meta: "anime · 2011", Matched: "S;G"}, all: "/search?q=sg"},
		{name: "a handful by default", path: "/api/search?q=many", status: http.StatusOK, hits: 6, all: "/search?q=many"},
		{name: "a smaller limit", path: "/api/search?q=many&limit=3", status: http.StatusOK, hits: 3, all: "/search?q=many"},
		{name: "a limit is capped", path: "/api/search?q=many&limit=500", status: http.StatusOK, hits: 10, all: "/search?q=many"},
		{name: "a bad limit is the default", path: "/api/search?q=many&limit=lots", status: http.StatusOK, hits: 6, all: "/search?q=many"},
		{name: "no query, no hits, and an empty list not null", path: "/api/search?q=+", status: http.StatusOK, hits: 0,
			all: "/search?q=", resultsJ: `"results":[]`},
		{name: "the full search link is escaped", path: "/api/search?q=a%26b%3Dc", status: http.StatusOK, hits: 0,
			all: "/search?q=a%26b%3Dc"},
		{name: "a failed search is a 500", path: "/api/search?q=boom", status: http.StatusInternalServerError},
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
			if res.StatusCode != tc.status {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.status)
			}
			if ct, cc := res.Header.Get("Content-Type"), res.Header.Get("Cache-Control"); ct != "application/json" || cc != "no-store" {
				t.Errorf("Content-Type %q, Cache-Control %q", ct, cc)
			}
			if tc.status != http.StatusOK {
				return
			}
			if tc.resultsJ != "" && !strings.Contains(string(raw), tc.resultsJ) {
				t.Errorf("body %s, want %s", raw, tc.resultsJ)
			}
			var got response
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Results) != tc.hits || got.All != tc.all {
				t.Errorf("got %d hits, all %q; want %d, %q", len(got.Results), got.All, tc.hits, tc.all)
			}
			if tc.first != (hit{}) && (len(got.Results) == 0 || got.Results[0] != tc.first) {
				t.Errorf("first = %+v, want %+v", got.Results, tc.first)
			}
		})
	}
}

func (fakeCatalog) TitleSuggestions(_ context.Context, arg sqlcgen.TitleSuggestionsParams) ([]string, error) {
	switch arg.Q {
	case "toaru index s3":
		return []string{"Toaru Majutsu no Index", "A Certain Magical Index"}, nil
	case "near boom":
		return nil, errors.New("connection reset")
	}
	return nil, nil
}

func (fakeCatalog) GroupEntries(_ context.Context, id uuid.UUID) ([]sqlcgen.GroupEntriesRow, error) {
	if id == sgMain {
		// the anime straight in the continuity, first in it
		anime, month, total, first := "anime", "month", int32(24), "#1"
		start := time.Date(2011, 4, 1, 0, 0, 0, 0, time.UTC)
		return []sqlcgen.GroupEntriesRow{{
			PublicID: 7, Slug: "steins-gate-anime", Title: "Steins;Gate", Medium: "anime", MediumLabel: anime, Format: "tv",
			FormatLabel: anime, StartDate: &start, StartPrecision: &month, UnitKind: "episode", UnitTotal: &total,
			PositionLabel: &first, BranchPublicID: 7, BranchSlug: "steins-gate-anime", BranchType: "entry",
		}}, nil
	}
	if id != sgFranchise {
		return nil, nil
	}
	main, continuity, primary, anime, role := "Main timeline", "continuity", "primary", "anime", "adaptation of the VN"
	month, total := "month", int32(24)
	start := time.Date(2011, 4, 1, 0, 0, 0, 0, time.UTC)
	return []sqlcgen.GroupEntriesRow{{
		PublicID: 7, Slug: "steins-gate-anime", Title: "Steins;Gate", Medium: "anime", MediumLabel: anime, Format: "tv",
		FormatLabel: anime, StartDate: &start, StartPrecision: &month, UnitKind: "episode", UnitTotal: &total, Role: &role,
		BranchPublicID: 3, BranchSlug: "steins-gate-main", BranchType: "group", BranchName: &main, BranchKind: &continuity,
		BranchStanding: &primary,
	}}, nil
}

func (fakeCatalog) GroupSiblings(_ context.Context, id uuid.UUID) ([]sqlcgen.GroupSiblingsRow, error) {
	if id != sgFranchise {
		return nil, nil
	}
	return []sqlcgen.GroupSiblingsRow{{PublicID: 5, Slug: "robotics-notes", Kind: "franchise", Name: "Robotics;Notes", Entries: 2}}, nil
}

// monogatariOrder is six steps, so a five-step preview leaves one out. The
// film has no episodes; the rest add up to 68.
func monogatariOrder() []sqlcgen.GroupOrderingsRow {
	note := "a prequel, but not first"
	step := func(pos int32, publicID int64, slug, title string, eps int32) sqlcgen.GroupOrderingsRow {
		r := sqlcgen.GroupOrderingsRow{OrderingID: monoOrder, Kind: "recommended", Name: "first-timer order", Origin: "staff",
			Position: pos, Branch: 1, Step: 1, PublicID: publicID, Slug: slug, Title: title, Necessity: "required", UnitKind: "none"}
		if eps > 0 {
			r.UnitKind, r.UnitTotal = "episode", &eps
		}
		return r
	}
	kizu := step(3, 112, "kizumonogatari-1", "Kizumonogatari I: Tekketsu", 0)
	kizu.Note = &note
	neko := step(4, 106, "nekomonogatari-black", "Nekomonogatari (Black)", 4)
	neko.Necessity = "recommended"
	return []sqlcgen.GroupOrderingsRow{
		step(1, 104, "bakemonogatari", "Bakemonogatari", 15),
		step(2, 105, "nisemonogatari", "Nisemonogatari", 11),
		kizu,
		neko,
		step(5, 107, "monogatari-second-season", "Monogatari Series: Second Season", 26),
		step(6, 110, "owarimonogatari", "Owarimonogatari", 12),
	}
}

// airingDescription is long enough that the meta description clips it.
const airingDescription = "A show still on the air, with a description long enough to need clipping in a search result. " +
	"It runs on past the point a search result would show, so the page has to cut it at a word and say so."

func (fakeCatalog) EntryCompanies(_ context.Context, id uuid.UUID) ([]sqlcgen.EntryCompaniesRow, error) {
	switch id {
	case sgAnime:
		return []sqlcgen.EntryCompaniesRow{{Name: "White Fox", Role: "studio"}}, nil
	case airing:
		return []sqlcgen.EntryCompaniesRow{{Name: "Studio A", Role: "studio"}, {Name: "Studio B", Role: "studio"},
			{Name: "Pub", Role: "publisher"}}, nil
	}
	return nil, nil
}

// EntryUnits gives the airing show two episodes out, the third due in three
// days and a bit, and the fourth later; air dates hang off now so the page's
// relative times hold.
func (fakeCatalog) EntryUnits(_ context.Context, id uuid.UUID) ([]sqlcgen.EntryUnitsRow, error) {
	if id != airing {
		return nil, nil
	}
	now := time.Now()
	at := func(d time.Duration) *time.Time { t := now.Add(d); return &t }
	title := func(s string) *string { return &s }
	return []sqlcgen.EntryUnitsRow{
		{Ordinal: 1, Kind: "episode", Label: "Ep. 1", Title: title("The Start"), AirDate: at(-8 * 24 * time.Hour)},
		{Ordinal: 2, Kind: "episode", Label: "Ep. 2", Title: title("The Middle"), AirDate: at(-2*24*time.Hour - time.Hour)},
		{Ordinal: 3, Kind: "episode", Label: "Ep. 3", AirDate: at(3*24*time.Hour + 4*time.Hour + 5*time.Minute)},
		{Ordinal: 4, Kind: "episode", Label: "Ep. 4", AirDate: at(10*24*time.Hour + time.Hour)},
	}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
