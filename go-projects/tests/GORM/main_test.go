package main

import (
	"log"
	base "m/tests/Base"
	"m/tests/GORM/entities"
	"m/tests/GORM/repository"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func startupTest(b *testing.B) (*gorm.DB, []entities.Resource, []entities.Project) {
	psqlInfo := "host=localhost port=5432 user=my_user password=my@Pass%1234 dbname=my_database sslmode=disable"
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	data := base.GetInputData(b)

	resources, err := base.Cast[[]entities.Resource](data.Resources)
	if err != nil {
		b.Fatalf("Failed to cast resources: %v", err)
	}

	projects, err := base.Cast[[]entities.Project](data.Projects)
	if err != nil {
		b.Fatalf("Failed to cast projects: %v", err)
	}

	return db, resources, projects
}

// Benchmark for inserting a resources.
func BenchmarkInsertResources(b *testing.B) {
	dbg, resources, _ := startupTest(b)

	db := base.SetupDB()
	defer db.Close()

	err := base.ClearAllProjectsAndResources(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	b.ResetTimer() // Start benchmark timer here to exclude setup time.

	for i := 0; i < b.N; i++ {
		for _, resource := range resources {
			_, err := repository.InsertResource(dbg, resource)
			if err != nil {
				b.Fatalf("Failed to insert resource: %v", err)
			}
		}
	}
}

// Benchmark for inserting a project.
func BenchmarkInsertProject(b *testing.B) {
	db, _, projects := startupTest(b)

	b.ResetTimer() // Start benchmark timer here to exclude setup time.

	for i := 0; i < b.N; i++ {
		for _, project := range projects {
			_, err := repository.InsertProject(db, project)
			if err != nil {
				b.Fatalf("Failed to insert project: %v", err)
			}
		}
	}
}

// BenchmarkReadProject measures the performance of the ReadProject method.
func BenchmarkReadProject(b *testing.B) {
	db, _, projects := startupTest(b)

	b.ResetTimer() // Start benchmark timer here to exclude setup time.

	for i := 0; i < b.N; i++ {
		for _, project := range projects {
			readProject, err := repository.ReadProject(db, project.ID)
			if err != nil {
				b.Fatalf("Failed to read project: %v", err)
			}

			if base.CompareObjectsAsJSON(project, *readProject) != nil {
				b.Errorf("Objects do not match.")
			}
		}
	}
}

// Benchmark for updating a project.
func BenchmarkUpdateProject(b *testing.B) {
	db, _, projects := startupTest(b)

	b.ResetTimer() // Start benchmark timer here to exclude setup time.

	for i := 0; i < b.N; i++ {
		for _, project := range projects {
			updatedProject := project
			updatedProject.Name = "new name"
			if len(updatedProject.Tasks) > 0 {
				updatedProject.Tasks[0].Deadline = time.Now()
				if len(updatedProject.Tasks[0].Resources) > 0 {
					updatedProject.Tasks[0].Resources[0].DailyCost = new(float64)
					*updatedProject.Tasks[0].Resources[0].DailyCost = 3.14
				}
			}
			err := repository.UpdateProject(db, &updatedProject)
			if err != nil {
				b.Fatalf("Failed to update project: %v", err)
			}
		}
	}
}

// Benchmark for deleting a project.
func BenchmarkDeleteProject(b *testing.B) {
	db, _, projects := startupTest(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, project := range projects {
			err := repository.DeleteProject(db, project.ID)
			if err != nil {
				b.Fatalf("Failed to delete project: %v", err)
			}
		}
	}
}
