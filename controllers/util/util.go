package util

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/a-h/templ"
	"io"
	"log/slog"
	"net/http"

	"os"
	TYPE "url-shortner/model/type"
)

//const (
//	base         uint64 = 62
//	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
//)

func Render(writer http.ResponseWriter, request *http.Request, componentContext templ.Component) error {
	return componentContext.Render(request.Context(), writer)
}

func Main(httpHandler TYPE.HTTPHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := httpHandler(writer, request); err != nil {
			slog.Error("HTTP handler error", "err", err, "path", request.URL.Path)
		}
	}
}

type StringInterpolator struct {
	sl TYPE.StringLiteral
}

func (si StringInterpolator) Interpolate(template string, variables map[string]string) string {
	f := func(ph string) string {
		return variables[ph]
	}
	return os.Expand(template, f)
}

//func ToBase62(num uint64) string {
//	encoded := ""
//	for num > 0 {
//		r := num % base
//		num = num / base
//		encoded += string(characterSet[r])
//	}
//	return encoded
//}

func CreateMd5Hash(s string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func CloseDbConnection(w http.ResponseWriter, db *sql.DB) bool {
	if e := db.Close(); e != nil {
		slog.Error("Unable to close the database connection: ", "err", e.Error())
		http.Error(w, "Unable to close the database connection", http.StatusInternalServerError)
		return true
	}
	return false
}
func GetCurrDomain(r *http.Request) string {

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

func ParseShortUrl(l, s string, request *http.Request) (string, error) {
	hostname := GetCurrDomain(request)
	slog.Debug("ParseShortController", "hostname", hostname)
	slog.Debug("ParseShortController", "short_url", s)
	var interpolator TYPE.StringLiteral = StringInterpolator{}
	returnStr := interpolator.Interpolate("${HOST_NAME}/${SHORT_URL}", map[string]string{
		"HOST_NAME": hostname,
		"SHORT_URL": s,
	})

	slog.Debug("ParseShortUrl", "return_str", returnStr)
	return returnStr, nil
}
