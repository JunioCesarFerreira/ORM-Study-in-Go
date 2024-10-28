package repository

import (
	"database/sql"
	"fmt"
	columnfieldmap "m/tests/SQLRepository/columnFieldMap"
	"reflect"
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
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", fields, entity.TableName())
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

	query := "INSERT INTO %s (%s, TENANT_ID) VALUES (%s) RETURNING IDENTIFIER"
	query = fmt.Sprintf(query, entity.TableName(), cols, placeholders)

	id := -1
	err := repo.db.QueryRow(query, values...).Scan(&id)
	return id, err
}

func (repo *SQLRepository) Insert(entity Entity) error {
	cols, values := repo.prepareFieldsAndValuesForInsert(entity)
	placeholders := repo.generatePlaceholders(len(values))

	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, entity.TableName(), cols, placeholders)

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

	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, entity.TableName(), cols, placeholders)

	_, err := repo.db.Exec(query, values...)
	return err
}

func (repo *SQLRepository) Update(entity Entity) error {
	fields, values := repo.prepareFieldsAndValuesForUpdate(entity)
	conditional, condValues := repo.buildConditional(entity, len(values)+1)

	query := "UPDATE %s SET %s WHERE %s"
	query = fmt.Sprintf(query, entity.TableName(), fields, conditional)

	_, err := repo.db.Exec(query, append(values, condValues...)...)
	return err
}

func (repo *SQLRepository) Delete(id int, entity Entity) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", entity.TableName())
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *SQLRepository) Links(links Links) error {
	del := "DELETE FROM %s WHERE %s = $1"
	del = fmt.Sprintf(del, links.TableName, links.MasterColName)

	_, err := repo.db.Exec(del, links.MasterId)
	if err != nil {
		return err
	}

	for _, link := range links.LinksIds {
		insert := "INSERT INTO %s (%s, %s) VALUES ($1, $2)"
		insert = fmt.Sprintf(insert, links.TableName, links.MasterColName, links.LinkColName)

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
			updatePairs = append(updatePairs, fmt.Sprintf("%s = $%d", field, count))
			values = append(values, repo.recValue(fields[i]))
		}
	}
	return strings.Join(updatePairs, ", "), values
}

func (repo *SQLRepository) generatePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
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
		comp := fmt.Sprintf("%s = $%d", cols[i], i+firstPlaceholder)
		comps = append(comps, comp)
		values = append(values, repo.recValue(fields[i]))
	}
	conditional := strings.Join(comps, " AND ")
	return conditional, values
}
