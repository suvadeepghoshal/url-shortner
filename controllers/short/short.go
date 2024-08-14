package short

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"url-shortner/controllers/util"
	"url-shortner/db"
	TYPE "url-shortner/model/type"
)

const (
	encodingChunkSize = 2
	decodingChunkSize = 3
)

func ShortController(writer http.ResponseWriter, request *http.Request) {
	slog.Info("inside ShortController")
	urlParams := TYPE.CommonResponse{
		Time:      time.Now(),
		UrlParams: TYPE.UrlParameter{},
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
	} else {
		slog.Info("inside Short Controller :: seeding successful")
	}

	shortUrl, e := getShortUrl(input.LongUrl)
	if e != nil {
		slog.Error("inside short controller :: not able to get short URL: ", "err", e.Error())
	}

	slog.Info("inside Short Controller :: length of", "shorUrl", len(shortUrl))
	urlParams.UrlParams.ShortUrl = shortUrl

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(urlParams)
	if err != nil {
		slog.Error("inside short controller :: unable to write response: ", "err", err.Error())
		return
	}
}

func getShortUrl(url string) (string, error) {
	bytes := []byte(url)

	byteLen := len(bytes)

	fmt.Println("length of the url: ", len(url))
	fmt.Println("length of the byte stream: ", len(bytes))

	var strBuilder strings.Builder

	for i := 0; i < byteLen; i += encodingChunkSize {
		chunk := bytes[i:min(i+encodingChunkSize, byteLen)]
		hx := hex.EncodeToString(chunk)
		parseUint, err := strconv.ParseUint(hx, 16, 64)
		if err != nil {
			return "", err
		}
		s := util.PadLeft(util.ToBase62(parseUint), "0", decodingChunkSize)
		strBuilder.WriteString(s)
	}
	return strBuilder.String(), nil
}
