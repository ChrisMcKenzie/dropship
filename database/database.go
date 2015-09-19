package database

import (
	"log"

	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/jinzhu/gorm"
	_ "github.com/mxk/go-sqlite/sqlite3"
	"github.com/spf13/viper"
)

var db gorm.DB

func Init() {
	dbPath := viper.GetString("database.path")
	log.Printf("Opening database at %s", dbPath)
	var err error
	db, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	// Enable Logger
	db.LogMode(true)
	db.CreateTable(&model.User{})
	db.CreateTable(&model.Repo{})
	db.AutoMigrate(&model.User{}, &model.Repo{})
}

func Query() *gorm.DB {
	return &db
}
