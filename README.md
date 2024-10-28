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

### CRUD Tests

In the `cmd` subdirectory, we implemented a program that runs all the complete benchmark tests. This program logs the results to a file named `benchmark_results.log`. To execute it, run the following command in the `go-projects` directory:

```sh
go run cmd/main.go
```

### Results

Using the program mentioned above, several test rounds were executed, and the results were averaged. The final outcome can be observed in the following table:

|    | Methodology   | Operation       | Time per Op (ns) | Bytes per Op | Allocs per Op |
|----|---------------|-----------------|------------------|--------------|---------------|
| 0  | DAONotation   | InsertResources | 1,095,199,700   | 1,342,200    | 36,377        |
| 1  | DAONotation   | InsertProject   | 6,619,269,900   | 3,434,856    | 84,948        |
| 2  | DAONotation   | ReadProject     | 736,151,200     | 6,382,984    | 109,964       |
| 3  | DAONotation   | UpdateProject   | 988,277,400     | 1,092,796    | 31,370        |
| 4  | DAONotation   | DeleteProject   | 11,452,412      | 4,342        | 103           |
| 5  | DirectStruct  | InsertResources | 1,183,314,100   | 698,448      | 17,853        |
| 6  | DirectStruct  | InsertProject   | 7,216,280,000   | 2,336,688    | 58,054        |
| 7  | DirectStruct  | ReadProject     | 91,569,631      | 6,757,562    | 220,025       |
| 8  | DirectStruct  | UpdateProject   | 1,000,177,050   | 503,960      | 11,732        |
| 9  | DirectStruct  | DeleteProject   | 11,524,834      | 3,705        | 83            |
| 10 | GORM          | InsertResources | 2,284,504,400   | 4,468,928    | 64,764        |
| 11 | GORM          | InsertProject   | 411,073,933     | 13,142,098   | 151,953       |
| 12 | GORM          | ReadProject     | 189,112,433     | 8,420,336    | 148,320       |
| 13 | GORM          | UpdateProject   | 42,444,254      | 99,141       | 1,235         |
| 14 | GORM          | DeleteProject   | 31,744,019      | 62,123       | 724           |
| 15 | SQLRepository | InsertResources | 1,678,694,100   | 1,384,672    | 25,985        |
| 16 | SQLRepository | InsertProject   | 12,395,136,800  | 3,521,720    | 81,748        |
| 17 | SQLRepository | ReadProject     | 196,602,150     | 9,319,502    | 262,171       |
| 18 | SQLRepository | UpdateProject   | 33,497,850      | 29,992       | 627           |
| 19 | SQLRepository | DeleteProject   | 21,979,775      | 6,534        | 118           |


The table presents performance in nanoseconds per operation (ns/op), memory usage in bytes per operation (B/op), and the number of memory allocations per operation (allocs/op), providing a comprehensive view of the efficiency of each tested approach.

The figure below shows in bar graph the results normalized by the maximum in each operation.

![picture](./output.png)

Based on the normalized data from the benchmarking experiment, the following conclusions can be drawn about the CRUD operation performance across different methodologies (`DAONotation`, `DirectStruct`, `GORM`, and `SQLRepository`):

**Time Performance**

- `GORM` demonstrates significant time inefficiency in some operations. It reaches peak performance in both `InsertResources` and `DeleteProject` operations (both normalized to 1.0), but shows a remarkably slow performance in other operations, particularly `InsertProject` and `UpdateProject`.
- `SQLRepository` is consistent with moderate time efficiency. While `InsertProject` is relatively time-intensive, other operations like `UpdateProject` demonstrate minimal time consumption, suggesting that this method is optimized for certain tasks.

**Memory Allocation (Bytes)**

- `GORM` exhibits the highest memory consumption overall, with several operations reaching a normalized value of 1.0, indicating maximum memory usage.
- `DAONotation` performs well in memory usage, with all values below 1.0 except for `UpdateProject`. It shows efficiency particularly in `DeleteProject`.
- `SQLRepository` offers a balanced memory allocation profile across operations, excelling particularly in `UpdateProject`, where memory consumption is minimal.

**Allocation Counts**

- `GORM` tends to allocate the most resources, evident in high `allocs_per_op` values for most operations, showing inefficiency in resource allocation.
- `DirectStruct` strikes a good balance in allocation efficiency, with lower normalized values for `DeleteProject` and `UpdateProject`.
- `SQLRepository` has the best allocation efficiency, especially for the `UpdateProject` operation.

For analysis details see the [notebook](ResultAnalysis.ipynb).

---

## Contributions

Contributions, corrections, and suggestions are welcome.

## License

This project is licensed under the [MIT License](LICENSE).
