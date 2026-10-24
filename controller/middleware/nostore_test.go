package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoStore(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"header set, then WriteHeader", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=60")
			w.WriteHeader(http.StatusOK)
		}},
		{"header set, then only Write", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Write([]byte("asset"))
		}},
		{"header set, then a flush", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "private, max-age=60")
			http.NewResponseController(w).Flush()
		}},
		{"no header at all", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "nope", http.StatusNotFound)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			NoStore(tt.handler).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

			if got := rec.Result().Header.Values("Cache-Control"); len(got) != 1 || got[0] != "no-store" {
				t.Errorf("Cache-Control = %q, want [no-store]", got)
			}
		})
	}
}
