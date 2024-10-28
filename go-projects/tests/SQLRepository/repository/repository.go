package repository

import (
	"database/sql"
	base "m/tests/Base"
	directStruct "m/tests/DirectStruct/repository"
	columnfieldmap "m/tests/SQLRepository/columnFieldMap"
	"m/tests/SQLRepository/entities"
)

// InsertResource inserts a new resource into the RESOURCES table.
func InsertResource(db *sql.DB, resource entities.Resource) (int, error) {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}
	err = repo.Insert(&resource)
	return resource.ID, err
}

// InsertProject inserts a new project along with its associated tasks.
func InsertProject(db *sql.DB, project entities.Project) (int, error) {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}
	err = repo.Insert(&project)
	if err != nil {
		return -1, err
	}
	fk := []columnfieldmap.ColumnFieldPair{
		{ColumnName: "project_id", Field: &project.ID},
	}
	baseLink := NewBaseLinks("TASK_RESOURCE", "TASK_ID", "RESOURCE_ID")
	for _, task := range project.Tasks {
		err = repo.InsertWithFK(&task, fk)
		if err != nil {
			panic(err)
		}
		var resourcesIds []int
		for _, resource := range task.Resources {
			resourcesIds = append(resourcesIds, resource.ID)
		}
		links := baseLink.NewLinks(task.ID, resourcesIds)
		repo.Links(links)
	}
	return project.ID, err
}

// ReadProject retrieves a project by ID, including its tasks and resources associated with each task.
func ReadProject(db *sql.DB, projectID int) (*entities.Project, error) {
	tmp, err := directStruct.ReadProject(db, projectID)
	if err != nil {
		return nil, err
	}
	ret, err := base.Cast[entities.Project](tmp)
	return &ret, err
}

// UpdateProject updates the details of a project by ID.
func UpdateProject(db *sql.DB, updatedProject *entities.Project) error {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}
	err = repo.Update(updatedProject)
	return err
}

// DeleteProject deletes a project by ID.
func DeleteProject(db *sql.DB, projectID int) error {
	var project entities.Project
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}
	return repo.Delete(projectID, &project)
}
