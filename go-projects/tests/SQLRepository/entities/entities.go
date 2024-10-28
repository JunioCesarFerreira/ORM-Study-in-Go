package entities

import (
	columnfieldmap "m/tests/SQLRepository/columnFieldMap"
	"time"
)

// Project represents the PROJECTS table
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

func (p *Project) TableName() string {
	return "projects"
}

func (p *Project) Mapped() []columnfieldmap.ColumnFieldPair {
	return []columnfieldmap.ColumnFieldPair{
		{ColumnName: "id", Field: &p.ID},
		{ColumnName: "name", Field: &p.Name},
		{ColumnName: "manager", Field: &p.Manager},
		{ColumnName: "start_date", Field: &p.StartDate},
		{ColumnName: "end_date", Field: &p.EndDate},
		{ColumnName: "budget", Field: &p.Budget},
		{ColumnName: "description", Field: &p.Description},
	}
}

func (p *Project) ColumnsNames() []string {
	return columnfieldmap.ColumnsNames(p)
}

func (p *Project) Fields() []interface{} {
	return columnfieldmap.Fields(p)
}

func (p *Project) PKColNames() []string {
	return []string{"id"}
}

func (p *Project) PKFields() []interface{} {
	return []interface{}{&p.ID}
}

// Task represents the TASKS table
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

func (t *Task) TableName() string {
	return "tasks"
}

func (t *Task) Mapped() []columnfieldmap.ColumnFieldPair {
	return []columnfieldmap.ColumnFieldPair{
		{ColumnName: "id", Field: &t.ID},
		{ColumnName: "name", Field: &t.Name},
		{ColumnName: "responsible", Field: &t.Responsible},
		{ColumnName: "deadline", Field: &t.Deadline},
		{ColumnName: "status", Field: &t.Status},
		{ColumnName: "priority", Field: &t.Priority},
		{ColumnName: "estimated_time", Field: &t.EstimatedTime},
		{ColumnName: "description", Field: &t.Description},
	}
}

func (t *Task) ColumnsNames() []string {
	return columnfieldmap.ColumnsNames(t)
}

func (t *Task) Fields() []interface{} {
	return columnfieldmap.Fields(t)
}

func (t *Task) PKColNames() []string {
	return []string{"id"}
}

func (t *Task) PKFields() []interface{} {
	return []interface{}{&t.ID}
}

// Resource represents the RESOURCES table
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

func (r *Resource) TableName() string {
	return "resources"
}

func (r *Resource) Mapped() []columnfieldmap.ColumnFieldPair {
	return []columnfieldmap.ColumnFieldPair{
		{ColumnName: "id", Field: &r.ID},
		{ColumnName: "type", Field: &r.Type},
		{ColumnName: "name", Field: &r.Name},
		{ColumnName: "daily_cost", Field: &r.DailyCost},
		{ColumnName: "status", Field: &r.Status},
		{ColumnName: "supplier", Field: &r.Supplier},
		{ColumnName: "quantity", Field: &r.Quantity},
		{ColumnName: "acquisition_date", Field: &r.AcquisitionDate},
	}
}

func (r *Resource) ColumnsNames() []string {
	return columnfieldmap.ColumnsNames(r)
}

func (r *Resource) Fields() []interface{} {
	return columnfieldmap.Fields(r)
}

func (r *Resource) PKColNames() []string {
	return []string{"id"}
}

func (r *Resource) PKFields() []interface{} {
	return []interface{}{&r.ID}
}
