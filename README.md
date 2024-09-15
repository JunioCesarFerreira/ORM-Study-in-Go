# Comparative Study of Database Access Performance in Golang

🌍 *[**English**](README.md) ∙ [Português](README_pt.md)*

## Description
This project explores different methods of data access in a PostgreSQL database using Go. Three different approaches were implemented and tested for reading data: a single SQL query, multiple CRUD DAO queries managed with reflection, and the GORM ORM.

### Implementations

- [OneQuery](./go-projects/tests/ClassOneQuery/repository/repository.go): For this approach, we implemented the database access methods using SQL commands written directly in the code. Input parameters are passed separately to the standard SQL library, which prevents SQL injection.

- [DAO](./go-projects/tests/ClassDAO/dao/dao.go): For this approach, SQL commands are constructed generically using reflection. In this implementation, we use tags to indicate the column names in the database.

- [ORM](./go-projects/tests/ClassWithGorm/repository/repository.go): For this approach, we used one of the most popular ORMs for Go, GORM. The implementation was done following the framework's specifications.

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
