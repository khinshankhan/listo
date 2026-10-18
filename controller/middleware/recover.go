package middleware

import (
	"io"
	"net/http"
	"runtime/debug"

	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/logstox/fields"
)

// Recover turns a panic into a 500, unless the response already started, in
// which case it's cut off.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &recorder{wrapper: wrapper{ResponseWriter: w}}

		defer func() {
			if err := recover(); err != nil {
				// an intentional abort, pass it on
				if err == http.ErrAbortHandler {
					panic(err)
				}

				logx.Default().Named("recover").Error("panic serving request",
					fields.String("method", r.Method),
					fields.String("path", r.URL.Path),
					fields.Any("panic", err),
					fields.String("stack", string(debug.Stack())))

				if !rw.wroteHeader {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

// recorder tracks whether the status line has been sent.
type recorder struct {
	wrapper
	wroteHeader bool
}

func (w *recorder) WriteHeader(code int) {
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *recorder) Write(b []byte) (int, error) {
	w.wroteHeader = true
	return w.ResponseWriter.Write(b)
}

// ReadFrom and Flush bypass Write, so both mark the response as started.
func (w *recorder) ReadFrom(r io.Reader) (int64, error) {
	w.wroteHeader = true
	return w.wrapper.ReadFrom(r)
}

func (w *recorder) Flush() {
	w.wroteHeader = true
	w.wrapper.Flush()
}
