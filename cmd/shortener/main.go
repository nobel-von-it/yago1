package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
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

func TestAddMap(t *testing.T) {
	testMap := make(map[string]string)
	key := "hello"
	value := "world"

	AddMap(testMap, key, value)
	if testMap[key] != value {
		t.Errorf("expected key %s to have value %s, but got %s", key, value, testMap[key])
	}

	new_value := "world!!!!"
	AddMap(testMap, key, new_value)
	if testMap[key] != new_value {
		t.Errorf("expected key %s to have value %s, but got %s", key, new_value, testMap[key])
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

func TestFindVal(t *testing.T) {
	testMap := map[string]string{"k1": "v1", "k2": "v2"}

	value := "v1"
	key := FindVal(testMap, value)
	if key != "k1" {
		t.Errorf("expected to find key k1 for value %s, but got %s", value, key)
	}

	value = "asldkjf"
	if key != "" {
		t.Errorf("expected empty string for non-existent value %s, but got %s", value, key)
	}
}

func GenShortUrl(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func TestGenShortUrl(t *testing.T) {
	url := GenShortUrl(defaultLen)
	if len(url) != defaultLen {
		t.Errorf("expected URL length to be %d, but got %d", defaultLen, len(url))
	}

	specLen := 10
	url = GenShortUrl(specLen)
	if len(url) != specLen {
		t.Errorf("expected URL length to be %d, but got %d", specLen, len(url))
	}

	for i := 0; i > 10; i++ {
		url = GenShortUrl(defaultLen)
		for _, c := range url {
			if !strings.ContainsAny(string(c), symbols) {
				t.Errorf("generated URL number %d contains invalid char %c", i, c)
			}
		}
	}
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
