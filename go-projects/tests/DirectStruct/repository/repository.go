package repository

import (
	"database/sql"
	"m/tests/DirectStruct/entities"
	"time"
)

// InsertResource inserts a single resource into the RESOURCES table.
func InsertResource(db *sql.DB, resource entities.Resource) (int, error) {
	query := `
		INSERT INTO RESOURCES (ID, TYPE, NAME, DAILY_COST, STATUS, SUPPLIER, QUANTITY, ACQUISITION_DATE)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING ID
	`
	var resourceID int
	err := db.QueryRow(query, resource.ID, resource.Type, resource.Name, resource.DailyCost, resource.Status, resource.Supplier, resource.Quantity, resource.AcquisitionDate).Scan(&resourceID)
	if err != nil {
		return 0, err
	}

	return resourceID, nil
}

// InsertProject inserts a project along with its tasks and linked resources.
func InsertProject(db *sql.DB, project entities.Project) (int, error) {
	// Insert the main project
	query := `
		INSERT INTO PROJECTS (ID, NAME, MANAGER, START_DATE, END_DATE, BUDGET, DESCRIPTION)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING ID
	`
	var projectID int
	err := db.QueryRow(query, project.ID, project.Name, project.Manager, project.StartDate, project.EndDate, project.Budget, project.Description).Scan(&projectID)
	if err != nil {
		return 0, err
	}

	// Insert each task associated with the project
	for _, task := range project.Tasks {
		taskQuery := `
			INSERT INTO TASKS (ID, NAME, RESPONSIBLE, DEADLINE, STATUS, PRIORITY, ESTIMATED_TIME, PROJECT_ID, DESCRIPTION)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING ID
		`
		var taskID int
		err := db.QueryRow(taskQuery, task.ID, task.Name, task.Responsible, task.Deadline, task.Status, task.Priority, task.EstimatedTime, projectID, task.Description).Scan(&taskID)
		if err != nil {
			return projectID, err
		}

		// Insert each resource associated with the task into the TASK_RESOURCE link table
		for _, resource := range task.Resources {
			linkQuery := `
				INSERT INTO TASK_RESOURCE (TASK_ID, RESOURCE_ID, QUANTITY_USED)
				VALUES ($1, $2, $3)
			`
			_, err := db.Exec(linkQuery, taskID, resource.ID, resource.Quantity)
			if err != nil {
				return projectID, err
			}
		}
	}

	return projectID, nil
}

// Reads a project by ID, including its tasks and resources.
func ReadProject(db *sql.DB, projectID int) (*entities.Project, error) {
	query := `
	SELECT 
		p.NAME, 
		p.MANAGER, 
		p.START_DATE, 
		p.END_DATE, 
		p.BUDGET, 
		p.DESCRIPTION, 
		t.ID, 
		t.NAME, 
		t.RESPONSIBLE, 
		t.DEADLINE, 
		t.STATUS, 
		t.PRIORITY, 
		t.ESTIMATED_TIME, 
		t.DESCRIPTION,
		r.ID, 
		r.TYPE, 
		r.NAME, 
		r.DAILY_COST, 
		r.STATUS, 
		r.SUPPLIER, 
		r.QUANTITY, 
		r.ACQUISITION_DATE
	FROM PROJECTS p
		LEFT JOIN TASKS t ON p.ID = t.PROJECT_ID
		LEFT JOIN TASK_RESOURCE tr ON t.ID = tr.TASK_ID
		LEFT JOIN RESOURCES r ON r.ID = tr.RESOURCE_ID
	WHERE p.ID = $1
	`

	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := &entities.Project{ID: projectID}
	taskMap := make(map[int]*entities.Task)

	var tID, rID sql.NullInt32
	var pName, pManager, tName, tResponsible, tStatus, tPriority, rType, rName, rStatus, rSupplier sql.NullString
	var pStartDate, tDeadline time.Time
	var pEndDate, rAcquisitionDate sql.NullTime
	var pBudget, rDailyCost sql.NullFloat64
	var tEstimatedTime, tDescription sql.NullString
	var rQuantity sql.NullInt32
	var pDescription sql.NullString

	for rows.Next() {
		err := rows.Scan(
			&pName, &pManager, &pStartDate, &pEndDate, &pBudget, &pDescription,
			&tID, &tName, &tResponsible, &tDeadline, &tStatus, &tPriority, &tEstimatedTime, &tDescription,
			&rID, &rType, &rName, &rDailyCost, &rStatus, &rSupplier, &rQuantity, &rAcquisitionDate,
		)
		if err != nil {
			return nil, err
		}

		project.Name = pName.String
		project.Manager = pManager.String
		project.StartDate = pStartDate
		project.EndDate = parseTimePtr(pEndDate)
		project.Budget = parseFloatPtr(pBudget)
		project.Description = parseStringPtr(pDescription)

		if tID.Valid {
			task, exists := taskMap[int(tID.Int32)]
			if !exists {
				task = &entities.Task{
					ID:            int(tID.Int32),
					Name:          tName.String,
					Responsible:   parseStringPtr(tResponsible),
					Deadline:      tDeadline,
					Status:        tStatus.String,
					Priority:      parseStringPtr(tPriority),
					EstimatedTime: parseStringPtr(tEstimatedTime),
					Description:   parseStringPtr(tDescription),
				}
				taskMap[int(tID.Int32)] = task
			}

			if rID.Valid {
				resource := entities.Resource{
					ID:              int(rID.Int32),
					Type:            rType.String,
					Name:            rName.String,
					DailyCost:       parseFloatPtr(rDailyCost),
					Status:          rStatus.String,
					Supplier:        parseStringPtr(rSupplier),
					Quantity:        parseIntPtr(rQuantity),
					AcquisitionDate: parseTimePtr(rAcquisitionDate),
				}
				task.Resources = append(task.Resources, resource)
			}
		}
	}

	for _, task := range taskMap {
		project.Tasks = append(project.Tasks, *task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return project, nil
}

// UpdateProject updates a project and its associated tasks.
func UpdateProject(db *sql.DB, project *entities.Project) error {
	// Update the main project attributes
	query := `
		UPDATE PROJECTS
		SET NAME = $1, MANAGER = $2, START_DATE = $3, END_DATE = $4, BUDGET = $5, DESCRIPTION = $6
		WHERE ID = $7
	`
	_, err := db.Exec(query, project.Name, project.Manager, project.StartDate, project.EndDate, project.Budget, project.Description, project.ID)
	if err != nil {
		return err
	}

	// Update each task associated with the project
	for _, task := range project.Tasks {
		taskQuery := `
			UPDATE TASKS
				SET NAME = $1, 
				RESPONSIBLE = $2, 
				DEADLINE = $3, 
				STATUS = $4, 
				PRIORITY = $5, 
				ESTIMATED_TIME = $6,
				DESCRIPTION = $7
			WHERE ID = $8 AND PROJECT_ID = $9
		`
		_, err := db.Exec(taskQuery, task.Name, task.Responsible, task.Deadline, task.Status, task.Priority, task.EstimatedTime, task.Description, task.ID, project.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Deletes a project by ID.
func DeleteProject(db *sql.DB, projectID int) error {
	query := `
		DELETE FROM PROJECTS
		WHERE ID = $1
	`
	_, err := db.Exec(query, projectID)
	if err != nil {
		return err
	}

	return nil
}

func parseStringPtr(input sql.NullString) *string {
	if input.Valid {
		return &input.String
	}
	return nil
}

func parseIntPtr(input sql.NullInt32) *int {
	if input.Valid {
		value := int(input.Int32)
		return &value
	}
	return nil
}

func parseTimePtr(input sql.NullTime) *time.Time {
	if input.Valid {
		return &input.Time
	}
	return nil
}

func parseFloatPtr(input sql.NullFloat64) *float64 {
	if input.Valid {
		return &input.Float64
	}
	return nil
}
