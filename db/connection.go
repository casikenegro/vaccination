package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConfigDB struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

var DB *gorm.DB

func DBConnection(configDB *ConfigDB) {
	var error error
	var DSN = "host=" + configDB.Host + " user=" + configDB.User + " password=" + configDB.Password + " dbname=" + configDB.DBName + " port=" + configDB.Port

	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("Database connection successful")
	}
}
