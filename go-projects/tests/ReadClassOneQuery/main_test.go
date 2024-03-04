package main

import (
	"database/sql"
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
		_, err := ReadClass(db, classID)
		if err != nil {
			b.Fatalf("Failed to read class: %v", err)
		}
	}
}
