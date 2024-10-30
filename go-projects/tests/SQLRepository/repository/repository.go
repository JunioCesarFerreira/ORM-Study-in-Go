package repository

import (
	"database/sql"
	columnfieldmap "m/tests/SQLRepository/columnFieldMap"
	"m/tests/SQLRepository/entities"
	"time"
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

// UpdateProject updates the details of a project by ID.
func UpdateProject(db *sql.DB, updatedProject *entities.Project) error {
	repo, err := NewSQLRepository(db)
	if err != nil {
		panic(err)
	}
	err = repo.Update(updatedProject)
	for _, task := range updatedProject.Tasks {
		err = repo.Update(&task)
		if err != nil {
			return err
		}
	}
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
