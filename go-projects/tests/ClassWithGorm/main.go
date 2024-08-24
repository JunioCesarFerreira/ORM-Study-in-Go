package main

import (
	"fmt"
	"log"
	"m/tests"
	"m/tests/ClassWithGorm/repository"

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

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	classID := 1 // Use the actual class ID you want to read
	class, err := repository.ReadClass(db, classID)
	if err != nil {
		log.Fatalf("failed to read class: %v", err)
	}

	tests.SaveResult("WithGorm.json", class)
}
