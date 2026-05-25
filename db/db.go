package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(
		sqlite.Open("db.db"),
		&gorm.Config{},
	)

	if err != nil {
		panic("❌ Failed to connect database.")
	}

	fmt.Println("✅ Database connected.")

	DB = db
}
