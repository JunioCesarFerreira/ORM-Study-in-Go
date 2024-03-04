package main

import (
	"log"
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

// BenchmarkReadClass mede o desempenho do método ReadClass.
func BenchmarkReadClass(b *testing.B) {
	db := setupDB()

	// Preparar o ID da classe para o benchmark.
	classID := 1

	b.ResetTimer() // Inicia o timer do benchmark aqui, para não incluir o tempo de setup.

	for i := 0; i < b.N; i++ {
		_, err := ReadClass(db, classID)
		if err != nil {
			b.Fatalf("Failed to read class: %v", err)
		}
	}
}
