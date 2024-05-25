package main

import (
	"github.com/go-chi/chi/v5"
	"nerd/shortener/flags"
	"nerd/shortener/handlers"
	"nerd/shortener/middlewares"
	"nerd/shortener/storage"
	"net/http"
)

var (
	config = flags.GetConfig()
	sugar  = flags.GetSugar()
	events = storage.GetEvents()
)

func main() {
	r := chi.NewRouter()
	r.Handle("/", middlewares.LoggerMiddleware(middlewares.GzipMiddleware(http.HandlerFunc(handlers.GetForm))))
	r.Handle("/{id}", middlewares.LoggerMiddleware(middlewares.GzipMiddleware(http.HandlerFunc(handlers.GetHandler))))
	r.Handle("/getall", middlewares.LoggerMiddleware(middlewares.GzipMiddleware(http.HandlerFunc(handlers.GetAllHandler))))
	r.Handle("/api/shorten", middlewares.LoggerMiddleware(middlewares.GzipMiddleware(http.HandlerFunc(handlers.PostFormHandler))))
	r.Handle("/api/shorten/storage", middlewares.LoggerMiddleware(middlewares.GzipMiddleware(http.HandlerFunc(handlers.JsonPostFormHandler))))

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
