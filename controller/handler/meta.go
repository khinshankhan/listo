package handler

import (
	"net/http"
	"time"

	"github.com/khinshankhan/listo/platform"
)

// Meta is the build and server clock.
type Meta struct {
	Module     string    `json:"module"`
	Version    string    `json:"version"`
	Commit     string    `json:"commit,omitempty"`
	Date       string    `json:"date,omitempty"`
	Dirty      bool      `json:"dirty"`
	SystemTime time.Time `json:"systemTime"`
}

func MetaHandler(w http.ResponseWriter, r *http.Request) {
	v := platform.Version("listo")
	writeJSON(w, http.StatusOK, Meta{
		Module:     v.Module,
		Version:    v.Version,
		Commit:     v.Commit,
		Date:       v.Date,
		Dirty:      v.Dirty,
		SystemTime: time.Now().UTC(),
	})
}
