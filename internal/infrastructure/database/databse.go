package database

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDbConnection(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable ", cfg.Host, cfg.User, cfg.Password, cfg.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(-1)
	}

	return db
}

func CloseDbConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("db close err: ", err)
		return
	}
	sqlDB.Close()
}
