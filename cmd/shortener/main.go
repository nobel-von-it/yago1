package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"math/rand"
	"nerd/yago1/cmd/shortener/config"
	"nerd/yago1/cmd/shortener/logging"
	"net/http"
	"strings"
	"time"
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

var shoring = make(map[string]string)
var cfg = config.ParseArgs()

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

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "url is empty", http.StatusBadRequest)
			return
		}
		short := GenShortUrl(defaultLen)
		addr := cfg.BaseUrl + "/" + short

		AddMap(shoring, url, short)

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err := w.Write([]byte(addr))
		if err != nil {
			return
		}
		sugar.Infoln("From:", addr, "To:", url)
	} else if r.Method == http.MethodGet {
		w.Write([]byte(form))
	} else {
		http.Error(w, "shiiit", http.StatusBadRequest)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		addr := FindVal(shoring, strings.TrimPrefix(r.URL.String(), "/"))
		if addr == "" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(400)
			w.Write([]byte("This id not found"))
		} else {
			w.Header().Set("location", addr)
			w.WriteHeader(http.StatusTemporaryRedirect)
			sugar.Infow("Redirecting...")
		}
	} else {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte("Incorrect request"))
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
			"size", rd.Size,
		)
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

	r.Handle("/", LoggerMiddleware(http.HandlerFunc(PostFormHandler)))
	r.Handle("/{id}", LoggerMiddleware(http.HandlerFunc(GetHandler)))
	r.Handle("/getall", LoggerMiddleware(http.HandlerFunc(GetAllHandler)))

	sugar.Infow("Starting server",
		"Server Address", cfg.ServerAddress,
		"Base URL", cfg.BaseUrl,
		"Is same", cfg.Same,
	)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
