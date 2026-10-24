package middleware

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// capableWriter supports everything and records what reached it.
type capableWriter struct {
	*httptest.ResponseRecorder
	hijacked    bool
	deadlineSet bool
}

func (w *capableWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	return nil, nil, nil
}

func (w *capableWriter) SetWriteDeadline(time.Time) error {
	w.deadlineSet = true
	return nil
}

// bareWriter is only an http.ResponseWriter: no flushing, no hijacking.
type bareWriter struct{ rec *httptest.ResponseRecorder }

func (b bareWriter) Header() http.Header {
	return b.rec.Header()
}

func (b bareWriter) Write(p []byte) (int, error) {
	return b.rec.Write(p)
}

func (b bareWriter) WriteHeader(code int) {
	b.rec.WriteHeader(code)
}

// Wrapped writers must pass capabilities through: supported ones work,
// unsupported ones error.
func TestWrappersPassCapabilitiesThrough(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	wrapped := map[string]Middleware{
		"Recover":   Recover,
		"BlackHole": BlackHole(nil, "/api"),
		"NoStore":   NoStore,
	}

	for name, mw := range wrapped {
		// BlackHole only wraps a *ServeMux
		serve := func(w http.ResponseWriter, inner http.HandlerFunc) {
			var h http.Handler = inner
			if name == "BlackHole" {
				mux := http.NewServeMux()
				mux.Handle("GET /api/x", inner)
				h = mux
			}
			mw(h).ServeHTTP(w, httptest.NewRequest("GET", "/api/x", nil))
		}

		t.Run(name+"/supported reaches the writer", func(t *testing.T) {
			cw := &capableWriter{ResponseRecorder: httptest.NewRecorder()}
			var flushErr, deadlineErr, hijackErr error
			serve(cw, func(w http.ResponseWriter, r *http.Request) {
				rc := http.NewResponseController(w)
				flushErr = rc.Flush()
				deadlineErr = rc.SetWriteDeadline(time.Time{})
				_, _, hijackErr = rc.Hijack()
			})

			if flushErr != nil || !cw.Flushed {
				t.Errorf("Flush: err = %v, reached = %v", flushErr, cw.Flushed)
			}
			if deadlineErr != nil || !cw.deadlineSet {
				t.Errorf("SetWriteDeadline: err = %v, reached = %v", deadlineErr, cw.deadlineSet)
			}
			if hijackErr != nil || !cw.hijacked {
				t.Errorf("Hijack: err = %v, reached = %v", hijackErr, cw.hijacked)
			}
		})

		t.Run(name+"/unsupported is an error", func(t *testing.T) {
			var flushErr, hijackErr error
			serve(bareWriter{httptest.NewRecorder()}, func(w http.ResponseWriter, r *http.Request) {
				rc := http.NewResponseController(w)
				flushErr = rc.Flush()
				_, _, hijackErr = rc.Hijack()
			})

			if !errors.Is(flushErr, http.ErrNotSupported) {
				t.Errorf("Flush: err = %v, want ErrNotSupported", flushErr)
			}
			if !errors.Is(hijackErr, http.ErrNotSupported) {
				t.Errorf("Hijack: err = %v, want ErrNotSupported", hijackErr)
			}
		})
	}
}
