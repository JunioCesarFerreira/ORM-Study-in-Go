package main

import (
	"fmt"
	"log"
	"m/tests"
	"m/tests/ReadClassWithGorm/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "my_user"
	password = "my@Pass%1234"
	dbname   = "my_database"
)

func ReadClass(db *gorm.DB, classID int) (*entities.Class, error) {
	var class entities.Class
	err := db.Preload("Objects.Items").First(&class, classID).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	classID := 1 // Use the actual class ID you want to read
	class, err := ReadClass(db, classID)
	if err != nil {
		log.Fatalf("failed to read class: %v", err)
	}

	tests.SaveResult("WithGorm.json", class)
}
