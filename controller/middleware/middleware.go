package middleware

import "net/http"

// Middleware wraps a handler with additional behavior.
type Middleware func(http.Handler) http.Handler

// Chain composes mw into one Middleware, first listed outermost, so stacks can
// be named and extended:
//
//	base := Chain(Recover, Logging)
//	protected := Chain(base, RequireAuth)
func Chain(mw ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			next = mw[i](next)
		}
		return next
	}
}
