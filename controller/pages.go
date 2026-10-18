package controller

import (
	"net/http"

	"github.com/khinshankhan/listo/controller/handler"
)

// pageRoutes is every server-rendered page.
func pageRoutes(mux *http.ServeMux, pages *handler.Pages) {
	mux.HandleFunc("GET /{$}", pages.Home)
	mux.HandleFunc("GET /there", pages.There)
}
