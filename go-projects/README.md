# Testes Go

Neste diretório está o subdiretório com os testes realizados.

Os arquivos `json` são os resultados das consultas realizadas em uma execução dos testes. Estes resultados servem para validar que todos os experimentos retornam os mesmos dados.

Os resultados do benchmark são apresentados no [README principal](../README.md).

## Como usar

Para executar as consultas dos testes em console, navegue até o nível `projects-go`. Os comandos a seguir servem para executar cada teste.

Teste com gorm:
```shell
go run tests/ReadClassWithGorm/main.go
```

Teste com consulta única:
```shell
go run tests/ReadClassOneQuery/main.go
```

Teste com métodos genéricos (crud) de execução de queries:
```shell
go run tests/ReadClassWithCrud/main.go
```