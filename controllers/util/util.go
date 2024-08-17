package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/a-h/templ"
	"io"
	"log/slog"
	"net/http"
	"os"
	TYPE "url-shortner/model/type"
)

const (
	base         uint64 = 62
	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
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
