# Database Schema

Neste diretório está o esquemático de banco de dados utilizado nos testes e exeperimentos.

Temos as seguintes entidades:

![DiagramaER](er-diagram.png)

O principal objetivo deste modelo é textar diferentes tipos de dados e relacionamentos `1-N` e `N-N`.

Veja o script de criação das tabelas [schema](schema.sql).

# Docker

Para gerar um container Docker com o banco de dados execute neste diretório os seguintes comandos:

```shell
docker build -t my-db-image .
```

```shell
docker run --name my-container-db -p 5432:5432 -d my-db-image
```

# Script Python de inserção de dados

O script [data_insert](data_insert.py) pode ser utilizado para gerar dados na estrutura proposta. Observe que é possível paramentrizar neste script as quantidade de dados a serem inseridos.

Para executar o script é necessário ter o Python 3.12 instalado.