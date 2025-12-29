package logger

import (
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

// middleware
func MyMiddleware(log *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("logger middleware enabled")
		entry := log.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)

		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 4096)
				stack := string(buf[:runtime.Stack(buf, false)])

				entry.Error("panic!!!",
			slog.Any("panic", r),
			slog.String("stack", stack))

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
			}
		}()

		t1 := time.Now()
		defer func() {
			entry.Info("request complited",
				slog.String("duration", time.Since(t1).String()))
		}()

		next(w, r)
	}
}
