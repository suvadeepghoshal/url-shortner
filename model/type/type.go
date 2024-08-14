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

type StringLiteral interface {
	Interpolate(template string, variables map[string]string) string
}

type HTTPHandler func(writer http.ResponseWriter, request *http.Request) error

type Url struct {
	id        uint   `gorm:"primaryKey"`
	shortUrl  string `gorm:"unique"`
	longUrl   string
	createdAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	expiry    time.Duration
}
