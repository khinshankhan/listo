package handler

import (
	"context"

	"github.com/khinshankhan/listo/data/sqlcgen"
)

// Pinger is the slice of the database pool the API needs so far.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Searcher is the slice of the queries the search API needs.
type Searcher interface {
	SearchNodes(ctx context.Context, arg sqlcgen.SearchNodesParams) ([]sqlcgen.SearchNodesRow, error)
}

// API serves the JSON endpoints under /api.
type API struct {
	db     Pinger
	search Searcher
}

func NewAPI(db Pinger, search Searcher) *API {
	return &API{db: db, search: search}
}
