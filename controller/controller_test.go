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

	"github.com/khinshankhan/listo/web"
)

func TestRoutes(t *testing.T) {
	bundle := &web.Bundle{
		Files: fstest.MapFS{"assets/app-abc.css": {Data: []byte("body{}")}},
		Tags:  web.Tags{CSS: []string{"/assets/app-abc.css"}, Script: "/assets/islands-abc.js"},
	}
	up := httptest.NewServer(newHTTPServer("", bundle, fakeDB{}).srv.Handler)
	t.Cleanup(up.Close)
	down := httptest.NewServer(newHTTPServer("", bundle, fakeDB{err: errors.New("connection refused")}).srv.Handler)
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
