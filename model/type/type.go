package model

import (
	"net/http"
	"time"
)

type CommonResponse struct {
	Time      time.Time    `json:"time"`
	UrlParams UrlParameter `json:"urlParams"`
}

type UrlParameter struct {
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
}

type HTTPHandler func(writer http.ResponseWriter, request *http.Request) error
