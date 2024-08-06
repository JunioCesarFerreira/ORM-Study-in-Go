package main

import (
	"database/sql"
	"fmt"
	"log"
	"m/tests"
	"m/tests/ClassDAOcrud/database"
	"m/tests/ClassDAOcrud/entities"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "my_user"
	password = "my@Pass%1234"
	dbname   = "my_database"
)

func ReadClass(db *sql.DB, classID int) (*entities.Class, error) {
	class := &entities.Class{}

	crud := database.NewDatabase(db)

	err := crud.Read("CLASSES", classID, class)
	if err != nil {
		return class, err
	}

	condition := "CLASS_ID = $1"
	args := []interface{}{classID}
	model := &entities.Object{}

	objects, err := crud.ReadMultiple("OBJECTS", condition, args, model)
	if err != nil {
		return class, err
	}

	for _, obj := range objects {
		tmpObj := *obj.(*entities.Object)
		condition := "OBJECT_ID = $1"
		args := []interface{}{tmpObj.Id}
		model := &entities.Item{}

		items, err := crud.ReadMultiple("ITEMS_BY_OBJECT_VIEW", condition, args, model)
		if err != nil {
			return class, err
		}

		for _, item := range items {
			tmpObj.Items = append(tmpObj.Items, *item.(*entities.Item))
		}

		class.Objects = append(class.Objects, tmpObj)
	}

	return class, nil
}

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
	class, err := ReadClass(db, classID)
	if err != nil {
		log.Fatal(err)
	}

	tests.SaveResult("DAOcrud.json", class)
}
