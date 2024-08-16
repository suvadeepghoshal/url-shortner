package short

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortner/controllers"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"
)

func UrlController(ctx *controllers.ControllerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("inside ShortController")

		urlParams := TYPE.CommonResponse{
			Time:      time.Now(),
			UrlParams: TYPE.UrlParameter{},
		}

		var input TYPE.UrlParameter
		reqParseErr := json.NewDecoder(request.Body).Decode(&input)
		if reqParseErr != nil {
			slog.Error("Unable to parse request body", "err", reqParseErr.Error())
			http.Error(writer, "Unable to parse request body", http.StatusBadRequest)
			return
		}

		if valErr := ctx.Validator.Struct(input); valErr != nil {
			slog.Error("Valid long_url is required in the request", "err", valErr.Error())
			http.Error(writer, "Valid long_url is required in the request", http.StatusBadRequest)
			return
		}

		longUrl := input.LongUrl

		// TODO: may be get the connection and seeding done at start of the application, when the user is in '/' route
		dbConn, conErr := db.ConnectDB()
		if conErr != nil {
			slog.Error("Unable to connect to the database: ", "err", conErr.Error())
			http.Error(writer, "Unable to connect to the database", http.StatusInternalServerError)
			return
		}

		migErr := dbConn.AutoMigrate(&TYPE.Url{})
		if migErr != nil {
			slog.Error("Unable to seed the database: ", "err", migErr.Error())
			http.Error(writer, "Unable to seed the database", http.StatusInternalServerError)
			return
		} else {
			slog.Info("Database seeding successful")
		}

		shortUrl, e := getShortUrl(longUrl)
		if e != nil {
			slog.Error("Unable to get short url", "err", e.Error())
			http.Error(writer, "Unable to get short url", http.StatusInternalServerError)
			return
		}
		slog.Info("ShortController", "short_url_length", len(shortUrl))

		urlParams.UrlParams.ShortUrl = shortUrl
		urlParams.UrlParams.LongUrl = longUrl

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(urlParams)
		if err != nil {
			slog.Error("Unable to write response: ", "err", err.Error())
			http.Error(writer, "Unable to write response", http.StatusInternalServerError)
			return
		}
	}
}

func getShortUrl(L string) (string, error) {
	sUrl := util.ToBase62(rand.Uint64())
	slog.Info("getShortUrl", "s_url", sUrl)
	parse, err := url.Parse(L)
	if err != nil {
		return "", err
	}
	prefix := parse.Scheme
	slog.Debug("getShortUrl", "prefix", prefix)
	var strBuilder strings.Builder
	strBuilder.WriteString(prefix + "://")
	strBuilder.WriteString(strings.ToLower(sUrl))
	return strBuilder.String(), nil
}
