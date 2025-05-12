package api

import "errors"

type Tupper struct {
	ID    string
	Model int
}

type User struct {
	ID      string
	Name    string
	Surname string
	Tuppers []Tupper
}

var testUsers = []User{
	{ID: "1234", Name: "Perico", Surname: "De los palotes", Tuppers: []Tupper{Tupper{ID: "tuper1", Model: 0}}},
	{ID: "5678", Name: "Manolo", Surname: "El del bombo"},
}

func GetUser(id string) (User, error) {
	for _, u := range testUsers {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("no such user")
}
