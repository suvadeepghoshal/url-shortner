package home

import (
    "net/http"
    "time"
    "url-shortner/controllers/util"
    TYPE "url-shortner/model/type"
    "url-shortner/view/home"
)

func HomeController(writer http.ResponseWriter, request *http.Request) error {
    _ = TYPE.CommonResponse{
        Time: time.Now(),
    }
    writer.WriteHeader(200)
    return util.Render(writer, request, home.Index())
}
