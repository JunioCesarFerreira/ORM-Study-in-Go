package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Defina os diretórios de teste e os comandos que deseja executar
	tests := []string{
		"tests/ClassWithGorm",
		"tests/ClassOneQuery",
		"tests/ClassDAO",
	}

	// Defina o nome do arquivo de log e abra-o para escrita
	logFileName := "benchmark_results.log"
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo de log: %v\n", err)
		return
	}
	defer logFile.Close()

	// Iterar sobre os testes, executar os comandos e registrar os resultados
	for _, testDir := range tests {
		fmt.Printf("Executando testes no diretório: %s\n", testDir)
		result := runBenchmark(testDir)
		fmt.Println(result)

		// Escrever os resultados no arquivo de log
		writeToLog(logFile, testDir, result)
	}

	fmt.Printf("Os resultados dos benchmarks foram registrados em: %s\n", logFileName)
}

func runBenchmark(testDir string) string {
	cmd := exec.Command("go", "test", "-benchmem", "-run=^_test$", "-bench", ".", "./...")
	cmd.Dir = testDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Erro ao executar benchmark em %s: %v\n", testDir, err)
	}

	return string(output)
}

func writeToLog(logFile *os.File, testDir string, result string) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("\n\n==== Benchmark Results for %s ====\nTimestamp: %s\n%s\n", testDir, timestamp, result)

	_, err := logFile.WriteString(logEntry)
	if err != nil {
		fmt.Printf("Erro ao escrever no arquivo de log: %v\n", err)
	}
}
