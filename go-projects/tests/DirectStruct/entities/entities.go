package entities

import "time"

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Manager     string     `json:"manager"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Budget      *float64   `json:"budget"`
	Description *string    `json:"description"`
	Tasks       []Task     `json:"tasks"` // Associated tasks
}

type Task struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Responsible   *string    `json:"responsible"`
	Deadline      time.Time  `json:"deadline"`
	Status        string     `json:"status"`
	Priority      *string    `json:"priority"`
	EstimatedTime *string    `json:"estimatedTime"`
	Description   *string    `json:"description"`
	Resources     []Resource `json:"resources"` // Resources used by the task
}

type Resource struct {
	ID              int        `json:"id"`
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	DailyCost       *float64   `json:"dailyCost"`
	Status          string     `json:"status"`
	Supplier        *string    `json:"supplier"`
	Quantity        *int       `json:"quantity"`
	AcquisitionDate *time.Time `json:"acquisitionDate"`
}
