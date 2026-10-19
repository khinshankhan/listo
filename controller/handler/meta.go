package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/platform"
	"github.com/khinshankhan/listo/platform/logx"
)

// so a hung database can't hang meta with it
const pingTimeout = time.Second

// Meta is the build, database status and server clock.
type Meta struct {
	Module     string    `json:"module"`
	Version    string    `json:"version"`
	Commit     string    `json:"commit,omitempty"`
	Date       string    `json:"date,omitempty"`
	Dirty      bool      `json:"dirty"`
	Database   string    `json:"database"`
	SystemTime time.Time `json:"systemTime"`
}

// Meta always answers 200, even with the database down.
func (a *API) Meta(w http.ResponseWriter, r *http.Request) {
	v := platform.Version("listo")
	writeJSON(w, http.StatusOK, Meta{
		Module:     v.Module,
		Version:    v.Version,
		Commit:     v.Commit,
		Date:       v.Date,
		Dirty:      v.Dirty,
		Database:   a.database(r.Context()),
		SystemTime: time.Now().UTC(),
	})
}

func (a *API) database(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := a.db.Ping(ctx); err != nil {
		logx.Default().Named("meta").Warn("database unreachable", fields.Error(err))
		return "unreachable"
	}
	return "ok"
}
