package entities

import "time"

type Class struct {
	Id      int      `db:"ID" json:"id"`
	Name    string   `db:"NAME" json:"name"`
	Objects []Object `json:"objects"`
}

type Object struct {
	Id       int        `db:"ID" json:"id"`
	Name     string     `db:"NAME" json:"name"`
	DateTime *time.Time `db:"DATETIME" json:"dateTime"`
	Value    float64    `db:"VALUE" json:"value"`
	Items    []Item     `json:"items"`
}

type Item struct {
	Id       int        `db:"ID" json:"id"`
	Name     string     `db:"NAME" json:"name"`
	DateTime *time.Time `db:"DATETIME" json:"dateTime"`
	Value    float64    `db:"VALUE" json:"value"`
}
