package main

import (
	"database/sql"
	base "m/tests/Base"
	"m/tests/DAONotation/entities"
	"m/tests/DAONotation/repository"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func startupTest(b *testing.B) (*sql.DB, []entities.Resource, []entities.Project) {
	db := base.SetupDB()

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
	db, resources, _ := startupTest(b)
	defer db.Close()

	err := base.ClearAllProjectsAndResources(db)
	if err != nil {
		b.Fatalf("Error cleaning database: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, resource := range resources {
			_, err := repository.InsertResource(db, resource)
			if err != nil {
				b.Fatalf("Failed to insert resource: %v", err)
			}
		}
	}
}

// Benchmark for inserting a project.
func BenchmarkInsertProject(b *testing.B) {
	db, _, projects := startupTest(b)
	defer db.Close()

	b.ResetTimer()
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
	defer db.Close()

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
	b.StopTimer()
}

// Benchmark for updating a project.
func BenchmarkUpdateProject(b *testing.B) {
	db, _, projects := startupTest(b)
	defer db.Close()

	b.ResetTimer()
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
	defer db.Close()

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
