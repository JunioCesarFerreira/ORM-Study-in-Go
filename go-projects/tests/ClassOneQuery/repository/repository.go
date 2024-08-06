package repository

import (
	"database/sql"
	"m/tests/ClassOneQuery/entities"
	"time"
)

func InsertClass(db *sql.DB, name string) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO CLASSES (NAME)
		VALUES ($1)
		RETURNING ID
	`
	var classID int
	err = tx.QueryRow(query, name).Scan(&classID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return classID, nil
}

func UpdateClass(db *sql.DB, classID int, newName string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE CLASSES
		SET NAME = $1
		WHERE ID = $2
	`
	_, err = tx.Exec(query, newName, classID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ReadClass(db *sql.DB, classID int) (*entities.Class, error) {
	query := `
	SELECT 
		c.ID, 
		c.NAME, 
		o.ID, 
		o.NAME, 
		o.DATETIME, 
		o.VALUE, 
		i.ID, 
		i.NAME, 
		i.DATETIME, 
		i.VALUE
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

func DeleteClass(db *sql.DB, classID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
		DELETE FROM CLASSES
		WHERE ID = $1
	`
	_, err = tx.Exec(query, classID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
