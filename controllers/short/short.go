package short

import (
	"encoding/hex"
	"encoding/json"
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

	var input TYPE.UrlParameter
	reqParseErr := json.NewDecoder(request.Body).Decode(&input)
	if reqParseErr != nil {
		slog.Error("Unable to parse request body", "err", reqParseErr.Error())
		return
	}

	longUrl := input.LongUrl

	dbConn, conErr := db.ConnectDB()
	if conErr != nil {
		slog.Error("Unable to connect to the database: ", "err", conErr.Error())
	}

	migErr := dbConn.AutoMigrate(&TYPE.Url{})
	if migErr != nil {
		slog.Error("Unable to seed the database: ", "err", migErr.Error())
	} else {
		slog.Info("Database seeding successful")
	}

	shortUrl, e := getShortUrl(longUrl)
	if e != nil {
		slog.Error("Not able to get short URL: ", "err", e.Error())
	}

	slog.Info("ShortController", "short_url_length", len(shortUrl))

	urlParams.UrlParams.ShortUrl = shortUrl
	urlParams.UrlParams.LongUrl = longUrl

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(urlParams)
	if err != nil {
		slog.Error("Unable to write response: ", "err", err.Error())
		return
	}
}

func getShortUrl(url string) (string, error) {
	bytes := []byte(url)

	byteLen := len(bytes)

	slog.Debug("getShortUrl", "url_length", len(url))
	slog.Debug("getShortUrl", "byte_stream_length", byteLen)

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
