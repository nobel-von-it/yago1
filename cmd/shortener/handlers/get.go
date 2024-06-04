package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

const form = `<html>
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

func GetForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(form))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var addr string
		ev := events.Find(strings.TrimPrefix(r.URL.String(), "/"))
		if ev != nil {
			addr = ev.OriginalUrl
		}
		if addr == "" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("This id not found"))
		} else {
			w.Header().Set("location", addr)
			w.WriteHeader(http.StatusTemporaryRedirect)
			sugar.Infow("Redirecting...")
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

func GetPing(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", config.DataBase)
	if err != nil {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error to connect database"))
	} else {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("connect database success"))
	}
	defer db.Close()
}
