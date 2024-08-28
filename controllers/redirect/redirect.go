package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortner/controllers"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RedirController(_ *controllers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("RedirController", "req_id", request.Context().Value("req_id"))

		var url TYPE.Url

		hash := chi.URLParam(request, "hash")

		pgDbConfig, configErr := db.LoadPgDbConfig()
		if configErr != nil {
			slog.Error("Unable to load pgDbConfig", "configErr", configErr)
			http.Error(writer, configErr.Error(), http.StatusInternalServerError)
			return
		}

		//pgDriver := db.PsqlDataBase{DbParams: pgDbConfig}
		pgDriver := db.NewPsqlDataBase(pgDbConfig)

		dbConn, conErr := pgDriver.Connection()
		if conErr != nil {
			slog.Error("Unable to connect to the database: ", "err", conErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		if e := dbConn.Where("short_url = ?", hash).First(&url).Error; e != nil {
			slog.Error("Unable to find the url", "hash", hash, "err", e.Error())
			if errors.Is(e, gorm.ErrRecordNotFound) {
				http.Error(writer, "url not found", http.StatusNotFound)
				return
			}
			http.Error(writer, "Something went wrong, Please try again", http.StatusInternalServerError)
			return
		}

		if len(url.LongUrl) != 0 {
			slog.Info("RedirController found the url", "hash", hash, "short_url_len", len(url.ShortUrl), "long_url_len", len(url.LongUrl))
		}

		//curr := db.GormDB{
		//	Gorm: dbConn,
		//}

		// TODO: It is rare to Close a DB, as the DB handle is meant to be long-lived and shared between many goroutines. It makes sense to close the connection once one user session ends or expires?
		//genDB, genErr := curr.GenDB()
		genDB, genErr := db.NewGormDB(dbConn).GenDB()
		if genErr != nil {
			slog.Error("Unable to generate the generic database instance", "genErr", genErr)
			http.Error(writer, "Unable to generate the generic database instance", http.StatusInternalServerError)
			return
		}

		if util.CloseDbConnection(writer, genDB) {
			return
		}

		//parsedShortUrl, parseErr := util.ParseShortUrl(url.LongUrl, url.ShortUrl, request)
		//if parseErr != nil {
		//	slog.Error("Unable to parse the short url", "err", parseErr.Error())
		//}
		//url.ShortUrl = parsedShortUrl

		// Giving permanent redirect to the user. Once the user clicks on the short URL, it will automatically take the user to the actual URL
		writer.Header().Set("Location", url.LongUrl)
		//writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusMovedPermanently)
		//err := json.NewEncoder(writer).Encode(url)
		//if err != nil {
		//	slog.Error("Unable to write response: ", "err", err.Error())
		//	http.Error(writer, "Unable to write response", http.StatusInternalServerError)
		//	return
		//}
	}
}
