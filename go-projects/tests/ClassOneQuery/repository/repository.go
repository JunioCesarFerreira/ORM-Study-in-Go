package repository

import (
	"database/sql"
	"m/tests/ClassOneQuery/entities"
	"time"
)

func InsertClass(db *sql.DB, name string) (int, error) {
	query := `
		INSERT INTO CLASSES (NAME)
		VALUES ($1)
		RETURNING ID
	`
	var classID int
	err := db.QueryRow(query, name).Scan(&classID)
	if err != nil {
		return 0, err
	}

	return classID, nil
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
		LEFT JOIN OBJECTS o ON c.ID = o.CLASS_ID
		LEFT JOIN OBJECT_ITEM_LINK l ON o.ID = l.OBJECT_ID
		LEFT JOIN ITEMS i ON i.ID = l.ITEM_ID
	WHERE c.ID = $1
	`

	rows, err := db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	class := &entities.Class{Id: 0}
	objectMap := make(map[int]*entities.Object)

	for rows.Next() {
		var cID, oID, iID sql.NullInt32
		var cName, oName, iName sql.NullString
		var oDateTime, iDateTime *time.Time
		var oValue, iValue sql.NullFloat64

		err := rows.Scan(&cID, &cName, &oID, &oName, &oDateTime, &oValue, &iID, &iName, &iDateTime, &iValue)
		if err != nil {
			return nil, err
		}

		if class.Id == 0 {
			class.Id = int(cID.Int32)
			class.Name = cName.String
		}

		if oID.Valid { // Check if the object exists
			object, ok := objectMap[int(oID.Int32)]
			if !ok {
				object = &entities.Object{Id: int(oID.Int32), Name: oName.String, DateTime: oDateTime, Value: oValue.Float64}
				objectMap[int(oID.Int32)] = object
			}

			if iID.Valid { // Check if the item exists
				item := entities.Item{Id: int(iID.Int32), Name: iName.String, DateTime: iDateTime, Value: iValue.Float64}
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

func UpdateClass(db *sql.DB, classID int, newName string) error {
	query := `
		UPDATE CLASSES
		SET NAME = $1
		WHERE ID = $2
	`
	_, err := db.Exec(query, newName, classID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteClass(db *sql.DB, classID int) error {
	query := `
		DELETE FROM CLASSES
		WHERE ID = $1
	`
	_, err := db.Exec(query, classID)
	if err != nil {
		return err
	}

	return nil
}
