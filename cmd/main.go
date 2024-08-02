package main

import (
    "log/slog"
    "net/http"
    //    "net/http"
    //    "time"
    "url-shortner/controllers/home"
    "url-shortner/controllers/util"
    "github.com/go-chi/chi/v5"
    //    TYPE "url-shortner/model/type"
)

func main() {
    slog.Info("inside main :: APP STARTED")
    //    commonResponse := TYPE.CommonResponse{
    //        Time: time.Now(),
    //    }
    router := chi.NewMux()
    router.Get("/", util.Main(home.HomeController))
    err := http.ListenAndServe(":1323", router)
    if err != nil {
        slog.Error("inside main :: App can not be served")
    }
}
