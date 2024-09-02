package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	TYPE "url-shortner/model/type"
)

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

func CreateMd5Hash(s string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GetCurrDomain(r *http.Request) string {

	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

func ParseShortUrl(s, hostname string) string {
	slog.Debug("ParseShortController", "hostname", hostname)
	slog.Debug("ParseShortController", "short_url", s)

	returnStr := fmt.Sprintf("%s/%s", hostname, s)

	slog.Debug("ParseShortUrl", "return_str", returnStr)
	return returnStr
}

func GenerateSessionSecret(length int) (string, error) {
	secret := make([]byte, length) // creates a byte slice called secret with length number of elements. Each element is initialized to zero (0x00).

	if _, e := rand.Read(secret); e != nil {
		return "", e
	}

	return hex.EncodeToString(secret), nil
}
