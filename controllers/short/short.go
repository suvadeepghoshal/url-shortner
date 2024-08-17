package short

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

		shortUrl, e := util.CreateMd5Hash(longUrl)
		if e != nil {
			slog.Error("Unable to create short url", "err", e.Error())
			http.Error(writer, "Unable to create short url", http.StatusInternalServerError)
			return
		}
		shortUrl = shortUrl[0:7]
		slog.Info("ShortController", "short_url", shortUrl)

		reqObj := TYPE.Url{ShortUrl: shortUrl[0:7], LongUrl: longUrl}

		result := dbConn.Create(&reqObj)
		if result.Error != nil {
			slog.Error("Unable to store url data in DB", "err", result.Error.Error())
			if strings.Contains(result.Error.Error(), "duplicate key value violates") {
				http.Error(writer, "Url already exists", http.StatusConflict)
			} else {
				http.Error(writer, "Unable to store url in DB", http.StatusInternalServerError)
			}
			return
		}

		if result.RowsAffected > 0 {
			slog.Info("Url data is stored in the DB")
		}

		parsedShortUrl, parseErr := parseShortUrl(longUrl, shortUrl, request)
		if parseErr != nil {
			slog.Error("Unable to parse short url", "err", parseErr.Error(), "url", longUrl)
			http.Error(writer, "Unable to parse short url", http.StatusInternalServerError)
			return
		}

		urlParams.UrlParams.ShortUrl = parsedShortUrl
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

func parseShortUrl(l, s string, request *http.Request) (string, error) {
	host := request.Host
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}
	var interpolator TYPE.StringLiteral = util.StringInterpolator{}
	returnStr := interpolator.Interpolate("${SCHEME}://${HOST}/${SHORT_URL}", map[string]string{
		"SCHEME":    scheme,
		"HOST":      host,
		"SHORT_URL": s,
	})

	slog.Debug("getShortUrl", "return_str", returnStr)
	return returnStr, nil
}
