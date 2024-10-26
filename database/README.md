# Database Schema

🌍 *[**English**](README.md) ∙ [Português](README_pt.md)*

In this directory, you will find the database schema used in the tests and experiments.

We have the following entities:

![ER Diagram](er-diagram.png)

The main objective of this model is to test different types of data and `1-N` and `N-N` relationships.

See the table creation script [schema](schema.sql).

# Docker

To create a Docker container with the database, run the following commands in this directory:

```shell
docker build -t my-db-image .
```

```shell
docker run --name my-container-db -p 5432:5432 -d my-db-image
```