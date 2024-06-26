package logging

import "net/http"

type (
	ResData struct {
		Status int
		Size   int
	}

	LogResponseWriter struct {
		http.ResponseWriter
		ResData *ResData
	}
)

func (r *LogResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResData.Size += size
	return size, err
}

func (r *LogResponseWriter) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.ResData.Status = code
}
