package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func init() {
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		defServAddr = serverAddress
	}
	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		defBaseUrl = baseUrl
	}
	if storagePath := os.Getenv("STORAGE_PATH"); storagePath != "" {
		defStoragePath = storagePath
	}

	flag.StringVar(&config.ServerAddress, "a", defServAddr, "address for http-server")
	flag.StringVar(&config.BaseUrl, "p", defBaseUrl, "user-input value for pre short url")
	flag.StringVar(&config.StoragePath, "f", defStoragePath, "full path to the storage")
	flag.BoolVar(&config.Same, "s", false, "are base url and server address same")

	if !isTestRun() {
		flag.Parse()
	}

	var logger *zap.Logger
	var err error
	if isTestRun() {
		logger = zap.NewNop()
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	sugar = logger.Sugar()
	if err = events.Load(); err != nil {
		panic(err)
	}

	if !isTestRun() {
		err := logger.Sync()
		if err != nil && !strings.Contains(err.Error(), "stderr") {
			panic(err)
		}
	}

	if config.Same && config.ServerAddress != config.BaseUrl {
		if config.ServerAddress != defServAddr {
			config.BaseUrl = config.ServerAddress
		} else if config.BaseUrl != defBaseUrl {
			config.ServerAddress = config.BaseUrl
		}
	}

	if err = dirExists(config.StoragePath); err != nil {
		panic(err)
	}
}

func dirExists(path string) error {
	if !strings.Contains(path, "/") {
		return errors.New("invalid path: " + path)
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(strings.Split(path, "/")[0], os.ModePerm)
		if err != nil {
			return err
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		return err
	}
	return nil
}

func isTestRun() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "test") {
			return true
		}
	}
	return false
}

type Config struct {
	ServerAddress string
	BaseUrl       string
	StoragePath   string
	Same          bool
}

type ResData struct {
	Status int
	Size   int
}
type LogResponseWriter struct {
	http.ResponseWriter
	ResData *ResData
}

func (r *LogResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResData.Size += size
	return size, err
}

func (r *LogResponseWriter) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.ResData.Status = code
}

const (
	symbols    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLen = 5
	form       = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/api/shorten" method="post">
            <label>Address</label><input type="text" name="url">
            <input type="submit" value="Generate">
        </form>
    </body>
	</html>`
)

var (
	config         Config
	events         Events
	sugar          *zap.SugaredLogger
	defServAddr    = "localhost:8080"
	defBaseUrl     = "localhost:8080"
	defStoragePath = "tmp/short-url-db.json"
)

//type Event struct {
//	Uuid        string `json:"uuid"`
//	ShortUrl    string `json:"short_url"`
//	OriginalUrl string `json:"original_url"`
//}
//
//func (e Event) Save(filename string) error {
//	data, err := json.MarshalIndent(e, "", "	")
//	if err != nil {
//		return err
//	}
//	return os.WriteFile(filename, data, 0644)
//}
//
//func (e *Event) Load(filename string) error {
//	data, err := os.ReadFile(filename)
//	if err != nil {
//		return err
//	}
//	if err = json.Unmarshal(data, e); err != nil {
//		return err
//	}
//	return nil
//}
//
//func ReadAllEvents() {
//	event := &Event{}
//	for event.Load(config.StoragePath) != nil {
//		AddMap(shoring, event.OriginalUrl, event.ShortUrl)
//	}
//}

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
	return config.BaseUrl + "/" + str
}

func Info() {
	sugar.Info("hello", "world")
}

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	url := chi.URLParam(r, "url")
	if len(url) == 0 {
		url = r.FormValue("url")
		if len(url) == 0 {
			http.Error(w, "url required", http.StatusBadRequest)
			return
		}
	}
	short := GenShortUrl(defaultLen)
	addr := ToAddr(short)
	events.Add(short, url)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("From: " + addr + " To: " + url + "\n"))

	sugar.Infoln("From:", addr, "To:", url)
}

func GetForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(form))
}

//func PostFormHandler(w http.ResponseWriter, r *http.Request) {
//	url := r.URL.Query().Get("url")
//	if r.Method == http.MethodPost || url != "" {
//		if url == "" {
//			url = r.FormValue("url")
//			if url == "" {
//				http.Error(w, "url is empty", http.StatusBadRequest)
//				return
//			}
//		}
//		short := GenShortUrl(defaultLen)
//		addr := ToAddr(short)
//
//		AddMap(shoring, url, short)
//
//		w.Header().Set("content-type", "text/plain")
//		w.WriteHeader(http.StatusCreated)
//
//		_, err := w.Write([]byte(addr))
//		if err != nil {
//			return
//		}
//		sugar.Infoln("From:", addr, "To:", url)
//	} else if r.Method == http.MethodGet {
//		_, _ = w.Write([]byte(form))
//	} else {
//		http.Error(w, "shiiit", http.StatusBadRequest)
//	}
//}

type RequestData struct {
	Url string `json:"url"`
}

type ResponseData struct {
	Result string `json:"result"`
}

func JsonPostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Need post method", http.StatusMethodNotAllowed)
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

	events.Add(short, req.Url)

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
		addr := events.Find(strings.TrimPrefix(r.URL.String(), "/")).OriginalUrl
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

	for _, e := range events.Events {
		ent := fmt.Sprintf("Addr [%s] with shortform [%s]\n", e.OriginalUrl, e.ShortUrl)
		_, err := w.Write([]byte(ent))
		if err != nil {
			return
		}
	}
}

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &ResData{
			Size:   0,
			Status: 0,
		}
		lw := LogResponseWriter{
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

type Event struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type Events struct {
	Events []Event `json:"events"`
}

func (es *Events) Save() error {
	data, err := json.MarshalIndent(es, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.StoragePath, data, 0666)
}

func (es *Events) Load() error {
	data, err := os.ReadFile(config.StoragePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		events.Events = make([]Event, 0)
		return nil
	}
	return json.Unmarshal(data, es)
}

func (es *Events) Add(short, url string) {
	es.Events = append(es.Events, Event{
		Uuid:        short,
		ShortUrl:    ToAddr(short),
		OriginalUrl: url,
	})
	err := es.Save()
	if err != nil {
		sugar.Infow("error on save", "err", err)
	}
}

func (es *Events) Find(uuid string) *Event {
	for _, e := range es.Events {
		if e.Uuid == uuid {
			return &e
		}
	}
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Handle("/", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(GetForm))))
	r.Handle("/{id}", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(GetHandler))))
	r.Handle("/getall", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(GetAllHandler))))
	r.Handle("/api/shorten", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(PostFormHandler))))
	r.Handle("/api/shorten/json", LoggerMiddleware(GzipMiddleware(http.HandlerFunc(JsonPostFormHandler))))

	if sugar == nil {
		panic("Logger not initialized")
	}
	if events.Events == nil {
		panic("Events not initialized")
	}

	sugar.Infow("Starting server",
		"Server Address", config.ServerAddress,
		"Base URL", config.BaseUrl,
		"Is same", config.Same,
	)

	if err := http.ListenAndServe(config.ServerAddress, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
}
