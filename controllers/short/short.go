package short

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"url-shortner/db"
	TYPE "url-shortner/model/type"
)

func ShortController(writer http.ResponseWriter, request *http.Request) {
	slog.Info("inside ShortController")
	urlParams := TYPE.CommonResponse{
		Time: time.Now(),
		UrlParams: TYPE.UrlParameter{
			ShortUrl: "https://dummy.com",
		},
	}
	// TODO: refactor the below first block to not hard code the response from url/short endpoint
	var input TYPE.UrlParameter
	reqParseErr := json.NewDecoder(request.Body).Decode(&input)
	if reqParseErr != nil {
		slog.Error("inside ShortController :: unable to parse request body", "err", reqParseErr.Error())
		return
	}
	urlParams.UrlParams.LongUrl = input.LongUrl

	dbConn, conErr := db.ConnectDB()
	if conErr != nil {
		slog.Error("inside short controller :: unable to connect to the database: ", "err", conErr.Error())
	}

	migErr := dbConn.AutoMigrate(&TYPE.Url{})
	if migErr != nil {
		slog.Error("inside short controller :: unable to seed the database: ", "err", migErr.Error())
	}

	// write the short URL as a response to the user

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(urlParams)
	if err != nil {
		slog.Error("inside short controller :: unable to write response: ", "err", err.Error())
		return
	}
}
