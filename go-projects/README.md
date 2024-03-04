# Testes em Go

Este diretório contém os subdiretórios e arquivos relacionados aos testes realizados para avaliar diferentes abordagens de acesso a dados em Go, utilizando um banco de dados PostgreSQL. Os testes incluem o uso direto de consultas SQL, o uso do ORM GORM e uma abordagem CRUD genérica.

## Resultados dos Testes

Os arquivos `.json` neste diretório representam os resultados das consultas obtidas durante a execução dos testes. Eles são cruciais para validar a consistência dos dados retornados por cada experimento, assegurando a precisão e a confiabilidade dos métodos de acesso a dados testados.

Para uma análise detalhada do desempenho de cada abordagem, consulte os resultados do benchmark disponíveis no [README principal](../README.md) do projeto. Os benchmarks oferecem insights valiosos sobre a eficiência de cada método em termos de tempo de execução e uso de recursos.

## Executando os Testes

### Executando os programas

Para executar os testes, é necessário estar no diretório `projects-go` do projeto. Abaixo estão os comandos para executar cada teste individualmente, permitindo que você avalie e compare as diferentes abordagens de acesso a dados.

Teste com gorm:
```bash
go run tests/ReadClassWithGorm/main.go
```

Teste com consulta única:
```bash
go run tests/ReadClassOneQuery/main.go
```

Teste com métodos genéricos (crud) de execução de queries:
```bash
go run tests/ReadClassWithCrud/main.go
```

### Executando os benchmarks

Para executar os testes com benchmark utilize a execução de testes do VS Code ou execute os seguintes comandos.

```bash
go test -benchmem -run=^$ -bench ^BenchmarkReadClass$ m/tests/ReadClassWithGorm
```

```bash
go test -benchmem -run=^$ -bench ^BenchmarkReadClass$ m/tests/ReadClassOneQuery
```

```bash
go test -benchmem -run=^$ -bench ^BenchmarkReadClass$ m/tests/ReadClassWithCrud
```

