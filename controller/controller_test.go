package controller

import (
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
	srv := httptest.NewServer(routes(bundle))
	t.Cleanup(srv.Close)

	cases := []struct {
		name   string
		path   string
		status int
		want   []string
	}{
		{"home", "/", http.StatusOK, []string{"hello world", `data-island="wave"`, `{&#34;name&#34;:&#34;world&#34;}`, "/assets/islands-abc.js"}},
		{"there", "/there", http.StatusOK, []string{"hello there", `data-island="wave"`}},
		{"unknown path is a real 404", "/nope", http.StatusNotFound, []string{"back home"}},
		{"there has no subpaths", "/there/x", http.StatusNotFound, []string{"back home"}},
		{"hashed asset", "/assets/app-abc.css", http.StatusOK, []string{"body{}"}},
		{"missing asset", "/assets/gone.js", http.StatusNotFound, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(srv.URL + tc.path)
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
