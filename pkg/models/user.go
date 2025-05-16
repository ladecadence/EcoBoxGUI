package models

type User struct {
	ID      string
	Name    string
	Surname string
	Tuppers []Tupper
}
