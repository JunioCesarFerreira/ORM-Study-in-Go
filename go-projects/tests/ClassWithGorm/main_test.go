package main

import (
	"fmt"
	"log"
	"m/tests/ClassWithGorm/repository"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Função para estabelecer conexão com o banco de dados (ajuste conforme necessário)
func setupDB() *gorm.DB {
	psqlInfo := "host=localhost port=5432 user=my_user password=my@Pass%1234 dbname=my_database sslmode=disable"
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func BenchmarkInsertClass(b *testing.B) {
	db := setupDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		class, err := repository.InsertClass(db, fmt.Sprintf("Test Class %d", i))
		if err != nil {
			b.Fatalf("Failed to insert class: %v", err)
		}

		// Cleanup
		if err := repository.DeleteClass(db, class.ID); err != nil {
			b.Fatalf("Failed to clean up inserted class: %v", err)
		}
	}
}

// BenchmarkReadClass mede o desempenho do método ReadClass.
func BenchmarkReadClass(b *testing.B) {
	db := setupDB()

	class, err := repository.InsertClass(db, "Initial Name")
	if err != nil {
		b.Fatalf("Failed to insert class: %v", err)
	}

	defer func() {
		// Cleanup
		if err := repository.DeleteClass(db, class.ID); err != nil {
			b.Fatalf("Failed to clean up inserted class: %v", err)
		}
	}()

	b.ResetTimer() // Inicia o timer do benchmark aqui, para não incluir o tempo de setup.

	for i := 0; i < b.N; i++ {
		classRead, err := repository.ReadClass(db, class.ID)
		if err != nil {
			b.Fatalf("Failed to read class: %v", err)
		}
		if class.Name != classRead.Name {
			b.Fatalf("Class name mismatch: expected %s, got %s", class.Name, classRead.Name)
		}
	}
}

func BenchmarkUpdateClass(b *testing.B) {
	db := setupDB()

	class, err := repository.InsertClass(db, "Initial Name")
	if err != nil {
		b.Fatalf("Failed to insert class: %v", err)
	}

	defer func() {
		// Cleanup
		if err := repository.DeleteClass(db, class.ID); err != nil {
			b.Fatalf("Failed to clean up inserted class: %v", err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repository.UpdateClass(db, class.ID, fmt.Sprintf("Updated Name %d", i))
		if err != nil {
			b.Fatalf("Failed to update class: %v", err)
		}
	}
}

func BenchmarkDeleteClass(b *testing.B) {
	db := setupDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		class, err := repository.InsertClass(db, "Class to Delete")
		if err != nil {
			b.Fatalf("Failed to insert class: %v", err)
		}

		err = repository.DeleteClass(db, class.ID)
		if err != nil {
			b.Fatalf("Failed to delete class: %v", err)
		}
	}
}
