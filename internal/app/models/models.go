package models

type Book struct {
	ID     string `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Author string `json:"author" db:"author"`
	Price  uint   `json:"price" db:"price"`
}
