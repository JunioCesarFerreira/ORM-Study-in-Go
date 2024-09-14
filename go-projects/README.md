# Tests in Go

🌍 *[**English**](README_en.md) ∙ [Português](README.md)*

This directory contains subdirectories and files related to tests conducted to evaluate different approaches to data access in Go, using a PostgreSQL database. The tests include direct SQL queries, the use of the GORM ORM, and a generic DAO approach.

## Results of Reading Tests

The `.json` files in this directory represent the query results obtained during the test executions. They are crucial for validating the consistency of the data returned by each experiment, ensuring the accuracy and reliability of the tested data access methods.

For a detailed analysis of the performance of each approach, please refer to the benchmark results available in the project's [main README](../README.md). The benchmarks provide valuable insights into the efficiency of each method in terms of execution time and resource usage.

## Running the Tests

### Running the programs

To run the tests, you need to be in the project's `projects-go` directory. Below are the commands to run each test individually, allowing you to evaluate and compare the different data access approaches.

Test with GORM:
```bash
go run tests/ClassWithGorm/main.go
```

Test with single query:
```bash
cd tests/ClassOneQuery
go test -benchmem -run=^_test$ -bench . ./...
```

Test with generic (CRUD) query execution methods:
```bash
go run tests/ClassDAO/main.go
```

### Running the benchmarks

To run the benchmark tests, use the test execution feature of VS Code or execute the following commands.

```bash
cd tests/ClassWithGorm
go test -benchmem -run=^_test$ -bench . ./...
```

```bash
cd tests/ClassOneQuery
go test -benchmem -run=^_test$ -bench . ./...
```

```bash
cd tests/ClassDAO
go test -benchmem -run=^_test$ -bench . ./...
```

#### Execution with Logging

In the `cmd` subdirectory, we implemented a program that runs all the complete benchmark tests. This program logs the results to a file named `benchmark_results.log`. To execute it, run the following command in the `go-projects` directory:

```sh
go run cmd/main.go
```