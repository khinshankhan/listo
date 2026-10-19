package handler

import "context"

// Pinger is the slice of the database pool the API needs so far.
type Pinger interface {
	Ping(ctx context.Context) error
}

// API serves the JSON endpoints under /api.
type API struct {
	db Pinger
}

func NewAPI(db Pinger) *API {
	return &API{db: db}
}
