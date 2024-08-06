# Comparative Study of Database Access Performance in Golang

🌍 *[Português](README.md) ∙ [English](README_en.md)*

## Description
This project explores different methods of data access in a PostgreSQL database using Go. Three different approaches were implemented and tested for reading data: a single SQL query, reflection DAO CRUD queries, and the GORM ORM.

## Test Environment

To make it easier to prepare the environment, we used PostgreSQL in a Docker container. While the go project was organized with each test in the `tests` directory. The details of these parts are detailed in the files:
- [README go](./go-projects/README_en.md)
- [README db](./database/README_en.md)

## Benchmark Results

The environment used in the tests has the following characteristics:

- **Operating System**: Windows
- **CPU Architecture**: AMD64
- **CPU**: Intel(R) Core(TM) i7-10510U @ 1.80GHz
- **Database**: PostgreSQL

---

### Initial Read-Only Tests

#### 1. Reading with Single SQL Query
```
Package: m/tests/ReadClassOneQuery
Runs: 
- 321 runs: 3243073 ns/op, 19293 B/op, 920 allocs/op
- 379 runs: 2810400 ns/op, 19288 B/op, 920 allocs/op
- 465 runs: 2960930 ns/op, 19291 B/op, 920 allocs/op
```

#### 2. Reading DAO CRUD
```
Package: m/tests/ReadClassWithCrud
Runs:
- 58 runs: 17997191 ns/op, 31874 B/op, 817 allocs/op
- 68 runs: 17424975 ns/op, 31870 B/op, 817 allocs/op
- 58 runs: 18148195 ns/op, 31874 B/op, 817 allocs/op
```

#### 3. Reading with GORM
```
Package: m/tests/ReadClassWithGorm
Runs:
- 256 runs: 4359351 ns/op, 74645 B/op, 1480 allocs/op
- 252 runs: 5066758 ns/op, 74634 B/op, 1480 allocs/op
- 242 runs: 4249418 ns/op, 74619 B/op, 1480 allocs/op
```

---

### Testing with CRUD

---

## Conclusion
The benchmarks reveal significant differences in performance and resource usage among the three tested approaches. The single SQL query reading proved to be the most efficient in terms of execution time and memory allocation. The manual DAO approach, although slower, maintained moderate memory usage. Finally, the GORM approach, despite being the most convenient in terms of development, resulted in higher execution time and greater resource usage.

---

## Contributions

Contributions, corrections, and suggestions are welcome.

## License

This project is licensed under the [MIT License](LICENSE).