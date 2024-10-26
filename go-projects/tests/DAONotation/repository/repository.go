package repository

import (
	"database/sql"
	"m/tests/DAONotation/dao"
	"m/tests/DAONotation/entities"
)

// InsertResource inserts a single resource into the RESOURCES table.
func InsertResource(db *sql.DB, resource entities.Resource) (int, error) {
	daoResource := dao.NewDAO(db)
	resourceId, err := daoResource.Create("RESOURCES", &resource)
	if err != nil {
		return -1, err
	}
	return resourceId, nil
}

// Inserts a project and its associated tasks and resources.
func InsertProject(db *sql.DB, project entities.Project) (int, error) {
	daoProject := dao.NewDAO(db)
	projectId, err := daoProject.Create("PROJECTS", &project)
	if err != nil {
		return projectId, err
	}

	for _, task := range project.Tasks {
		taskId, err := daoProject.CreateChild("TASKS", &task, "PROJECT_ID", projectId)
		if err != nil {
			return projectId, err
		}
		for _, resource := range task.Resources {
			_, err := daoProject.CreateWithLinkSingleSide(taskId, "RESOURCES", "TASK_RESOURCE", resource.ID, "TASK_ID", "RESOURCE_ID")
			if err != nil {
				return projectId, err
			}
		}
	}
	return projectId, nil
}

// Updates an existing project.
func UpdateProject(db *sql.DB, project *entities.Project) error {
	daoProject := dao.NewDAO(db)
	err := daoProject.Update("PROJECTS", project)
	if err != nil {
		return err
	}
	for _, task := range project.Tasks {
		err := daoProject.Update("TASKS", &task)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletes a project by ID.
func DeleteProject(db *sql.DB, projectID int) error {
	daoProject := dao.NewDAO(db)
	return daoProject.Delete("PROJECTS", projectID)
}

// Reads a project, along with its associated tasks and resources, by project ID.
func ReadProject(db *sql.DB, projectID int) (*entities.Project, error) {
	project := &entities.Project{}
	daoProject := dao.NewDAO(db)

	err := daoProject.Read("PROJECTS", projectID, project)
	if err != nil {
		return project, err
	}

	condition := "PROJECT_ID = $1"
	args := []interface{}{projectID}
	modelTask := &entities.Task{}

	// Fetch tasks associated with the project
	tasks, err := daoProject.ReadMultiple("TASKS", condition, args, modelTask)
	if err != nil {
		return project, err
	}

	for _, task := range tasks {
		tmpTask := *task.(*entities.Task)
		condition := "TASK_ID = $1"
		args := []interface{}{tmpTask.ID}
		modelResource := &entities.Resource{}

		// Fetch resources associated with each task
		resources, err := daoProject.ReadMultiple("TASK_RESOURCE_VIEW", condition, args, modelResource)
		if err != nil {
			return project, err
		}

		for _, resource := range resources {
			tmpTask.Resources = append(tmpTask.Resources, *resource.(*entities.Resource))
		}

		project.Tasks = append(project.Tasks, tmpTask)
	}

	return project, nil
}
