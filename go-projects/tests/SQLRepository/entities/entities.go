package entities

import "time"

// Class representa a tabela CLASSES
type Class struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
}

// Implementação da interface Entity para Class
func (c *Class) TableName() string {
	return "classes"
}

func (c *Class) ColumnsNames() []string {
	return []string{"id", "name"}
}

func (c *Class) Fields() []interface{} {
	return []interface{}{&c.Id, &c.Name}
}

// Object representa a tabela OBJECTS
type Object struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	DateTime *time.Time `json:"dateTime"`
	Value    float64    `json:"value"`
	ClassID  int        `json:"classId"`
	Items    []Item     `json:"items"`
}

// Implementação da interface Entity para Object
func (o *Object) TableName() string {
	return "objects"
}

func (o *Object) ColumnsNames() []string {
	return []string{"id", "name", "value", "datetime", "class_id"}
}

func (o *Object) Fields() []interface{} {
	return []interface{}{&o.Id, &o.Name, &o.Value, &o.DateTime, &o.ClassID}
}

// Item representa a tabela ITEMS
type Item struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	DateTime *time.Time `json:"dateTime"`
	Value    float64    `json:"value"`
}

// Implementação da interface Entity para Item
func (i *Item) TableName() string {
	return "items"
}

func (i *Item) ColumnsNames() []string {
	return []string{"id", "name", "value", "datetime"}
}

func (i *Item) Fields() []interface{} {
	return []interface{}{&i.Id, &i.Name, &i.Value, &i.DateTime}
}
