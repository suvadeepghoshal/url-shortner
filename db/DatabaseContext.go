package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	model "url-shortner/model/type"
)

type Driver interface {
	GetConnection() (*gorm.DB, error)
}

type PgDriver struct {
	config model.DbParams
}

func (d *PgDriver) GetConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", d.config.DbUsername, d.config.DbPassword, d.config.DbName, d.config.DbHost, d.config.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPgDriver() *PgDriver {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
		return &PgDriver{}
	}

	return &PgDriver{config: model.DbParams{
		DbName:     os.Getenv("DB_NAME"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
	}}
}
