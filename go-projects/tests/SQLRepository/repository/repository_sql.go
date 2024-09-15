package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

// Entity represents the general behavior for all entities
type Entity interface {
	TableName() string
	ColumnsNames() []string
	Fields() []interface{}
}

// Repository defines the interface for database operations specific to an entity
type Repository interface {
	Get(id int, entity Entity) (Entity, error)
	Add(entity Entity) error
	Update(entity Entity) error
	Delete(id int, entity Entity) error
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
	fields, values := repo.prepareFieldsAndValues(entity)
	placeholders := repo.generatePlaceholders(len(values))

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", entity.TableName(), fields, placeholders)

	id := -1
	err := repo.db.QueryRow(query, values...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *SQLRepository) Update(entity Entity) error {
	fields := repo.prepareFieldsForUpdate(entity)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $1", entity.TableName(), fields)
	_, err := repo.db.Exec(query, entity.Fields()...)
	return err
}

func (repo *SQLRepository) Delete(id int, entity Entity) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", entity.TableName())
	_, err := repo.db.Exec(query, id)
	return err
}

func (repo *SQLRepository) prepareFieldsAndValues(entity Entity) (string, []interface{}) {
	columns := entity.ColumnsNames()
	fields := entity.Fields()
	var cols []string
	var values []interface{}

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		if col != "id" {
			cols = append(cols, col)
			values = append(values, fields[i])
		}
	}
	return strings.Join(cols, ", "), values
}

func (repo *SQLRepository) prepareFieldsForUpdate(entity Entity) string {
	columns := entity.ColumnsNames()
	var updatePairs []string

	for i := 0; i < len(columns); i++ {
		field := columns[i]
		if field != "id" { // Assuming 'ID' is not included in UPDATE queries
			updatePairs = append(updatePairs, fmt.Sprintf("%s = $%d", field, i+1))
		}
	}
	return strings.Join(updatePairs, ", ")
}

func (repo *SQLRepository) generatePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}
