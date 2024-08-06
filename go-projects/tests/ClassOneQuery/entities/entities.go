package entities

import "time"

type Class struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
}

type Object struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	DateTime *time.Time `json:"dateTime"`
	Value    float64    `json:"value"`
	Items    []Item     `json:"items"`
}

type Item struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	DateTime *time.Time `json:"dateTime"`
	Value    float64    `json:"value"`
}
