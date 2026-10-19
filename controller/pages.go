package controller

import (
	"net/http"

	"github.com/khinshankhan/listo/controller/handler"
)

// pageRoutes is every server-rendered page.
func pageRoutes(mux *http.ServeMux, pages *handler.Pages) {
	mux.HandleFunc("GET /{$}", pages.Home)
	mux.HandleFunc("GET /there", pages.There)

	// the id is the key, the slug a redirect away
	mux.HandleFunc("GET /entries/{id}", pages.Entry)
	mux.HandleFunc("GET /entries/{id}/{slug}", pages.Entry)
}
