package tests

import (
	"encoding/json"
	"log"
	"os"
)

func SaveResult(fileName string, obj interface{}) {
	// Abre (ou cria) o arquivo onde o JSON será salvo
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Cria um encoder que escreverá no arquivo
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Configura o encoder para formatar o JSON com indentação

	// Serializa o objeto em JSON e escreve no arquivo
	err = encoder.Encode(obj)
	if err != nil {
		log.Fatalf("Failed to encode object to JSON: %v", err)
	}

	log.Printf("Object saved to %s", fileName)
}
