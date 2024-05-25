package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gw GzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}(gz)
			r.Body = gz

			h.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(w)
			defer func(gz *gzip.Writer) {
				err := gz.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}(gz)

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			gzw := GzipWriter{
				ResponseWriter: w,
				Writer:         gz,
			}

			h.ServeHTTP(gzw, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
