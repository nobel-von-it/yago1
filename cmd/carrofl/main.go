package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

var cars = map[string]string{
	"id1": "Renault Logan",
	"id2": "Renault Duster",
	"id3": "BMW X6",
	"id4": "BMW M5",
	"id5": "VW Passat",
	"id6": "VW Jetta",
	"id7": "Audi A4",
	"id8": "Audi Q7",
}

func initList() (list []string) {
	for _, c := range cars {
		list = append(list, c)
	}
	return
}

func findCar(id string) string {
	if c, ok := cars[id]; ok {
		return c
	}
	return "unknown id " + id
}

func carsHandle(w http.ResponseWriter, r *http.Request) {
	list := initList()
	if _, err := io.WriteString(w, strings.Join(list, ", ")); err != nil {
		return
	}
}

func carHandle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is missed", http.StatusBadRequest)
		return
	}
	if _, err := w.Write([]byte(findCar(id))); err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}

func main() {
	s := http.NewServeMux()
	s.HandleFunc("/cars", carsHandle)
	s.HandleFunc("/car", carHandle)

	log.Fatalln(http.ListenAndServe(":8080", s))
}
