package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"url-shortner/controllers"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"
)

func RedirController(_ *controllers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("RedirController", "req_id", request.Context().Value("req_id"))

		var url TYPE.Url

		hash := chi.URLParam(request, "hash")

		conn, connErr := db.NewPgDriver().GetConnection()
		if connErr != nil {
			slog.Error("RedirController", "conErr", connErr)
		}

		if connErr != nil {
			slog.Error("Unable to connect to the database: ", "err", connErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		if e := conn.Where("short_url = ?", hash).First(&url).Error; e != nil {
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

		parsedShortUrl, parseErr := util.ParseShortUrl(url.LongUrl, url.ShortUrl, request)
		if parseErr != nil {
			slog.Error("Unable to parse the short url", "err", parseErr.Error())
		}
		url.ShortUrl = parsedShortUrl

		// Giving permanent redirect to the user. Once the user clicks on the short URL, it will automatically take the user to the actual URL
		writer.Header().Set("Location", url.LongUrl)
		writer.WriteHeader(http.StatusMovedPermanently)
	}
}
