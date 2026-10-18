package handler

import (
	"encoding/json"
	"net/http"

	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/platform/logx"
)

// writeJSON writes an uncached API response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logx.Default().Named("json").Error("writing response", fields.Error(err))
	}
}
