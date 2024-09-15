package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Define the test directories and the commands you want to run
	tests := []string{
		"tests/ClassWithGorm",
		"tests/ClassOneQuery",
		"tests/ClassDAO",
		"tests/SQLRepository",
	}

	// Set the log file name and open it for writing
	logFileName := "benchmark_results.log"
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer logFile.Close()

	for _, testDir := range tests {
		fmt.Printf("Running tests in the directory: %s\n", testDir)
		result := runBenchmark(testDir)
		fmt.Println(result)

		writeToLog(logFile, testDir, result)
	}

	fmt.Printf("Benchmark results were recorded in: %s\n", logFileName)
}

func runBenchmark(testDir string) string {
	cmd := exec.Command("go", "test", "-benchmem", "-run=^_test$", "-bench", ".", "./...")
	cmd.Dir = testDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running benchmark on %s: %v\n", testDir, err)
	}

	return string(output)
}

func writeToLog(logFile *os.File, testDir string, result string) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("\n\n==== Benchmark Results for %s ====\nTimestamp: %s\n%s\n", testDir, timestamp, result)

	_, err := logFile.WriteString(logEntry)
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
