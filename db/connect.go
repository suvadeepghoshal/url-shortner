package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"url-shortner/controllers/util"
	model "url-shortner/model/type"
)

func ConnectDB() (*gorm.DB, error) {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
	}

	var interpolator model.StringLiteral = util.StringInterpolator{}

	dsn := interpolator.Interpolate("user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_NAME} host=${DB_HOST} port=${DB_PORT} sslmode=disable",
		map[string]string{
			"DB_NAME":     os.Getenv("DB_NAME"),
			"DB_USERNAME": os.Getenv("DB_USERNAME"),
			"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
			"DB_HOST":     os.Getenv("DB_HOST"),
			"DB_PORT":     os.Getenv("DB_PORT"),
		})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
