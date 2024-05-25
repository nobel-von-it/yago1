package middlewares

import (
	"nerd/shortener/flags"
	"nerd/shortener/logging"
	"net/http"
	"time"
)

var sugar = flags.GetSugar()

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &logging.ResData{
			Size:   0,
			Status: 0,
		}
		lw := logging.LogResponseWriter{
			ResponseWriter: w,
			ResData:        rd,
		}

		h.ServeHTTP(&lw, r)

		end := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", rd.Status,
			"duration", end,
		)
	})
}
