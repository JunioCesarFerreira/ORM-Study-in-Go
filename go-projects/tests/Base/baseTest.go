package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

type BaseProject struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Manager     string     `json:"manager"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Budget      *float64   `json:"budget"`
	Description *string    `json:"description"`
	Tasks       []BaseTask `json:"tasks"`
}

type BaseTask struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Responsible   *string        `json:"responsible"`
	Deadline      time.Time      `json:"deadline"`
	Status        string         `json:"status"`
	Priority      *string        `json:"priority"`
	EstimatedTime *string        `json:"estimatedTime"`
	Description   *string        `json:"description"`
	Resources     []BaseResource `json:"resources"`
}

type BaseResource struct {
	ID              int        `json:"id"`
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	DailyCost       *float64   `json:"dailyCost"`
	Status          string     `json:"status"`
	Supplier        *string    `json:"supplier"`
	Quantity        *int       `json:"quantity"`
	AcquisitionDate *time.Time `json:"acquisitionDate"`
}

type TestInput struct {
	Resources []BaseResource `json:"resources"`
	Projects  []BaseProject  `json:"projects"`
}

func openInputJson() ([]byte, error) {
	jsonFile, err := os.Open("../input.json")
	if err != nil {
		log.Fatalf("Error opening JSON file: %s", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	return byteValue, err
}

func GetInputData(b *testing.B) TestInput {
	byteValue, err := openInputJson()
	if err != nil {
		b.Fatalf("Error reading test JSON file: %s", err)
	}

	var data TestInput
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		b.Fatalf("Error deserializing test JSON file: %s", err)
	}

	return data
}

func OpenInputData() (*TestInput, error) {
	byteValue, err := openInputJson()
	if err != nil {
		log.Fatalf("Error reading test JSON file: %s", err)
		return nil, err
	}

	var data TestInput
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("Error deserializing test JSON file: %s", err)
		return nil, err
	}

	return &data, nil
}

func Cast[T any](input any) (output T, err error) {
	// Serialize the input to JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return output, err
	}

	// Deserialize the JSON back to the specified type T
	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		return output, err
	}

	return output, nil
}

func CompareObjectsAsJSON(obj1, obj2 interface{}) error {
	// Convert both objects to JSON
	json1, err := json.Marshal(obj1)
	if err != nil {
		return fmt.Errorf("erro ao converter obj1 para JSON: %v", err)
	}

	json2, err := json.Marshal(obj2)
	if err != nil {
		return fmt.Errorf("erro ao converter obj2 para JSON: %v", err)
	}

	// Count the frequency of each byte in both json arrays
	freq1 := make(map[byte]int)
	freq2 := make(map[byte]int)

	for _, b := range json1 {
		freq1[b]++
	}
	for _, b := range json2 {
		freq2[b]++
	}

	// Check if frequencies are the same
	if !reflect.DeepEqual(freq1, freq2) {
		return errors.New("arrays have different byte frequencies and cannot be reordered")
	}

	return nil
}
