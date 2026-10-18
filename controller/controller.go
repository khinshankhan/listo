package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/khinshankhan/listo/controller/handler"
	"github.com/khinshankhan/listo/controller/middleware"
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
	return newHTTPServer(config.Addr(), bundle).run(ctx)
}

type httpServer struct {
	srv *http.Server
}

func newHTTPServer(addr string, bundle *web.Bundle) *httpServer {
	// first listed is outermost
	handler := middleware.Chain(
		middleware.Recover,
	)(routes(bundle))

	return &httpServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
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

func routes(bundle *web.Bundle) *http.ServeMux {
	mux := http.NewServeMux()
	pages := handler.NewPages(bundle.Tags)

	mux.Handle("GET /assets/", handler.Assets(bundle.Files))

	mux.HandleFunc("GET /{$}", pages.Home)
	mux.HandleFunc("GET /there", pages.There)

	// anything else is a real 404, not a catch-all shell
	mux.HandleFunc("/", pages.NotFound)

	return mux
}
