package main

import (
	"database/sql"
	"encoding/json"
	"m/tests"
	"m/tests/ClassDAO/entities"
	"m/tests/ClassDAO/repository"
	"testing"
	"time"

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

func getInputData(b *testing.B) []entities.Class {

	byteValue, err := tests.OpenInputJson()
	if err != nil {
		b.Fatalf("Error reading test JSON file: %s", err)
	}

	var data []entities.Class
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		b.Fatalf("Error deserializing test JSON file: %s", err)
	}

	return data
}

func compareClasses(b *testing.B, c1, c2 entities.Class) bool {
	if c1.Id != c2.Id || c1.Name != c2.Name {
		b.Fatalf("Class name mismatch: expected %s, got %s", c1.Name, c2.Name)
		return false
	}

	if len(c1.Objects) != len(c2.Objects) {
		b.Fatalf("Object length mismatch: expected %d, got %d", len(c1.Objects), len(c2.Objects))
		return false
	}

	for i := range c1.Objects {
		obj1 := c1.Objects[i]
		obj2 := c2.Objects[i]

		if obj1.Name != obj2.Name || obj1.Value != obj2.Value {
			b.Fatalf("Object mismatch: expected %v, got %v", c1.Objects, c2.Objects)
			return false
		}

		if (obj1.DateTime != nil && obj2.DateTime != nil) && !obj1.DateTime.Equal(*obj2.DateTime) {
			b.Fatalf("Object mismatch: expected %v, got %v", c1.Objects, c2.Objects)
			return false
		}

		if len(obj1.Items) != len(obj2.Items) {
			b.Fatalf("Object mismatch: expected %v, got %v", c1.Objects, c2.Objects)
			return false
		}

		for j := range obj1.Items {
			item1 := obj1.Items[j]
			item2 := obj2.Items[j]

			if item1.Name != item2.Name || item1.Value != item2.Value {
				b.Fatalf("Item mismatch: expected %v, got %v", item1, item2)
				return false
			}

			if (item1.DateTime != nil && item2.DateTime != nil) && !item1.DateTime.Equal(*item2.DateTime) {
				b.Fatalf("Item mismatch: expected %v, got %v", item1, item2)
				return false
			}
		}
	}

	return true
}

// Benchmark para inserir uma classe.
func BenchmarkInsertClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	err := tests.ClearDatabase(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	data := getInputData(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, class := range data {
			_, err := repository.InsertClass(db, class)
			if err != nil {
				b.Fatalf("Failed to insert class: %v", err)
			}
		}
	}
}

// BenchmarkReadClass mede o desempenho do método ReadClass.
func BenchmarkReadClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	err := tests.ClearDatabase(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	data := getInputData(b)

	for i := 0; i < b.N; i++ {
		for j, class := range data {
			classID, err := repository.InsertClass(db, class)
			if err != nil {
				b.Fatalf("Failed to insert class: %v", err)
			}
			data[j].Id = classID
		}
	}

	b.ResetTimer() // Inicia o timer do benchmark aqui, para não incluir o tempo de setup.

	for i := 0; i < b.N; i++ {
		for _, class := range data {
			readClass, err := repository.ReadClass(db, class.Id)
			if err != nil {
				b.Fatalf("Failed to read class: %v", err)
			}

			compareClasses(b, class, *readClass)
		}
	}
}

// Benchmark para atualizar uma classe.
func BenchmarkUpdateClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	err := tests.ClearDatabase(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	data := getInputData(b)

	for i := 0; i < b.N; i++ {
		for j, class := range data {
			classID, err := repository.InsertClass(db, class)
			if err != nil {
				b.Fatalf("Failed to insert class: %v", err)
			}
			data[j].Id = classID
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, class := range data {
			change := class
			change.Name = "new name"
			if len(change.Objects) > 0 {
				*change.Objects[0].DateTime = time.Now()
				if len(change.Objects[0].Items) > 0 {
					*&change.Objects[0].Items[0].Value = 3.14
				}
			}
			err := repository.UpdateClass(db, &change)
			if err != nil {
				b.Fatalf("Failed to update class: %v", err)
			}
		}
	}
}

// Benchmark para deletar uma classe.
func BenchmarkDeleteClass(b *testing.B) {
	db := setupDB()
	defer db.Close()

	err := tests.ClearDatabase(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	data := getInputData(b)

	for j, class := range data {
		classID, err := repository.InsertClass(db, class)
		if err != nil {
			b.Fatalf("Failed to insert class: %v", err)
		}
		data[j].Id = classID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, class := range data {
			err = repository.DeleteClass(db, class.Id)
			if err != nil {
				b.Fatalf("Failed to delete class: %v", err)
			}
		}
	}
}
