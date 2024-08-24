package main

import (
	"database/sql"
	"m/tests/ClassDAO/repository"
	"testing"

	_ "github.com/lib/pq"
)

// Função para estabelecer conexão com o banco de dados (ajuste conforme necessário)
func setupDB() *sql.DB {
	psqlInfo := "host=localhost port=5432 user=my_user password=my@Pass%1234 dbname=my_database sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

// Benchmark para inserir uma classe.
func BenchmarkInsertClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		classID, err := repository.InsertClass(db, "Test Class")
		if err != nil {
			b.Fatalf("Failed to insert class: %v", err)
		}

		// Cleanup
		if err := repository.DeleteClass(db, classID); err != nil {
			b.Fatalf("Failed to clean up inserted class: %v", err)
		}
	}
}

// BenchmarkReadClass mede o desempenho do método ReadClass.
func BenchmarkReadClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	// Certifique-se de que a conexão com o banco de dados está funcionando.
	if err := db.Ping(); err != nil {
		b.Fatalf("Failed to ping database: %v", err)
	}

	// Preparar o ID da classe para o benchmark.
	classID := 1

	b.ResetTimer() // Inicia o timer do benchmark aqui, para não incluir o tempo de setup.

	for i := 0; i < b.N; i++ {
		_, err := repository.ReadClass(db, classID)
		if err != nil {
			b.Fatalf("Failed to read class: %v", err)
		}
	}
}

// Benchmark para atualizar uma classe.
func BenchmarkUpdateClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	classID, err := repository.InsertClass(db, "Initial Name")
	if err != nil {
		b.Fatalf("Failed to insert class: %v", err)
	}

	class, err := repository.ReadClass(db, classID)
	class.Name = "Updated Name"

	defer func() {
		// Cleanup
		if err := repository.DeleteClass(db, classID); err != nil {
			b.Fatalf("Failed to clean up inserted class: %v", err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repository.UpdateClass(db, class)
		if err != nil {
			b.Fatalf("Failed to update class: %v", err)
		}
	}
}

// Benchmark para deletar uma classe.
func BenchmarkDeleteClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		classID, err := repository.InsertClass(db, "Class to Delete")
		if err != nil {
			b.Fatalf("Failed to insert class: %v", err)
		}

		err = repository.DeleteClass(db, classID)
		if err != nil {
			b.Fatalf("Failed to delete class: %v", err)
		}
	}
}
