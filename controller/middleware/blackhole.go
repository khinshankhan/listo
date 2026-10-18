package middleware

import (
	"net/http"
	"strings"

	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/logstox/fields"
)

// BlackHole replaces the mux's default 404. Under the given path prefixes the
// reply is terse plain text; anywhere else it's fallback, when there is one.
// "/api" covers "/api", "/api/" and "/api/x", not "/apifoo".
//
// NOTE: a catchall route would match any method, so the mux would never 405.
func BlackHole(fallback http.Handler, prefixes ...string) Middleware {
	cleaned := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		if p = strings.TrimSuffix(p, "/"); p != "" {
			cleaned = append(cleaned, p)
		}
	}

	return func(next http.Handler) http.Handler {
		// only a ServeMux can tell an unrouted path from a handler's own 404
		mux, ok := next.(*http.ServeMux)
		if !ok {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inHole := covered(cleaned, r.URL.Path)
			if !inHole && fallback == nil {
				next.ServeHTTP(w, r)
				return
			}

			// no pattern also covers a method mismatch, dispatch picks 404 or 405
			if _, pattern := mux.Handler(r); pattern != "" {
				next.ServeHTTP(w, r)
				return
			}

			bw := &blackHoleWriter{wrapper: wrapper{ResponseWriter: w}}
			next.ServeHTTP(bw, r)

			if !bw.swallowed {
				return
			}
			if !inHole {
				fallback.ServeHTTP(w, r)
				return
			}

			// structured, a path with newlines could forge log lines
			logx.Default().Named("blackhole").Info("no route",
				fields.String("method", r.Method),
				fields.String("path", r.URL.Path))
			http.Error(w, "ERROR 404: DATA NOT FOUND IN THIS SECTOR.", http.StatusNotFound)
		})
	}
}

func covered(prefixes []string, path string) bool {
	for _, p := range prefixes {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// blackHoleWriter swallows the mux's 404 so a replacement can be written.
type blackHoleWriter struct {
	wrapper
	swallowed   bool
	wroteHeader bool
}

func (w *blackHoleWriter) WriteHeader(code int) {
	w.wroteHeader = true

	// a 405 carries Allow, a 404 never does
	if code == http.StatusNotFound && w.Header().Get("Allow") == "" {
		w.swallowed = true
		return
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *blackHoleWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.swallowed {
		// discard "404 page not found", report success to avoid a short write
		return len(b), nil
	}
	return w.ResponseWriter.Write(b)
}
