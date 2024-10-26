package entities

import "time"

type Project struct {
	ID          int        `db:"ID" json:"id"`
	Name        string     `db:"NAME" json:"name"`
	Manager     string     `db:"MANAGER" json:"manager"`
	StartDate   time.Time  `db:"START_DATE" json:"startDate"`
	EndDate     *time.Time `db:"END_DATE" json:"endDate"`
	Budget      *float64   `db:"BUDGET" json:"budget"`
	Description *string    `db:"DESCRIPTION" json:"description"`
	Tasks       []Task     `json:"tasks"` // Associated tasks
}

type Task struct {
	ID            int        `db:"ID" json:"id"`
	Name          string     `db:"NAME" json:"name"`
	Responsible   *string    `db:"RESPONSIBLE" json:"responsible"`
	Deadline      time.Time  `db:"DEADLINE" json:"deadline"`
	Status        string     `db:"STATUS" json:"status"`
	Priority      *string    `db:"PRIORITY" json:"priority"`
	EstimatedTime *string    `db:"ESTIMATED_TIME" json:"estimatedTime"`
	Description   *string    `db:"DESCRIPTION" json:"description"`
	Resources     []Resource `json:"resources"` // Resources used by the task
}

type Resource struct {
	ID              int        `db:"ID" json:"id"`
	Type            string     `db:"TYPE" json:"type"`
	Name            string     `db:"NAME" json:"name"`
	DailyCost       *float64   `db:"DAILY_COST" json:"dailyCost"`
	Status          string     `db:"STATUS" json:"status"`
	Supplier        *string    `db:"SUPPLIER" json:"supplier"`
	Quantity        *int       `db:"QUANTITY" json:"quantity"`
	AcquisitionDate *time.Time `db:"ACQUISITION_DATE" json:"acquisitionDate"`
}
