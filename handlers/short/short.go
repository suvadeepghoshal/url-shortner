package short

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"url-shortner/db"
	"url-shortner/handlers"
	"url-shortner/handlers/util"
	TYPE "url-shortner/model/type"
	"url-shortner/repository"
	"url-shortner/service"
)

func UrlController(ctx *handlers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("UrlController", "reqId", request.Context().Value("req_id"))

		var urlParam TYPE.UrlParameter
		reqParseErr := json.NewDecoder(request.Body).Decode(&urlParam)
		if reqParseErr != nil {
			slog.Error("Unable to parse request body", "err", reqParseErr.Error())
			http.Error(writer, "Unable to parse request body", http.StatusBadRequest)
			return
		}

		if valErr := ctx.Validator.Struct(urlParam); valErr != nil {
			slog.Error("Valid long_url is required in the request", "err", valErr.Error())
			http.Error(writer, "Valid long_url is required in the request", http.StatusBadRequest)
			return
		}

		conn, connErr := db.NewPgDriver().GetConnection()
		if connErr != nil {
			slog.Error("Unable to connect to the database: ", "err", connErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		var url TYPE.Url

		repo := &repository.RepoService{Db: conn}
		srv := &service.UrlService{Repo: repo}

		url.LongUrl = urlParam.LongUrl

		ctx := context.WithValue(request.Context(), "hostname", util.GetCurrDomain(request))
		if e := srv.MakeShortUrl(ctx, &url); e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		urlParam.ShortUrl = url.ShortUrl

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(urlParam)
		if err != nil {
			slog.Error("Unable to write response: ", "err", err.Error())
			http.Error(writer, "Unable to write response", http.StatusInternalServerError)
			return
		}
	}
}
