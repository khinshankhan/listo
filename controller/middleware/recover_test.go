package middleware

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRecover(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	t.Run("panic becomes 500", func(t *testing.T) {
		h := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("boom")
		}))

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", rec.Code)
		}
	})

	t.Run("no panic passes through", func(t *testing.T) {
		h := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("fine"))
		}))

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

		if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "fine" {
			t.Errorf("got %d %q, want 200 \"fine\"", rec.Code, rec.Body.String())
		}
	})

	t.Run("ErrAbortHandler is re-panicked", func(t *testing.T) {
		defer func() {
			if err := recover(); err != http.ErrAbortHandler {
				t.Errorf("recovered %v, want ErrAbortHandler to propagate", err)
			}
		}()

		h := Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(http.ErrAbortHandler)
		}))
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	})
}
