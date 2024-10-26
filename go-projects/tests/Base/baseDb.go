package base

import (
	"database/sql"
	"fmt"
)

func SetupDB() *sql.DB {
	psqlInfo := "host=localhost port=5432 user=my_user password=my@Pass%1234 dbname=my_database sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func ClearAllProjectsAndResources(db *sql.DB) error {
	queries := []string{
		"DELETE FROM PROJECTS;",
		"DELETE FROM RESOURCES;",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("erro ao executar query '%s': %v", query, err)
		}
	}
	return nil
}
