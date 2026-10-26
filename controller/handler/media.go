package handler

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/platform/logx"
)

var mediaTypes = map[string]string{
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

// Media serves stored images by key. Keys are content hashes, so an image is
// cached forever; originals stay private, only re-encoded variants go out.
func Media(store media.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if strings.Contains(key, "-orig.") {
			http.NotFound(w, r)
			return
		}

		rc, err := store.Open(r.Context(), key)
		if errors.Is(err, media.ErrBadKey) || errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			logx.Default().Named("media").Error("opening image",
				fields.String("key", key),
				fields.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		w.Header().Set("Content-Type", mediaTypes[path.Ext(key)])
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		_, _ = io.Copy(w, rc)
	})
}
