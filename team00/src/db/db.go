package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Uuid      string
	Frequency float64
	Timestamp int64
}

func Connect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to psql: %v", err)
	}

	if err = db.AutoMigrate(&Record{}); err != nil {
		return nil, err
	}
	return db, nil
}
