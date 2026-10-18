package middleware

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// fullWriter implements every optional interface, unlike
// httptest.ResponseRecorder.
type fullWriter struct {
	*httptest.ResponseRecorder
	hijacked    bool
	pushed      string
	readFrom    bool
	deadlineSet bool
}

func (w *fullWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	return nil, nil, nil
}

func (w *fullWriter) Push(target string, _ *http.PushOptions) error {
	w.pushed = target
	return nil
}

func (w *fullWriter) ReadFrom(r io.Reader) (int64, error) {
	w.readFrom = true
	return io.Copy(w.ResponseRecorder.Body, r)
}

// SetWriteDeadline is what http.ResponseController looks for after Unwrap.
func (w *fullWriter) SetWriteDeadline(time.Time) error {
	w.deadlineSet = true
	return nil
}

func newFullWriter() *fullWriter {
	return &fullWriter{ResponseRecorder: httptest.NewRecorder()}
}

func TestWrappersPreserveCapabilities(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	wrapped := map[string]func(http.Handler) http.Handler{
		"Recover":   Recover,
		"BlackHole": BlackHole(nil, "/api"),
	}

	for name, mw := range wrapped {
		t.Run(name, func(t *testing.T) {
			var (
				gotFlusher, gotHijacker, gotPusher, gotReaderFrom bool
				ctrlErr                                           error
			)

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, gotFlusher = w.(http.Flusher)
				_, gotHijacker = w.(http.Hijacker)
				_, gotPusher = w.(http.Pusher)
				_, gotReaderFrom = w.(io.ReaderFrom)

				ctrlErr = http.NewResponseController(w).SetWriteDeadline(time.Time{})
			})

			// BlackHole only wraps a *ServeMux
			var h http.Handler = inner
			if name == "BlackHole" {
				m := http.NewServeMux()
				m.Handle("GET /api/x", inner)
				h = m
			}

			fw := newFullWriter()
			mw(h).ServeHTTP(fw, httptest.NewRequest("GET", "/api/x", nil))

			if !gotFlusher || !gotHijacker || !gotPusher || !gotReaderFrom {
				t.Errorf("dropped interfaces: Flusher=%v Hijacker=%v Pusher=%v ReaderFrom=%v",
					gotFlusher, gotHijacker, gotPusher, gotReaderFrom)
			}
			if ctrlErr != nil {
				t.Errorf("ResponseController.SetWriteDeadline = %v, want nil (Unwrap missing?)", ctrlErr)
			}
			if !fw.deadlineSet {
				t.Error("SetWriteDeadline never reached the underlying writer")
			}
		})
	}

}

func TestWrapperForwards(t *testing.T) {
	fw := newFullWriter()
	w := &wrapper{ResponseWriter: fw}

	if _, _, err := w.Hijack(); err != nil || !fw.hijacked {
		t.Errorf("Hijack did not reach the underlying writer (err=%v)", err)
	}
	if err := w.Push("/style.css", nil); err != nil || fw.pushed != "/style.css" {
		t.Errorf("Push did not reach the underlying writer (err=%v, target=%q)", err, fw.pushed)
	}
	if _, err := w.ReadFrom(strings.NewReader("hello")); err != nil || !fw.readFrom {
		t.Errorf("ReadFrom did not reach the underlying writer (err=%v)", err)
	}
	if got := fw.Body.String(); got != "hello" {
		t.Errorf("ReadFrom wrote %q, want %q", got, "hello")
	}
	if w.Unwrap() != http.ResponseWriter(fw) {
		t.Error("Unwrap did not return the underlying writer")
	}
}

// a missing capability is an error, not a panic
func TestWrapperReportsUnsupported(t *testing.T) {
	// bare writer, only Header/Write/WriteHeader
	w := &wrapper{ResponseWriter: bareWriter{httptest.NewRecorder()}}

	if _, _, err := w.Hijack(); err == nil {
		t.Error("Hijack on a non-hijackable writer returned nil error")
	}
	if err := w.Push("/x", nil); err == nil {
		t.Error("Push on a non-pushable writer returned nil error")
	}
	w.Flush() // must not panic

	if n, err := w.ReadFrom(strings.NewReader("abc")); err != nil || n != 3 {
		t.Errorf("ReadFrom fallback = (%d, %v), want (3, nil)", n, err)
	}
}

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

// once a response has started a panic just cuts it off
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
			w.WriteHeader(http.StatusOK)
			http.NewResponseController(w).Flush()
		}},
		{"after ReadFrom", func(w http.ResponseWriter) {
			io.Copy(w, strings.NewReader("partial"))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := newFullWriter()
			Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.start(w)
				panic("boom")
			})).ServeHTTP(fw, httptest.NewRequest("GET", "/", nil))

			if fw.Code != http.StatusOK {
				t.Errorf("status = %d, want 200 (already sent)", fw.Code)
			}
			if strings.Contains(fw.Body.String(), "Internal Server Error") {
				t.Errorf("appended an error to a started response: %q", fw.Body.String())
			}
		})
	}
}
