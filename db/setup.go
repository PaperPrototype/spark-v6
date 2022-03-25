package db

import (
	"main/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB

func Setup() {
	url := helpers.GetDatabaseURL()
	db, err := gorm.Open(postgres.Open(url))
	if err != nil {
		panic(err)
	}

	gormDB = db

	migrate()
}
