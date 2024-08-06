package entities

import "time"

type Class struct {
	ID      int      `gorm:"primaryKey" json:"id"`
	Name    string   `json:"name"`
	Objects []Object `gorm:"foreignKey:ClassID" json:"objects"`
}

type Object struct {
	ID       int        `gorm:"primaryKey" json:"id"`
	Name     string     `json:"name"`
	Value    float64    `json:"value"`
	Datetime *time.Time `json:"dateTime"`
	ClassID  int        `json:"-"`
	Items    []Item     `gorm:"many2many:object_item_link;" json:"items"`
}

type Item struct {
	ID       int        `gorm:"primaryKey" json:"id"`
	Name     string     `json:"name"`
	Value    float64    `json:"value"`
	Datetime *time.Time `json:"dateTime"`
}
