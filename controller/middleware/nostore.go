package middleware

import "net/http"

// NoStore makes every response uncacheable, whatever the handler set, so a
// page in development is never a stale copy.
func NoStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&noStoreWriter{ResponseWriter: w}, r)
	})
}

// noStoreWriter rewrites Cache-Control as the headers go out. Once they're
// sent, the rewrite is a no-op.
type noStoreWriter struct {
	http.ResponseWriter
}

func (w *noStoreWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *noStoreWriter) WriteHeader(code int) {
	w.Header().Set("Cache-Control", "no-store")
	w.ResponseWriter.WriteHeader(code)
}

func (w *noStoreWriter) Write(b []byte) (int, error) {
	w.Header().Set("Cache-Control", "no-store")
	return w.ResponseWriter.Write(b)
}

// FlushError is what http.ResponseController calls; a flush sends the
// headers too.
func (w *noStoreWriter) FlushError() error {
	w.Header().Set("Cache-Control", "no-store")
	return http.NewResponseController(w.ResponseWriter).Flush()
}
