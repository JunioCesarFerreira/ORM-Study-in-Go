package main

import (
	"database/sql"
	"fmt"
	"log"
	"m/tests"
	"m/tests/ReadClassOneQuery/entities"
	"time"

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
	query := `
	SELECT c.ID, c.NAME, o.ID, o.NAME, o.DATETIME, o.VALUE, i.ID, i.NAME, i.DATETIME, i.VALUE
	FROM CLASSES c
	INNER JOIN OBJECTS o ON c.ID = o.CLASS_ID
	INNER JOIN OBJECT_ITEM_LINK l ON o.ID = l.OBJECT_ID
	INNER JOIN ITEMS i ON i.ID = l.ITEM_ID
	WHERE c.ID = $1
	`

	rows, err := db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	class := &entities.Class{}
	objectMap := make(map[int]*entities.Object)

	for rows.Next() {
		var cID, oID, iID int
		var cName, oName, iName string
		var oDateTime, iDateTime *time.Time
		var oValue, iValue float64

		err := rows.Scan(&cID, &cName, &oID, &oName, &oDateTime, &oValue, &iID, &iName, &iDateTime, &iValue)
		if err != nil {
			return nil, err
		}

		if class.Id == 0 {
			class.Id = cID
			class.Name = cName
		}

		if oID != 0 { // Check if the object exists
			object, ok := objectMap[oID]
			if !ok {
				object = &entities.Object{Id: oID, Name: oName, DateTime: oDateTime, Value: oValue}
				objectMap[oID] = object
			}

			if iID != 0 { // Check if the item exists
				item := entities.Item{Id: iID, Name: iName, DateTime: iDateTime, Value: iValue}
				object.Items = append(object.Items, item)
			}
		}
	}

	for _, obj := range objectMap {
		class.Objects = append(class.Objects, *obj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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

	tests.SaveResult("OneQuery.json", class)
}
