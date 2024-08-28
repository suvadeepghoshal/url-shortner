package db

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"url-shortner/controllers/util"
	model "url-shortner/model/type"
)

func (p *PsqlDataBase) Connection() (*gorm.DB, error) {
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

func (g *GormDB) GenDB() (*sql.DB, error) {
	w, err := g.Gorm.DB()
	if err != nil {
		return nil, err
	}
	return w, nil
}
