package db

import (
	"gorm.io/gorm"
	TYPE "url-shortner/model/type"
)

func CreateOperation(db *gorm.DB, url *TYPE.Url) error {
	result := db.Create(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
