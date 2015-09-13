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

	db.AutoMigrate(&model.User{}, &model.Repo{})
}
