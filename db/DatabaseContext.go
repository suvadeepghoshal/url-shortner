package db

import (
	"database/sql"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log/slog"
	"os"
	model "url-shortner/model/type"
)

type DBDriver interface {
	Connection() (*gorm.DB, error)
	GenDB() (*sql.DB, error)
}

type PsqlDataBase struct {
	DbParams model.DbParams
}

// GormDB this type helps to create a SQL DB (https://pkg.go.dev/database/sql#DB) instance from an active gorm DB instance
type GormDB struct {
	Gorm *gorm.DB
}

func NewGormDB(gdb *gorm.DB) *GormDB {
	return &GormDB{gdb}
}

func NewPsqlDataBase(dbp model.DbParams) *PsqlDataBase {
	return &PsqlDataBase{dbp}
}

// LoadPgDbConfig Load multiple DB configs based on requirements
func LoadPgDbConfig() (model.DbParams, error) {
	// TODO: check if we can get the godotenv as a shared state available globally
	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
		return model.DbParams{}, envErr
	}

	return model.DbParams{
		DbName:     os.Getenv("DB_NAME"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT")}, nil
}
