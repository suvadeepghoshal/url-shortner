package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortner/controllers"
	"url-shortner/controllers/api/inity"
	"url-shortner/controllers/auth/provider/google_prov"
	"url-shortner/controllers/redirect"
	"url-shortner/controllers/short"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	// TODO: integrate with slogenv
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("inside main :: APP STARTED")

	validate := validator.New()

	// TODO: Check if the godotenv can be load once in the main and shared across as a state
	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
		return
	}

	// Initialize Contexts required by the Controllers
	ctx := &controllers.ControllerContext{
		Validator: validate,
	}

	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Heartbeat("/ping"))                      // Gives the status of the application
	router.Use(middleware.AllowContentEncoding("deflate", "gzip")) // AllowContentEncoding enforces a whitelist of request Content-Encoding
	router.Use(middleware.AllowContentType("application/json"))    // AllowContentType enforces a whitelist of request Content-Types
	router.Use(middleware.CleanPath)                               // CleanPath middleware will clean out double slash mistakes from a user's request path

	// API routes
	router.Get("/", inity.InitController(ctx))
	router.Get("/{hash}", redirect.RedirController(ctx))
	router.Post("/url/short", short.UrlController(ctx))
	router.Get("/auth", google_prov.HandleGoogleAuth)
	router.Get("/auth/callback", google_prov.HandleGoogleAuthCallBack)

	err := http.ListenAndServe(os.Getenv("APP_PORT"), router)
	if err != nil {
		slog.Error("App can not be served", "err", err)
	}

}
