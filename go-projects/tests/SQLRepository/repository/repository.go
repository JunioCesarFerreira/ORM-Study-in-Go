package repository

import (
	"database/sql"
	"m/tests/SQLRepository/entities"
)

// InsertClass insere uma nova classe no banco de dados
func InsertClass(db *sql.DB, name string) (int, error) {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}

	class := &entities.Class{
		Name: name,
	}

	// Adiciona a classe usando o método Add do SQLRepository
	id, err := repo.Add(class)
	if err != nil {
		return 0, err
	}
	class.Id = id

	// Retorna o ID da nova classe inserida
	return class.Id, nil
}

// ReadClass lê uma classe do banco de dados por ID
func ReadClass(db *sql.DB, classID int) (*entities.Class, error) {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}

	class := &entities.Class{}

	// Busca a classe pelo ID usando o método Get do SQLRepository
	err = repo.Get(classID, class)
	if err != nil {
		return nil, err
	}

	// Converte para o tipo Class e retorna
	return class, nil
}

// UpdateClass atualiza o nome de uma classe existente
func UpdateClass(db *sql.DB, classID int, newName string) error {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}

	// Usamos o método Update do SQLRepository para atualizar no banco de dados
	err = repo.Update(&entities.Class{Id: classID, Name: newName})
	if err != nil {
		return err
	}

	return nil
}

// DeleteClass exclui uma classe do banco de dados por ID
func DeleteClass(db *sql.DB, classID int) error {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}

	class := &entities.Class{}

	// Usamos o método Delete do SQLRepository para remover a classe
	err = repo.Delete(classID, class)
	if err != nil {
		return err
	}

	return nil
}
