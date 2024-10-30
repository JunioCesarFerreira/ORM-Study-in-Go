package repository

import (
	"database/sql"
	columnfieldmap "m/tests/SQLRepository/columnFieldMap"
	"reflect"
	"strconv"
	"strings"
)

// Entity represents the general behavior for all entities
type Entity interface {
	TableName() string
	ColumnsNames() []string
	Fields() []interface{}
	PKColNames() []string
	PKFields() []interface{}
}

// SQLRepository is the concrete implementation of Repository for SQL databases
type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) (*SQLRepository, error) {
	return &SQLRepository{db: db}, nil
}

func (repo *SQLRepository) Get(id int, entity Entity) error {
	fields := strings.Join(entity.ColumnsNames(), ", ")
	query := "SELECT " + fields + " FROM " + entity.TableName() + " WHERE id = $1"
	row := repo.db.QueryRow(query, id)
	err := row.Scan(entity.Fields()...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SQLRepository) Add(entity Entity) (int, error) {
	cols, values := repo.prepareFieldsAndValuesForAdd(entity)
	placeholders := repo.generatePlaceholders(len(values))

	query := "INSERT INTO " + entity.TableName() + " (" + cols + ") VALUES (" + placeholders + ") RETURNING IDENTIFIER"

	id := -1
	err := repo.db.QueryRow(query, values...).Scan(&id)
	return id, err
}

func (repo *SQLRepository) Insert(entity Entity) error {
	cols, values := repo.prepareFieldsAndValuesForInsert(entity)
	placeholders := repo.generatePlaceholders(len(values))

	query := "INSERT INTO " + entity.TableName() + " (" + cols + ") VALUES (" + placeholders + ")"

	_, err := repo.db.Exec(query, values...)
	return err
}

func (repo *SQLRepository) InsertWithFK(entity Entity, fks []columnfieldmap.ColumnFieldPair) error {
	cols, values := repo.prepareFieldsAndValuesForInsert(entity)

	for _, fk := range fks {
		cols += ", " + fk.ColumnName
		values = append(values, repo.recValue(fk.Field))
	}

	placeholders := repo.generatePlaceholders(len(values))

	query := "INSERT INTO " + entity.TableName() + " (" + cols + ") VALUES (" + placeholders + ")"

	_, err := repo.db.Exec(query, values...)
	return err
}

func (repo *SQLRepository) Update(entity Entity) error {
	fields, values := repo.prepareFieldsAndValuesForUpdate(entity)
	conditional, condValues := repo.buildConditional(entity, len(values)+1)

	query := "UPDATE " + entity.TableName() + " SET " + fields + " WHERE " + conditional

	_, err := repo.db.Exec(query, append(values, condValues...)...)
	return err
}

func (repo *SQLRepository) Delete(id int, entity Entity) error {
	query := "DELETE FROM " + entity.TableName() + " WHERE id = $1"
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *SQLRepository) Links(links Links) error {
	del := "DELETE FROM " + links.TableName + " WHERE " + links.MasterColName + " = $1"

	_, err := repo.db.Exec(del, links.MasterId)
	if err != nil {
		return err
	}

	for _, link := range links.LinksIds {
		insert := "INSERT INTO " + links.TableName + " (" + links.MasterColName + ", " + links.LinkColName + ") VALUES ($1, $2)"

		_, err := repo.db.Exec(insert, links.MasterId, link)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *SQLRepository) prepareFieldsAndValuesForAdd(entity Entity) (string, []interface{}) {
	columns := entity.ColumnsNames()
	fields := entity.Fields()
	pkCols := entity.PKColNames()
	var cols []string
	var values []interface{}

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		if !repo.sliceContainsFold(pkCols, col) {
			cols = append(cols, col)
			values = append(values, repo.recValue(fields[i]))
		}
	}
	return strings.Join(cols, ", "), values
}

func (repo *SQLRepository) prepareFieldsAndValuesForInsert(entity Entity) (string, []interface{}) {
	columns := entity.ColumnsNames()
	fields := entity.Fields()
	var cols []string
	var values []interface{}

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		cols = append(cols, col)
		values = append(values, repo.recValue(fields[i]))
	}
	return strings.Join(cols, ", "), values
}

func (repo *SQLRepository) prepareFieldsAndValuesForUpdate(entity Entity) (string, []interface{}) {
	columns := entity.ColumnsNames()
	fields := entity.Fields()
	pkCols := entity.PKColNames()
	var updatePairs []string
	var values []interface{}

	count := 0
	for i := 0; i < len(columns); i++ {
		field := columns[i]
		if !repo.sliceContainsFold(pkCols, columns[i]) {
			count++
			updatePairs = append(updatePairs, field+" = $"+strconv.Itoa(count))
			values = append(values, repo.recValue(fields[i]))
		}
	}
	return strings.Join(updatePairs, ", "), values
}

func (repo *SQLRepository) generatePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(placeholders, ", ")
}

func (repo *SQLRepository) recValue(input any) any {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value.Interface()
}

func (repo *SQLRepository) sliceContainsFold(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

func (repo *SQLRepository) buildConditional(entity Entity, firstPlaceholder int) (string, []interface{}) {
	cols := entity.PKColNames()
	fields := entity.PKFields()
	if len(cols) != len(fields) {
		panic("Invalid input: the number of columns and fields must be the same.")
	}
	var comps []string
	var values []interface{}
	for i := 0; i < len(cols); i++ {
		comp := cols[i] + " = $" + strconv.Itoa(i+firstPlaceholder)
		comps = append(comps, comp)
		values = append(values, repo.recValue(fields[i]))
	}
	conditional := strings.Join(comps, " AND ")
	return conditional, values
}
