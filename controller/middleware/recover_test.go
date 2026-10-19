package middleware

import (
	"errors"
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

// once a response has started, a panic has to abort it
func TestRecoverAfterResponseStarted(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	tests := []struct {
		name  string
		start func(w http.ResponseWriter)
	}{
		{"after Write", func(w http.ResponseWriter) {
			w.Write([]byte("partial"))
		}},
		{"after Flush", func(w http.ResponseWriter) {
			w.Write([]byte("partial"))
			http.NewResponseController(w).Flush()
		}},
		// flushing alone starts the response
		{"after Flush alone", func(w http.ResponseWriter) {
			http.NewResponseController(w).Flush()
		}},
		{"after ReadFrom", func(w http.ResponseWriter) {
			io.Copy(w, strings.NewReader("partial"))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.start(w)
				panic("boom")
			})))
			t.Cleanup(srv.Close)

			res, err := http.Get(srv.URL)
			if err != nil {
				return // cut off before the status line: aborted
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err == nil {
				t.Errorf("got a complete %d %q, want the response cut off", res.StatusCode, body)
			}
			if strings.Contains(string(body), "Internal Server Error") {
				t.Errorf("appended an error to a started response: %q", body)
			}
		})
	}
}

// an unsupported flush sends nothing, so a panic still gets a 500
func TestRecoverAfterUnsupportedFlush(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	rec := httptest.NewRecorder()
	var flushErr error
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("panic escaped Recover: %v", err)
			}
		}()
		Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flushErr = http.NewResponseController(w).Flush()
			panic("boom")
		})).ServeHTTP(bareWriter{rec}, httptest.NewRequest("GET", "/", nil))
	}()

	if !errors.Is(flushErr, http.ErrNotSupported) {
		t.Fatalf("Flush: err = %v, want ErrNotSupported", flushErr)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}
