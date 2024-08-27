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
		//reqId := request.Context().Value("req_id").(string)
		//slog.Info("UrlController", "reqId", reqId)

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

		shortUrl, e := util.CreateMd5Hash(longUrl)
		if e != nil {
			slog.Error("Unable to create short url", "err", e.Error())
			http.Error(writer, "Unable to create short url", http.StatusInternalServerError)
			return
		}
		shortUrl = shortUrl[0:7]
		slog.Info("UrlController", "short_url_hash_length", len(shortUrl))

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

		parsedShortUrl, parseErr := util.ParseShortUrl(longUrl, shortUrl, request)
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
