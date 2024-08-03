package model

import (
    "net/http"
    "time"
)

type CommonResponse struct {
    Time time.Time `json:"time"`
}

type HTTPHandler func(writer http.ResponseWriter, request *http.Request) error
