package redirect

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"url-shortner/db"
	"url-shortner/handlers"
	TYPE "url-shortner/model/type"
	"url-shortner/repository"
	"url-shortner/service"
)

func RedirController(_ *handlers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("RedirController", "req_id", request.Context().Value("req_id"))

		var url TYPE.Url

		hash := chi.URLParam(request, "hash")

		conn, connErr := db.NewPgDriver().GetConnection()

		if connErr != nil {
			slog.Error("Unable to connect to the database: ", "err", connErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		repo := &repository.RepoService{
			Db: conn,
		}

		ctx := context.WithValue(request.Context(), "hash", hash)
		srv := &service.UrlService{Repo: repo}

		if e := srv.GetLongUrl(ctx, &url); e != nil {
			slog.Error("RedirController", "err", e.Error())
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("RedirController found the url", "hash", hash, "short_url_len", len(url.ShortUrl), "long_url_len", len(url.LongUrl))

		// Giving permanent redirect to the user. Once the user clicks on the short URL, it will automatically take the user to the actual URL
		writer.Header().Set("Location", url.LongUrl)
		writer.WriteHeader(http.StatusMovedPermanently)
	}
}
