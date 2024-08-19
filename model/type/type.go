package model

import (
	"gorm.io/gorm"
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

type DBDriver interface {
	Connection() (*gorm.DB, error)
}
type DbParams struct {
	DbName     string
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     string
}

type HTTPHandler func(writer http.ResponseWriter, request *http.Request) error

type Url struct {
	ID        uint   `gorm:"primaryKey"`
	ShortUrl  string `gorm:"unique"`
	LongUrl   string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Expiry    time.Duration
}
