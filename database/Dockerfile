FROM postgres:latest

# Dados de conexão do banco de dados
ENV POSTGRES_DB=my_database
ENV POSTGRES_USER=my_user
ENV POSTGRES_PASSWORD=my@Pass%1234

# Copia esquemático de tabelas do banco de dados
COPY schema.sql /docker-entrypoint-initdb.d/schema.sql

RUN chmod -R 755 /docker-entrypoint-initdb.d
