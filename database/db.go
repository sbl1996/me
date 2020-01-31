package database

import (
	"os"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	gorm.Open("postgres", dbURL)
	db, err := gorm.Open("postgres", dbURL)
	DB = db
	return db, err
}
