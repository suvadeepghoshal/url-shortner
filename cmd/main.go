package main

import (
	"log/slog"
	"net/http"
	"url-shortner/controllers/api/inity"
	"url-shortner/controllers/home"
	"url-shortner/controllers/short"
	"url-shortner/controllers/util"

	"github.com/go-chi/chi/v5"
)

func main() {
	slog.Info("inside main :: APP STARTED")
	router := chi.NewMux()
	// views
	router.Get("/", util.Main(home.HomeController))

	// API routes
	router.Get("/init", inity.InitController)
	router.Post("/url/short", short.ShortController)

	err := http.ListenAndServe(":1323", router)
	if err != nil {
		slog.Error("inside main :: App can not be served: ", "err", err)
	}
}
