# Testes em Go

🌍 *[**Português**](README.md) ∙ [English](README_en.md)*

Este diretório contém os subdiretórios e arquivos relacionados aos testes realizados para avaliar diferentes abordagens de acesso a dados em Go, utilizando um banco de dados PostgreSQL. Os testes incluem o uso direto de consultas SQL, o uso do ORM GORM e uma abordagem DAO genérica.

## Resultados dos Testes de Leitura

Os arquivos `.json` neste diretório representam os resultados das consultas obtidas durante a execução dos testes. Eles são cruciais para validar a consistência dos dados retornados por cada experimento, assegurando a precisão e a confiabilidade dos métodos de acesso a dados testados.

Para uma análise detalhada do desempenho de cada abordagem, consulte os resultados do benchmark disponíveis no [README principal](../README.md) do projeto. Os benchmarks oferecem insights valiosos sobre a eficiência de cada método em termos de tempo de execução e uso de recursos.

## Executando os Testes

### Executando os programas

Para executar os testes, é necessário estar no diretório `projects-go` do projeto. Abaixo estão os comandos para executar cada teste individualmente, permitindo que você avalie e compare as diferentes abordagens de acesso a dados.

Teste com gorm:
```bash
go run tests/ClassWithGorm/main.go
```

Teste com consulta única:
```bash
go run tests/ClassOneQuery/main.go
```

Teste com métodos DAO de execução de queries:
```bash
go run tests/ClassDAO/main.go
```

### Executando os benchmarks

Para executar os testes com benchmark utilize a execução de testes do VS Code ou execute os seguintes comandos.

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

#### Execução com Log

No subdiretódio `cmd` implementamos um programa que executa todos os testes completos com benchmark. Este programa registra os resultados em um arquivo `benchmark_results.log`. Para executar, no diretório `go-projects` execute o comando:

```sh
go run cmd/main.go
```