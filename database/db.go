package database

import (
	"L0-wb/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Connect() {
	connStr := "host=localhost port=5432 user=admin password=admin dbname=l0_database sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to open DB conn: ", err)
	}

	err = db.AutoMigrate(
		&models.Order{},
		&models.Payment{},
		&models.Item{},
		&models.Delivery{},
	)
	if err != nil {
		log.Fatal("failed to migrate models to DB: ", err)
	}

	log.Println("DB connection established")
}
