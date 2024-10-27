# Comparative Study of Database Access Performance in Golang

🌍 *[**English**](README.md) ∙ [Português](README_pt.md)*

## Description
This project explores different methods of data access in a PostgreSQL database using Go. Three different approaches were implemented and tested for reading data: a single SQL query, multiple CRUD DAO queries managed with reflection, and the GORM ORM.

### Structs Implementations

#### DAO Notation

Uses only structure declarations with additional tags. Since the idea behind this DAO is simplicity, it includes only tags indicating column names in the database.

```go
package entities

import "time"

type Project struct {
    ID          int       `db:"ID" json:"id"`
    Name        string    `db:"NAME" json:"name"`
    Manager     string    `db:"MANAGER" json:"manager"`
    StartDate   time.Time `db:"START_DATE" json:"startDate"`
    EndDate     *time.Time `db:"END_DATE" json:"endDate"`
    Budget      *float64  `db:"BUDGET" json:"budget"`
    Description *string   `db:"DESCRIPTION" json:"description"`
    Tasks       []Task    `json:"tasks"` // Associated tasks
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
    ID              int       `db:"ID" json:"id"`
    Type            string    `db:"TYPE" json:"type"`
    Name            string    `db:"NAME" json:"name"`
    DailyCost       *float64  `db:"DAILY_COST" json:"dailyCost"`
    Status          string    `db:"STATUS" json:"status"`
    Supplier        *string   `db:"SUPPLIER" json:"supplier"`
    Quantity        *int      `db:"QUANTITY" json:"quantity"`
    AcquisitionDate *time.Time `db:"ACQUISITION_DATE" json:"acquisitionDate"`
}
```

#### DirectStruct

In this approach, we only declare the structures without any additional tags.

```go
package entities

import "time"

type Project struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Manager     string    `json:"manager"`
    StartDate   time.Time `json:"startDate"`
    EndDate     *time.Time `json:"endDate"`
    Budget      *float64  `json:"budget"`
    Description *string   `json:"description"`
    Tasks       []Task    `json:"tasks"` // Associated tasks
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
    ID              int       `json:"id"`
    Type            string    `json:"type"`
    Name            string    `json:"name"`
    DailyCost       *float64  `json:"dailyCost"`
    Status          string    `json:"status"`
    Supplier        *string   `json:"supplier"`
    Quantity        *int      `json:"quantity"`
    AcquisitionDate *time.Time `json:"acquisitionDate"`
}
```

#### GORM

Uses structure declarations with additional tags. In this case, tags can be complex as they must describe relationships and database definitions.

```go
package entities

import "time"

type Project struct {
    ID          int       `gorm:"primaryKey" json:"id"`
    Name        string    `json:"name"`
    Manager     string    `json:"manager"`
    StartDate   time.Time `json:"startDate"`
    EndDate     *time.Time `json:"endDate"`
    Budget      *float64  `json:"budget"`
    Description *string   `json:"description"`
    Tasks       []Task    `gorm:"foreignKey:ProjectID" json:"tasks"` // Associated tasks
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
    ID              int       `gorm:"primaryKey" json:"id"`
    Type            string    `json:"type"`
    Name            string    `json:"name"`
    DailyCost       *float64  `json:"dailyCost"`
    Status          string    `json:"status"`
    Supplier        *string   `json:"supplier"`
    Quantity        *int      `json:"quantity"`
    AcquisitionDate *time.Time `json:"acquisitionDate"`
    Tasks           []Task    `gorm:"many2many:task_resource;" json:"tasks"` // Tasks that use this resource
}
```

#### SQLRepository

Entities are declared as simple structures but must implement an interface that defines mapping methods for the database.

```go
package entities

import "time"

// Project represents the PROJECTS table
type Project struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Manager     string    `json:"manager"`
    StartDate   time.Time `json:"startDate"`
    EndDate     *time.Time `json:"endDate"`
    Budget      *float64  `json:"budget"`
    Description *string   `json:"description"`
    Tasks       []Task    `json:"tasks"` // Associated tasks
}

func (p *Project) TableName() string {
    return "projects"
}

func (p *Project) ColumnsNames() []string {
    return []string{"id", "name", "manager", "start_date", "end_date", "budget", "description"}
}

func (p *Project) Fields() []interface{} {
    return []interface{}{&p.ID, &p.Name, &p.Manager, &p.StartDate, &p.EndDate, &p.Budget, &p.Description}
}

// Task represents the TASKS table
type Task struct {
    ID            int       `json:"id"`
    Name          string    `json:"name"`
    Responsible   *string   `json:"responsible"`
    Deadline      time.Time `json:"deadline"`
    Status        string    `json:"status"`
    Priority      *string   `json:"priority"`
    EstimatedTime *string   `json:"estimatedTime"`
    Description   *string   `json:"description"`
    Resources     []Resource `json:"resources"` // Resources used by the task
}

func (t *Task) TableName() string {
    return "tasks"
}

func (t *Task) ColumnsNames() []string {
    return []string{"id", "name", "responsible", "deadline", "status", "priority", "estimated_time", "project_id", "description"}
}

func (t *Task) Fields() []interface{} {
    return []interface{}{&t.ID, &t.Name, &t.Responsible, &t.Deadline, &t.Status, &t.Priority, &t.EstimatedTime, &t.Description}
}

// Resource represents the RESOURCES table
type Resource struct {
    ID              int       `json:"id"`
    Type            string    `json:"type"`
    Name            string    `json:"name"`
    DailyCost       *float64  `json:"dailyCost"`
    Status          string    `json:"status"`
    Supplier        *string   `json:"supplier"`
    Quantity        *int      `json:"quantity"`
    AcquisitionDate *time.Time `json:"acquisitionDate"`
}

func (r *Resource) TableName() string {
    return "resources"
}

func (r *Resource) ColumnsNames() []string {
    return []string{"id", "type", "name", "daily_cost", "status", "supplier", "quantity", "acquisition_date"}
}

func (r *Resource) Fields() []interface{} {
    return []interface{}{&r.ID, &r.Type, &r.Name, &r.DailyCost, &r.Status, &r.Supplier, &r.Quantity, &r.AcquisitionDate}
}
```

## Test Environment

To facilitate the setup, we used PostgreSQL in a Docker container. The Go project was organized with each test in the `tests` directory. Details of these components are provided in the following files:
- [README Go](./go-projects/README.md)
- [README DB](./database/README.md)

## Benchmark Results

The test environment has the following characteristics:
- **Operating System**: Windows
- **CPU Architecture**: AMD64
- **CPU**: Intel(R) Core(TM) i7-10510U @ 1.80GHz
- **Database**: PostgreSQL

---

### Initial Read Tests

#### 1. Read with a Single SQL Query
```
Package: m/tests/ReadClassOneQuery
Executions: 
- 660 executions: 1876615 ns/op, 11064 B/op, 517 allocs/op
- 771 executions: 1436036 ns/op, 11066 B/op, 517 allocs/op
- 387 executions: 3240193 ns/op, 11064 B/op, 517 allocs/op
```

#### 2. Read with DAO Implemented Using Reflection
```
Package: m/tests/ReadClassWithCrud
Executions:
- 96 executions: 12052747 ns/op, 18664 B/op, 491 allocs/op
- 100 executions: 10449300 ns/op, 18668 B/op, 491 allocs/op
- 82 executions: 15597262 ns/op, 18661 B/op, 491 allocs/op
```

#### 3. Read with GORM
```
Package: m/tests/ReadClassWithGorm
Executions:
- 298 executions: 4154921 ns/op, 51744 B/op, 955 allocs/op
- 188 executions: 6620905 ns/op, 51794 B/op, 957 allocs/op
- 196 executions: 5753415 ns/op, 51777 B/op, 957 allocs/op
```
---

### CRUD Tests

In the `cmd` subdirectory, we implemented a program that runs all the complete benchmark tests. This program logs the results to a file named `benchmark_results.log`. To execute it, run the following command in the `go-projects` directory:

```sh
go run cmd/main.go
```

### Results

Using the program mentioned above, several test rounds were executed, and the results were averaged. The final outcome can be observed in the following figure:

![picture](./resource/output.png)

The chart presents performance in nanoseconds per operation (ns/op), memory usage in bytes per operation (B/op), and the number of memory allocations per operation (allocs/op), providing a comprehensive view of the efficiency of each tested approach.

---

## Conclusion
The benchmarks reveal significant differences in performance and resource usage among the three tested approaches. As expected, reading with a single SQL query is the most efficient approach. However, in terms of resource allocation, the DAO implementation had comparable memory allocation to the single query for this example. Finally, the GORM approach, while being the most convenient in terms of development, resulted in longer execution times and higher resource usage. Additionally, during implementation, we observed that using transactions can significantly degrade performance.

---

## Contributions

Contributions, corrections, and suggestions are welcome.

## License

This project is licensed under the [MIT License](LICENSE).
