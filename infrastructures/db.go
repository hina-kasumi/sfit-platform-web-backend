package infrastructures

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

var myDB *gorm.DB

func OpenDbConnection(username string, password string, dbName string, host string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable ", host, username, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(-1)
	}

	myDB = db
	return db
}

func GetDB() *gorm.DB {
	return myDB
}
