package models

type User struct {
	Response
	ID   int    `json:"id"`
	Name string `json:"nombre"`
	Code string
}
