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
