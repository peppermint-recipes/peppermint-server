package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/theErikss0n/peppermint-server/services/recipe-service/config"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Init(dbConnection *config.DBConfig) error {
	connString := fmt.Sprintf("%s:%s@(%s:%d)/%s?%s&%s&%s",
		dbConnection.Username,
		dbConnection.Password,
		dbConnection.Host,
		dbConnection.Port,
		dbConnection.Name,
		"charset=utf8mb4",
		"parseTime=true",
		"loc=Local",
	)

	dbCon, err := gorm.Open(dbConnection.Dialect, connString)
	if err != nil {
		log.Println("Connection Failed to Open")

		return err
	} else {
		log.Println("Connection Established")
	}

	dbCon.AutoMigrate(&Recipe{})

	db = dbCon

	return nil
}
