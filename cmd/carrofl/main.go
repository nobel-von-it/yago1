package main

import (
	"github.com/go-chi/chi/v5"
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

//func carHandle(w http.ResponseWriter, r *http.Request) {
//	if _, err := w.Write([]byte(findCar(chi.URLParam(r, "id")))); err != nil {
//		return
//	}
//}

func brandHandle(w http.ResponseWriter, r *http.Request) {
	list := make([]string, 0)
	brand := strings.ToLower(chi.URLParam(r, "brand"))
	for _, c := range cars {
		if strings.Split(strings.ToLower(c), " ")[0] == brand {
			list = append(list, c)
		}
	}
	io.WriteString(w, strings.Join(list, ", "))
}

func modelHandle(w http.ResponseWriter, r *http.Request) {
	car := strings.ToLower(chi.URLParam(r, "brand") + " " + chi.URLParam(r, "model"))
	for _, c := range cars {
		if strings.ToLower(c) == car {
			io.WriteString(w, c)
			return
		}
	}
	http.Error(w, "unknown model: "+car, http.StatusNotFound)
}

func main() {
	s := chi.NewRouter()

	s.Route("/cars", func(s chi.Router) {
		s.Get("/", carsHandle)
		s.Route("/{brand}", func(s chi.Router) {
			s.Get("/", brandHandle)
			s.Get("/{model}", modelHandle)
		})
	})

	log.Fatalln(http.ListenAndServe(":8080", s))
}
