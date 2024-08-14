package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"url-shortner/controllers/api/inity"
	"url-shortner/controllers/home"
	"url-shortner/controllers/short"
	"url-shortner/controllers/util"
)

func main() {
	// TODO: integrate with slogenv
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("inside main :: APP STARTED")

	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
	}

	router := chi.NewMux()

	// views
	router.Get("/", util.Main(home.HomeController))

	// API routes
	router.Get("/init", inity.InitController)
	router.Post("/url/short", short.ShortController)

	err := http.ListenAndServe(os.Getenv("APP_PORT"), router)
	if err != nil {
		slog.Error("App can not be served", "err", err)
	}

}
