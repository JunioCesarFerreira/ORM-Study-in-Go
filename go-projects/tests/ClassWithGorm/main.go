package main

import (
	"fmt"
	"log"
	"m/tests"
	"m/tests/ClassWithGorm/entities"

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

func InsertClass(db *gorm.DB, name string) (*entities.Class, error) {
	class := entities.Class{Name: name}
	if err := db.Create(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func ReadClass(db *gorm.DB, classID int) (*entities.Class, error) {
	var class entities.Class
	err := db.Preload("Objects.Items").First(&class, classID).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func UpdateClass(db *gorm.DB, classID int, newName string) error {
	return db.Model(&entities.Class{}).Where("id = ?", classID).Update("name", newName).Error
}

func DeleteClass(db *gorm.DB, classID int) error {
	return db.Delete(&entities.Class{}, classID).Error
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
