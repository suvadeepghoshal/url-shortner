package inity

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"url-shortner/controllers"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func InitController(_ *controllers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		slog.Info("inside InitController")
		commonResponse := TYPE.CommonResponse{
			Time: time.Now(),
		}

		dbConn, genObj, conErr := db.ConnectDB()
		if conErr != nil {
			slog.Error("Unable to connect to the database: ", "err", conErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		// Seeding Database at init
		if migErr := dbConn.AutoMigrate(&TYPE.Url{}); migErr != nil {
			slog.Error("Unable to seed the database: ", "err", migErr.Error())
			http.Error(writer, "Unable to seed the database", http.StatusInternalServerError)
			return
		} else {
			slog.Info("Database seeding successful")
		}

		if util.CloseDbConnection(writer, genObj) {
			return
		}

		// TODO: make the secret for the session and store in the env
		secret, e := util.GenerateSessionSecret(32)
		if e != nil {
			slog.Error("Unable to generate session secret", "err", e)
		}

		// TODO: Check if the godotenv can be load once in the main and shared across as a state
		envErr := godotenv.Load(".env")
		if envErr != nil {
			slog.Error("Error loading env file", "err", envErr)
			return
		}

		store := sessions.NewCookieStore([]byte(secret))
		gothic.Store = store

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(commonResponse)
		if err != nil {
			slog.Error("Unable to write response: ", "err", err.Error())
			http.Error(writer, "Unable to write response", http.StatusInternalServerError)
			return
		}
	}
}
