package entities

import "time"

type Project struct {
	ID          int        `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name        string     `json:"name"`
	Manager     string     `json:"manager"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Budget      *float64   `json:"budget"`
	Description *string    `json:"description"`
	Tasks       []Task     `gorm:"foreignKey:ProjectID" json:"tasks"` // Associated tasks
}

type Task struct {
	ID            int        `gorm:"primaryKey" json:"id"`
	ProjectID     int        `json:"-"`
	Name          string     `json:"name"`
	Responsible   *string    `json:"responsible"`
	Deadline      time.Time  `json:"deadline"`
	Status        string     `json:"status"`
	Priority      *string    `json:"priority"`
	EstimatedTime *string    `json:"estimatedTime"`
	Description   *string    `json:"description"`
	Resources     []Resource `gorm:"many2many:task_resource;" json:"resources"` // Resources used by the task
}

type Resource struct {
	ID              int        `gorm:"primaryKey" json:"id"`
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	DailyCost       *float64   `json:"dailyCost"`
	Status          string     `json:"status"`
	Supplier        *string    `json:"supplier"`
	Quantity        *int       `json:"quantity"`
	AcquisitionDate *time.Time `json:"acquisitionDate"`
	Tasks           []Task     `gorm:"many2many:task_resource;" json:"tasks"` // Tasks that use this resource
}
