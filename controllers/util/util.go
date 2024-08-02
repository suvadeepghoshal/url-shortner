package util

import (
    "github.com/a-h/templ"
    "log/slog"
    "net/http"
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
