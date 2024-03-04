# Estudo comparativo de desenpenho de acesso a banco de dados em Golang

## Descrição
Este projeto explora diferentes métodos de acesso a dados em um banco de dados PostgreSQL usando Go. Foram implementados e testados três abordagens diferentes para ler dados: uma consulta SQL única, múltiplas consultas CRUD gerenciadas manualmente e o ORM GORM.

## Ambiente de Teste
- **Sistema Operacional**: Windows
- **Arquitetura do CPU**: AMD64
- **CPU**: Intel(R) Core(TM) i7-10510U @ 1.80GHz
- **Banco de Dados**: PostgreSQL

## Resultados dos Benchmarks

### 1. Leitura com Consulta SQL Única
```
Pacote: m/tests/ReadClassOneQuery
Execuções: 
- 321 execuções: 3243073 ns/op, 19293 B/op, 920 allocs/op
- 379 execuções: 2810400 ns/op, 19288 B/op, 920 allocs/op
- 465 execuções: 2960930 ns/op, 19291 B/op, 920 allocs/op
```

### 2. Leitura com CRUD Manual
```
Pacote: m/tests/ReadClassWithCrud
Execuções:
- 58 execuções: 17997191 ns/op, 31874 B/op, 817 allocs/op
- 68 execuções: 17424975 ns/op, 31870 B/op, 817 allocs/op
- 58 execuções: 18148195 ns/op, 31874 B/op, 817 allocs/op
```

### 3. Leitura com GORM
```
Pacote: m/tests/ReadClassWithGorm
Execuções:
- 256 execuções: 4359351 ns/op, 74645 B/op, 1480 allocs/op
- 252 execuções: 5066758 ns/op, 74634 B/op, 1480 allocs/op
- 242 execuções: 4249418 ns/op, 74619 B/op, 1480 allocs/op
```

## Conclusão
Os benchmarks revelam diferenças significativas no desempenho e no uso de recursos entre as três abordagens testadas. A leitura com consulta SQL única se mostrou a mais eficiente em termos de tempo de execução e alocação de memória. A abordagem CRUD manual, embora mais lenta, manteve um uso moderado de memória. Por fim, a abordagem com GORM, apesar de ser a mais prática em termos de desenvolvimento, resultou em maior tempo de execução e maior uso de recursos.
