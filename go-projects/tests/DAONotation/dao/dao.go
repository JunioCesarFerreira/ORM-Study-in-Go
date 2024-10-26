package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type DAO struct {
	Db *sql.DB
}

func NewDAO(db *sql.DB) DAO {
	return DAO{Db: db}
}

// Create performs insert considering auto-increment ID via database.
func (d DAO) Create(tableName string, entity interface{}) (int, error) {
	// Use reflection to iterate over entity fields
	val := reflect.ValueOf(entity).Elem()
	typeOfT := val.Type()

	var fieldNames []string
	var fieldValues []interface{}
	var placeholders []string

	count := 1
	for i := 0; i < val.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" {
			fieldNames = append(fieldNames, dbTag)
			fieldValues = append(fieldValues, val.Field(i).Interface())
			placeholders = append(placeholders, fmt.Sprintf("$%d", count))
			count++
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING ID",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	var newID int
	err := d.Db.QueryRow(query, fieldValues...).Scan(&newID)
	if err != nil {
		return -1, err
	}

	if idField := val.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(newID))
	}

	return newID, nil
}

func (d DAO) CreateChild(tableName string, entity interface{}, foreignKey string, foreignKeyValue int) (int, error) {
	// Use reflection to iterate over entity fields
	val := reflect.ValueOf(entity).Elem()
	typeOfT := val.Type()

	var fieldNames []string
	var fieldValues []interface{}
	var placeholders []string

	count := 1
	for i := 0; i < val.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" {
			fieldNames = append(fieldNames, dbTag)
			fieldValues = append(fieldValues, val.Field(i).Interface())
			placeholders = append(placeholders, fmt.Sprintf("$%d", count))
			count++
		}
	}

	// Add the foreign key to the list of fields and values
	fieldNames = append(fieldNames, foreignKey)
	fieldValues = append(fieldValues, foreignKeyValue)
	placeholders = append(placeholders, fmt.Sprintf("$%d", count))

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING ID",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	var newID int
	err := d.Db.QueryRow(query, fieldValues...).Scan(&newID)
	if err != nil {
		return -1, err
	}

	if idField := val.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(newID))
	}

	return newID, nil
}

func (d DAO) CreateWithLinkSingleSide(existingParentId int, childTable string, linkTable string, childId int, parentForeignKey string, childForeignKey string) (int, error) {
	// Insert into the link table (e.g., OBJECT_ITEM_LINK) using the existing parent object ID
	linkQuery := fmt.Sprintf(
		"INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		linkTable,
		parentForeignKey,
		childForeignKey,
	)

	_, err := d.Db.Exec(linkQuery, existingParentId, childId)
	if err != nil {
		return childId, err
	}

	return childId, nil
}

// Read fetches an entity by ID and fills the passed struct with the found data.
func (d DAO) Read(tableName string, id interface{}, entity interface{}) error {
	val := reflect.ValueOf(entity).Elem()
	typeOfEntity := val.Type()

	var columnNames []string
	var scanTargets []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typeOfEntity.Field(i)
		tag := field.Tag.Get("db")

		if tag != "" {
			columnNames = append(columnNames, tag)
			// Prepare scan targets for query results
			scanTargets = append(scanTargets, val.Field(i).Addr().Interface())
		}
	}

	// Build SQL query string
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(columnNames, ", "),
		tableName,
	)

	// Execute SQL query
	row := d.Db.QueryRow(query, id)
	if err := row.Scan(scanTargets...); err != nil {
		return err
	}

	return nil
}

// Update updates any struct in the database.
func (d DAO) Update(tableName string, entity interface{}) error {
	val := reflect.ValueOf(entity).Elem()
	typeOfEntity := val.Type()

	var setClauses []string
	var fieldValues []interface{}
	var idFieldValue interface{}
	var idFieldName string

	for i := 0; i < val.NumField(); i++ {
		field := typeOfEntity.Field(i)
		tag := field.Tag.Get("db")
		fieldValue := val.Field(i).Interface()

		if tag == "ID" {
			idFieldValue = fieldValue
			idFieldName = tag
			continue
		}

		if tag != "" {
			setClause := fmt.Sprintf("%s = $%d", tag, len(fieldValues)+1)
			setClauses = append(setClauses, setClause)
			fieldValues = append(fieldValues, fieldValue)
		}
	}

	// Add ID to the end of values for WHERE clause
	fieldValues = append(fieldValues, idFieldValue)

	// Build SQL query string
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d",
		tableName,
		strings.Join(setClauses, ", "),
		idFieldName,
		len(fieldValues),
	)

	// Execute SQL query
	_, err := d.Db.Exec(query, fieldValues...)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes an entity by ID from the specified table.
func (d DAO) Delete(tableName string, id interface{}) error {
	// Build SQL query string
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	// Execute SQL query
	result, err := d.Db.Exec(query, id)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

// ReadMultiple fetches multiple entities based on an SQL condition and arguments.
// The function accepts an empty struct as a model for the results.
func (d DAO) ReadMultiple(tableName string, condition string, args []interface{}, model interface{}) ([]interface{}, error) {
	sliceType := reflect.SliceOf(reflect.TypeOf(model))
	resultsSlice := reflect.MakeSlice(sliceType, 0, 0)

	val := reflect.New(reflect.TypeOf(model).Elem()).Elem()
	typeOfModel := val.Type()

	var columnNames []string
	var scanTargets []reflect.Value

	for i := 0; i < val.NumField(); i++ {
		field := typeOfModel.Field(i)
		tag := field.Tag.Get("db")

		if tag != "" {
			columnNames = append(columnNames, tag)
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columnNames, ", "),
		tableName,
		condition,
	)

	rows, err := d.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		newElem := reflect.New(reflect.TypeOf(model).Elem()).Elem()
		scanTargets = make([]reflect.Value, len(columnNames))

		for i := range columnNames {
			scanTargets[i] = newElem.Field(i).Addr()
		}

		scanArgs := make([]interface{}, len(scanTargets))
		for i, v := range scanTargets {
			scanArgs[i] = v.Interface()
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		resultsSlice = reflect.Append(resultsSlice, newElem.Addr())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	results := make([]interface{}, resultsSlice.Len())
	for i := 0; i < resultsSlice.Len(); i++ {
		results[i] = resultsSlice.Index(i).Interface()
	}

	return results, nil
}
