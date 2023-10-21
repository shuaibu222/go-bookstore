package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var myDb *gorm.DB

func Connect() {

	config, err := LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	// dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error while connecting to the database: ", err)
	}
	myDb = db
}

func GetDb() *gorm.DB {
	return myDb
}
