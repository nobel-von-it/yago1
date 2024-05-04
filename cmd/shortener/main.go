package main

import (
	"fmt"
	"math/rand"
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
            <label>Address</label><input type="text" name="address">
            <input type="submit" value="Generate">
        </form>
    </body>
	</html>`
)

var shoring map[string]string = make(map[string]string)

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

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.FormValue("address") != "" {
		addr := r.FormValue("address")
		short := GenShortUrl(defaultLen)
		AddMap(shoring, addr, short)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(r.Host + "/" + short))
		if err != nil {
			return
		}
		fmt.Println("address added")
	} else {
		_, err := w.Write([]byte(form))
		if err != nil {
			return
		}
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
	if r.Method == http.MethodGet {
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

}

func main() {
	s := http.NewServeMux()

	s.HandleFunc("/", PostHandler)
	s.HandleFunc("/getall", GetAllHandler)
	s.HandleFunc("/{id}", GetHandler)

	err := http.ListenAndServe(":8080", s)
	if err != nil {
		panic(err)
	}
}
