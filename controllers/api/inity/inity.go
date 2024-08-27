package inity

import (
	"encoding/gob"
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

		pgDbConfig, configErr := db.LoadPgDbConfig()
		if configErr != nil {
			slog.Error("Unable to load pgDbConfig", "configErr", configErr)
			http.Error(writer, configErr.Error(), http.StatusInternalServerError)
			return
		}

		pgDriver := db.PsqlDataBase{DbParams: pgDbConfig}

		dbConn, conErr := pgDriver.Connection()
		if conErr != nil {
			slog.Error("Unable to connect to the database: ", "err", conErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		// Seeding Database at init
		// TODO: may be seed only if tables are not found
		if migErr := dbConn.AutoMigrate(&TYPE.Url{}); migErr != nil {
			slog.Error("Unable to seed the database: ", "err", migErr.Error())
			http.Error(writer, "Unable to seed the database", http.StatusInternalServerError)
			return
		} else {
			slog.Info("Database seeding successful")
		}

		curr := db.GormDB{
			Gorm: dbConn,
		}

		// TODO: It is rare to Close a DB, as the DB handle is meant to be long-lived and shared between many goroutines. It makes sense to close the connection once one user session ends or expires?
		genDB, genErr := curr.GenDB()
		if genErr != nil {
			slog.Error("Unable to generate the generic database instance", "genErr", genErr)
			http.Error(writer, "Unable to generate the generic database instance", http.StatusInternalServerError)
			return
		}

		if util.CloseDbConnection(writer, genDB) {
			return
		}

		secret, e := util.GenerateSessionSecret(32)
		if e != nil {
			slog.Error("Unable to generate session secret", "err", e)
		}
		slog.Debug("InitController", "secret", secret)

		// TODO: Check if the godotenv can be load once in the main and shared across as a state
		envErr := godotenv.Load(".env")
		if envErr != nil {
			slog.Error("Error loading env file", "err", envErr)
			return
		}

		// creating a new store to hold protected resources after auth
		store := sessions.NewCookieStore([]byte(secret))
		slog.Debug("InitController", "cookie_store", store)
		gothic.Store = store
		// This is needed as we will store the expiry time of the auth tokens in the session, and by Go lang uses gob package to encode/decode session data, which is not capable of automatically handling certain types (eg: time.Time)
		gob.Register(time.Time{})

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
