package middleware

import (
	"net/http"
	"time"

	"api/logger"

	"github.com/go-chi/chi/v5/middleware"
)

func LogError(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		if ww.Status() >= 400 {
			logger.Logger().Error().
				Int("status", ww.Status()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Dur("duration", duration).
				Msg("HTTP error")
		}
	})
}
