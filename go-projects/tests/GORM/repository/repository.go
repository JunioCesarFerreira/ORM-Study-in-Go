package repository

import (
	"m/tests/GORM/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InsertResource inserts a new resource into the RESOURCES table.
func InsertResource(db *gorm.DB, resource entities.Resource) (int, error) {
	if err := db.Create(&resource).Error; err != nil {
		return -1, err
	}
	return resource.ID, nil
}

// InsertProject inserts a new project along with its associated tasks.
func InsertProject(db *gorm.DB, project entities.Project) (int, error) {
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&project).Error; err != nil {
		return -1, err
	}
	return project.ID, nil
}

// ReadProject retrieves a project by ID, including its tasks and resources associated with each task.
func ReadProject(db *gorm.DB, projectID int) (*entities.Project, error) {
	var project entities.Project
	err := db.Preload("Tasks.Resources").First(&project, projectID).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates the details of a project by ID.
func UpdateProject(db *gorm.DB, updatedProject *entities.Project) error {
	// Start a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Update the main project fields
		if err := tx.Model(&entities.Project{}).Where("id = ?", updatedProject.ID).Updates(updatedProject).Error; err != nil {
			return err
		}

		// Update tasks associated with the project
		for _, task := range updatedProject.Tasks {
			task.ProjectID = updatedProject.ID // Ensure ProjectID is set
			if err := tx.Save(&task).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// DeleteProject deletes a project by ID.
func DeleteProject(db *gorm.DB, projectID int) error {
	return db.Delete(&entities.Project{}, projectID).Error
}
