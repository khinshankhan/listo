package handler

import (
	"io/fs"
	"net/http"
)

// Assets serves the hashed Vite build; names change with content, so it's
// cached forever.
func Assets(files fs.FS) http.Handler {
	server := http.FileServerFS(files)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		server.ServeHTTP(w, r)
	})
}
