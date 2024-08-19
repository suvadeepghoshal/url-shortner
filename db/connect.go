package db

import (
	"database/sql"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"time"
	"url-shortner/controllers/util"
	model "url-shortner/model/type"
)

func ConnectDB() (*gorm.DB, *sql.DB, error) {
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
		return nil, nil, err
	}

	sqlDB, err1 := db.DB()
	if err1 != nil {
		return nil, nil, err1
	}

	//Maintaining connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, sqlDB, nil
}

func (p PsqlDataBase) Connection() (*gorm.DB, error) {
	var interpolator model.StringLiteral = util.StringInterpolator{}

	dsn := interpolator.Interpolate("user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_NAME} host=${DB_HOST} port=${DB_PORT} sslmode=disable",
		map[string]string{
			"DB_NAME":     p.DbParams.DbName,
			"DB_USERNAME": p.DbParams.DbUsername,
			"DB_PASSWORD": p.DbParams.DbPassword,
			"DB_HOST":     p.DbParams.DbHost,
			"DB_PORT":     p.DbParams.DbPort,
		})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (g GormDB) GenDB() (*sql.DB, error) {
	w, err := g.Gorm.DB()
	if err != nil {
		return nil, err
	}
	return w, nil
}
