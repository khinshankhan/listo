package controller

import (
	"net/http"

	"github.com/khinshankhan/listo/controller/handler"
)

// apiRoutes is the JSON API, everything under /api.
func apiRoutes(mux *http.ServeMux, api *handler.API) {
	mux.HandleFunc("GET /api/meta", api.Meta)
	mux.HandleFunc("GET /api/search", api.Search)
}
