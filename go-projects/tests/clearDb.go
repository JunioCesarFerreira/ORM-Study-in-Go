package tests

import (
	"database/sql"
	"fmt"
)

// Função para limpar todas as tabelas
func ClearDatabase(db *sql.DB) error {
	// Ordem correta para deletar dados, respeitando dependências de chaves estrangeiras
	queries := []string{
		"DELETE FROM OBJECT_ITEM_LINK;",
		"DELETE FROM ITEMS;",
		"DELETE FROM OBJECTS;",
		"DELETE FROM CLASSES;",
	}

	// Executando as queries em ordem
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("erro ao executar query '%s': %v", query, err)
		}
	}
	return nil
}
