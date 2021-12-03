package models

import (
	//"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//DB 資料庫連接
var DB *gorm.DB

//Setup 資料庫連接設定
func Setup() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	AutoMigrateAll()
}

//AutoMigrateAll 自動產生 table
func AutoMigrateAll() {
	DB.AutoMigrate(&Class{})
	DB.AutoMigrate(&ClassUser{})
	DB.AutoMigrate(&Problem{})
	DB.AutoMigrate(&Test{})
}
func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
