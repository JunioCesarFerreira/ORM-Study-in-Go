# Database

Neste diretório está o esquemático de banco de dados utilizado nos testes e exeperimentos.

Temos as seguintes entidades:

[DiagramaER](er-diagram.png)

# Docker

Para gerar um container Docker com o banco de dados execute neste diretório os seguintes comandos:

```shell
docker build -t my-db-image .
```

```shell
docker run --name my-container-db -p 5432:5432 -d my-db-image
```