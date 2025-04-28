package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// PanicRecovery is a middleware that recovers from panics while logging the error and returning an internal server error.
// It should be the first middleware in the chain.
// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-5-middleware/
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				slog.Error("Server panic:",
					slog.String("Request method:", r.Method),
					slog.String("URI:", r.RequestURI),
					slog.String("Remote address:", r.RemoteAddr),
					slog.Any("Headers:", r.Header),
				)
				slog.Debug("HTTP handler panic:", slog.Any("debug", debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler, accessLogger *slog.Logger) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLogger.Debug("",
			slog.String("Request method:", r.Method),
			slog.String("URI:", r.RequestURI),
			slog.String("Remote address:", r.RemoteAddr),
			slog.Any("Headers:", r.Header),
		)
		next.ServeHTTP(w, r)
	})
}
