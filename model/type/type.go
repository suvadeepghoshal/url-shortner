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
	LongUrl  string `json:"long_url" validate:"required,url"`
}

type StringLiteral interface {
	Interpolate(template string, variables map[string]string) string
}

type DbParams struct {
	DbName     string
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     string
}

// HTTPHandler TODO: shall we move it to util.go to sync? (abstract func)
type HTTPHandler func(writer http.ResponseWriter, request *http.Request) error

type Url struct {
	ID        uint   `gorm:"primaryKey"`
	ShortUrl  string `gorm:"unique"`
	LongUrl   string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Expiry    time.Duration
}

type User struct {
	ID            string
	Email         string
	Picture       string
	VerifiedEmail bool
}
