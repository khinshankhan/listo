package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/khinshankhan/listo/controller/handler"
	"github.com/khinshankhan/listo/controller/middleware"
	"github.com/khinshankhan/listo/data"
	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/config"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/web"
	"github.com/khinshankhan/logstox/fields"
)

// so a hung request can't hold shutdown forever
const shutdownGrace = 10 * time.Second

// RunHTTPServer serves pages and assets until the context is cancelled.
func RunHTTPServer(ctx context.Context) error {
	bundle, err := web.Load()
	if err != nil {
		return err
	}

	db, err := data.Open(ctx, config.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close()

	return newHTTPServer(config.Addr(), bundle, db, sqlcgen.New(db)).run(ctx)
}

type httpServer struct {
	srv *http.Server
}

func newHTTPServer(addr string, bundle *web.Bundle, db handler.Pinger, catalog handler.Catalog) *httpServer {
	api := handler.NewAPI(db, catalog)
	pages := handler.NewPages(bundle.Tags, catalog)

	// first listed is outermost
	// NOTE: BlackHole wraps the mux directly, it reads the mux's routes
	root := middleware.Chain(
		middleware.Recover,
		middleware.BlackHole(http.HandlerFunc(pages.NotFound), "/api"),
	)(routes(bundle, api, pages))

	return &httpServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: root,
		},
	}
}

// run serves until ctx is cancelled, then shuts down gracefully.
func (s *httpServer) run(ctx context.Context) error {
	listening := make(chan error, 1)
	go func() {
		logx.Default().Named("http").Info("listening", fields.String("addr", s.srv.Addr))
		listening <- s.srv.ListenAndServe()
	}()

	// NOTE: return the listener error, a swallowed "address already in use"
	// leaves the process up with nothing serving
	select {
	case err := <-listening:
		return err

	case <-ctx.Done():
		logx.Default().Named("http").Info("shutting down")

		// ctx is already cancelled
		shutCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
		defer cancel()

		if err := s.srv.Shutdown(shutCtx); err != nil {
			return fmt.Errorf("shutting down: %w", err)
		}

		if err := <-listening; !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

// routes wires up each router; anything unrouted falls to the blackhole.
func routes(bundle *web.Bundle, api *handler.API, pages *handler.Pages) *http.ServeMux {
	mux := http.NewServeMux()

	apiRoutes(mux, api)
	mux.Handle("GET /assets/", handler.Assets(bundle.Files))
	pageRoutes(mux, pages)

	return mux
}
