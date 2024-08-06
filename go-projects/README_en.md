# Tests in Go

🌍 *[Português](README.md) ∙ [English](README_en.md)*

This directory contains subdirectories and files related to tests conducted to evaluate different approaches to data access in Go, using a PostgreSQL database. The tests include direct SQL queries, the use of the GORM ORM, and a generic CRUD DAO approach.

## Test Results

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
go run tests/ClassDAOcrud/main.go
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
cd tests/ClassDAOcrud
go test -benchmem -run=^_test$ -bench . ./...
```