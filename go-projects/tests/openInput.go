package tests

import (
	"io"
	"log"
	"os"
)

func OpenInputJson() ([]byte, error) {
	jsonFile, err := os.Open("../input.json")
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo JSON: %s", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	return byteValue, err
}
