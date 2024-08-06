package main

import (
	"database/sql"
	"fmt"
	"log"
	"m/tests"
	"m/tests/ClassOneQuery/repository"

	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database.")

	classID := 1 // O ID da classe que queremos buscar
	class, err := repository.ReadClass(db, classID)
	if err != nil {
		log.Fatal(err)
	}

	tests.SaveResult("OneQuery.json", class)
}
