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

func NewDatabase(db *sql.DB) DAO {
	return DAO{Db: db}
}

// Create realiza insert considerando ID serial via banco de dados.
func (d DAO) Create(tableName string, entity interface{}) (int, error) {
	// Usar reflexão para iterar sobre os campos da entidade
	val := reflect.ValueOf(entity).Elem()
	typeOfT := val.Type()

	var fieldNames []string
	var fieldValues []interface{}
	var placeholders []string

	for i := 0; i < val.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" && dbTag != "id" { // Excluir o campo ID para autoincremento
			fieldNames = append(fieldNames, dbTag)
			fieldValues = append(fieldValues, val.Field(i).Interface())
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
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

// Read busca uma entidade pelo ID e preenche a struct passada com os dados encontrados.
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
			// Preparar os alvos de scan para os resultados da query
			scanTargets = append(scanTargets, val.Field(i).Addr().Interface())
		}
	}

	// Construir a string de query SQL
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(columnNames, ", "),
		tableName,
	)

	// Executar a query SQL
	row := d.Db.QueryRow(query, id)
	if err := row.Scan(scanTargets...); err != nil {
		return err
	}

	return nil
}

// Update atualiza qualquer struct no banco de dados.
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

	// Adiciona o ID ao final dos valores para a cláusula WHERE
	fieldValues = append(fieldValues, idFieldValue)

	// Construir a string de query SQL
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d",
		tableName,
		strings.Join(setClauses, ", "),
		idFieldName,
		len(fieldValues),
	)

	// Executar a query SQL
	_, err := d.Db.Exec(query, fieldValues...)
	if err != nil {
		return err
	}

	return nil
}

// Delete remove uma entidade pelo ID da tabela especificada.
func (d DAO) Delete(tableName string, id interface{}) error {
	// Construir a string de query SQL
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	// Executar a query SQL
	result, err := d.Db.Exec(query, id)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

// ReadMultiple busca múltiplas entidades com base em uma condição SQL e argumentos.
// A função aceita uma struct vazia como modelo para os resultados.
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
