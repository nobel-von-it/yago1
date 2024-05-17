package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"nerd/shortener/config"
	"nerd/shortener/logging"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	symbols    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLen = 5
	form       = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/" method="post">
            <label>Address</label><input type="text" name="url">
            <input type="submit" value="Generate">
        </form>
    </body>
	</html>`
)

var (
	shoring = make(map[string]string)
	cfg     = config.ParseArgs()
)

func AddMap(mp map[string]string, key, value string) {
	if key == "" || value == "" || key == value {
		return
	}
	mp[key] = value
}

func FindVal(mp map[string]string, val string) string {
	for k, v := range mp {
		if val == v {
			return k
		}
	}
	return ""
}

func GenShortUrl(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func ToAddr(str string) string {
	return cfg.BaseUrl + "/" + str
}

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if r.Method == http.MethodPost || url != "" {
		if url == "" {
			url = r.FormValue("url")
			if url == "" {
				http.Error(w, "url is empty", http.StatusBadRequest)
				return
			}
		}
		short := GenShortUrl(defaultLen)
		addr := ToAddr(short)

		AddMap(shoring, url, short)

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err := w.Write([]byte(addr))
		if err != nil {
			return
		}
		sugar.Infoln("From:", addr, "To:", url)
	} else if r.Method == http.MethodGet {
		_, _ = w.Write([]byte(form))
	} else {
		http.Error(w, "shiiit", http.StatusBadRequest)
	}
}

type RequestData struct {
	Url string `json:"url"`
}

type ResponseData struct {
	Result string `json:"result"`
}

func JsonPostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Need post method", http.StatusBadRequest)
		return
	}

	var req RequestData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Parse json error", http.StatusBadRequest)
		return
	}

	sugar.Infow("Got json:", "url", req.Url)

	if req.Url == "" {
		http.Error(w, "Url is empty", http.StatusBadRequest)
		return
	}
	short := GenShortUrl(defaultLen)
	addr := ToAddr(short)

	AddMap(shoring, req.Url, short)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := ResponseData{
		Result: addr,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
		return
	}
	sugar.Infoln("From:", addr, "To:", req.Url)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		addr := FindVal(shoring, strings.TrimPrefix(r.URL.String(), "/"))
		if addr == "" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(400)
			_, _ = w.Write([]byte("This id not found"))
		} else {
			w.Header().Set("location", addr)
			w.WriteHeader(http.StatusTemporaryRedirect)
			sugar.Infow("Redirecting...")
		}
	} else {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(400)
		_, _ = w.Write([]byte("Incorrect request"))
	}
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

	for k, v := range shoring {
		ent := fmt.Sprintf("Addr [%s] with shortform [%s]\n", k, v)
		_, err := w.Write([]byte(ent))
		if err != nil {
			return
		}
	}
}

var sugar zap.SugaredLogger

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

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	sugar = *logger.Sugar()

	r := chi.NewRouter()

	r.Handle("/", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(PostFormHandler))))
	r.Handle("/{id}", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(GetHandler))))
	r.Handle("/getall", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(GetAllHandler))))
	r.Handle("/api/shorten", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(JsonPostFormHandler))))

	sugar.Infow("Starting server",
		"Server Address", cfg.ServerAddress,
		"Base URL", cfg.BaseUrl,
		"Is same", cfg.Same,
	)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
