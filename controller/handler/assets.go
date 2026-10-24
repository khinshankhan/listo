package handler

import (
	"io/fs"
	"net/http"
	"strings"
)

// Assets serves the hashed Vite build; names change with content, so it's
// cached forever.
func Assets(files fs.FS) http.Handler {
	server := http.FileServerFS(files)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// a directory listing changes every build, so it can't be immutable
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		server.ServeHTTP(w, r)
	})
}

// favicon is the wordmark's "l." as shapes, since a favicon can't load our
// fonts: a dark stem and the accent full stop on the cream ground.
const favicon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">` +
	`<rect width="32" height="32" rx="7" fill="#f2efe4"/>` +
	`<rect x="9" y="5" width="5.5" height="22" rx="2.75" fill="#1c1d19"/>` +
	`<circle cx="21" cy="23.75" r="3.25" fill="#3f6b3a"/>` +
	`</svg>`

// Favicon serves the icon. It changes about never, so a day of caching is
// plenty.
func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = w.Write([]byte(favicon))
}
