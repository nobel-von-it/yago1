package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

const form = `<html>
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
const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const defaultLen = 5

var shoring map[string]string = make(map[string]string)

func AddMap(mp map[string]string, key, value string) {
	_, ok := mp[key]
	if !ok {
		mp[key] = value
	}
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
	if r.Method == http.MethodPost {
		addr := r.FormValue("address")
		short := GenShortUrl(defaultLen)
		AddMap(shoring, addr, short)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(r.Host + "/" + short))
		fmt.Println("address added")
	} else {
		w.Write([]byte(form))
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
			fmt.Println("got address")
		}
	} else {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte("Incorrect request"))
	}
}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)

		for k, v := range shoring {
			ent := fmt.Sprintf("Addr [%s] with shortform [%s]\n", k, v)
			w.Write([]byte(ent))
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
