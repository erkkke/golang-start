package models

type Category struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CategoriesFilter struct {
	Query *string `json:"query"`

}