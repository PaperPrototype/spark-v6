package db

import (
	"main/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func Setup() {
	url := helpers.GetDatabaseURL()
	db, err := gorm.Open(postgres.Open(url))
	if err != nil {
		panic(err)
	}

	GormDB = db

	migrate()
}
