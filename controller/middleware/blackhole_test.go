package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testMux() *http.ServeMux {
	mux := http.NewServeMux()
	ok := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	mux.HandleFunc("GET /api/healthz", ok)
	// a real route that starts with the prefix
	mux.HandleFunc("GET /apidocs", ok)
	return mux
}

func TestBlackHole(t *testing.T) {
	h := BlackHole(nil, "/api", "/images")(testMux())

	const sector = "ERROR 404: DATA NOT FOUND IN THIS SECTOR."

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name: "registered route is untouched", method: "GET", path: "/api/healthz",
			wantStatus: 200, wantBody: "ok",
		},
		{
			name: "method mismatch still gets 405", method: "POST", path: "/api/healthz",
			wantStatus: 405, wantAllow: "GET, HEAD",
		},
		{
			name: "unknown path under prefix", method: "GET", path: "/api/nope",
			wantStatus: 404, wantBody: sector,
		},
		{
			name: "prefix root with slash", method: "GET", path: "/api/",
			wantStatus: 404, wantBody: sector,
		},
		{
			name: "bare prefix", method: "GET", path: "/api",
			wantStatus: 404, wantBody: sector,
		},
		{
			name: "second prefix", method: "GET", path: "/images/x.png",
			wantStatus: 404, wantBody: sector,
		},
		{
			name: "route sharing the prefix string still resolves", method: "GET", path: "/apidocs",
			wantStatus: 200, wantBody: "ok",
		},
		{
			name: "prefix is a path segment, not a substring", method: "GET", path: "/apifoo",
			wantStatus: 404, wantBody: "404 page not found",
		},
		{
			name: "outside every prefix", method: "GET", path: "/other",
			wantStatus: 404, wantBody: "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(rec.Body.String()); tt.wantBody != "" && got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
			if got := rec.Header().Get("Allow"); got != tt.wantAllow {
				t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
			}
		})
	}
}

// a handler's own 404 passes through
func TestBlackHoleKeepsHandler404(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/domains/{domain}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"no such domain"}`, http.StatusNotFound)
	})

	rec := httptest.NewRecorder()
	BlackHole(nil, "/api")(mux).ServeHTTP(rec, httptest.NewRequest("GET", "/api/domains/nope.test", nil))

	if got := strings.TrimSpace(rec.Body.String()); got != `{"error":"no such domain"}` {
		t.Errorf("handler's own 404 was replaced: body = %q", got)
	}
}

func TestBlackHoleFallback(t *testing.T) {
	page := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("lost page"))
	})
	h := BlackHole(page, "/api")(testMux())

	const sector = "ERROR 404: DATA NOT FOUND IN THIS SECTOR."

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "unrouted page gets the fallback", method: "GET", path: "/nope", wantStatus: 404, wantBody: "lost page"},
		{name: "prefix substring is a page, not the hole", method: "GET", path: "/apifoo", wantStatus: 404, wantBody: "lost page"},
		{name: "the hole still wins under its prefix", method: "GET", path: "/api/nope", wantStatus: 404, wantBody: sector},
		{name: "routed page is untouched", method: "GET", path: "/apidocs", wantStatus: 200, wantBody: "ok"},
		{name: "method mismatch still gets 405", method: "POST", path: "/apidocs", wantStatus: 405},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(rec.Body.String()); tt.wantBody != "" && got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
