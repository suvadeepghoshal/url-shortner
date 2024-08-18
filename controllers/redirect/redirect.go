package redirect

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"
)

func RedirController(writer http.ResponseWriter, request *http.Request) {
	slog.Info("inside RedirController")

	var url TYPE.Url

	hash := chi.URLParam(request, "hash")

	dbConn, genObj, conErr := db.ConnectDB()
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

	if util.CloseDbConnection(writer, genObj) {
		return
	}

	parsedShortUrl, parseErr := util.ParseShortUrl(url.LongUrl, url.ShortUrl, request)
	if parseErr != nil {
		slog.Error("Unable to parse the short url", "err", parseErr.Error())
	}
	url.ShortUrl = parsedShortUrl

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(url)
	if err != nil {
		slog.Error("Unable to write response: ", "err", err.Error())
		http.Error(writer, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
