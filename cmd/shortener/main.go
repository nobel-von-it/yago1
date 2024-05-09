package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"nerd/yago1/cmd/shortener/config"
	"net/http"
	"strings"
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
		fmt.Println("address added:", addr, "->", url)
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
			_, err := w.Write([]byte("This id not found"))
			if err != nil {
				return
			}
		} else {
			w.Header().Set("location", addr)
			w.WriteHeader(http.StatusTemporaryRedirect)
			fmt.Println("got address")
		}
	} else {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(400)
		_, err := w.Write([]byte("Incorrect request"))
		if err != nil {
			return
		}
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

func main() {
	r := chi.NewRouter()

	r.HandleFunc("/", PostFormHandler)
	r.Get("/{id}", GetHandler)
	r.Get("/getall", GetAllHandler)

	log.Println("cfg.ServerAddress =", cfg.ServerAddress)
	log.Println("cfg.BaseUrl =", cfg.BaseUrl)
	log.Fatalln(http.ListenAndServe(cfg.ServerAddress, r))
}
