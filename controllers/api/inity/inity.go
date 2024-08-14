package inity

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	TYPE "url-shortner/model/type"
)

func InitController(writer http.ResponseWriter, _ *http.Request) {
	slog.Info("inside InitController")
	commonResponse := TYPE.CommonResponse{
		Time: time.Now(),
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(commonResponse)
	if err != nil {
		slog.Error("Unable to write response: ", "err", err.Error())
	}
}
