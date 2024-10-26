package main

import (
	"database/sql"
	"fmt"
	"log"
	"m/tests"
	base "m/tests/Base"
	"m/tests/DAONotation/entities"
	"m/tests/DAONotation/repository"

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

	base.ClearAllProjectsAndResources(db)

	data, err := base.OpenInputData()
	if err != nil {
		log.Fatal(err)
	}

	resources, err := base.Cast[[]entities.Resource](data.Resources)
	if err != nil {
		log.Fatalf("Failed to cast resources: %v", err)
	}

	for _, resource := range resources {
		_, err := repository.InsertResource(db, resource)
		if err != nil {
			log.Fatalf("Failed to insert resource: %v", err)
		}
	}

	firstProject, err := base.Cast[entities.Project](data.Projects[0])

	projectId, err := repository.InsertProject(db, firstProject)
	if err != nil {
		log.Fatal(err)
	}

	project, err := repository.ReadProject(db, projectId)
	if err != nil {
		log.Fatal(err)
	}

	tests.SaveResult("result_main_execution.json", project)
}
