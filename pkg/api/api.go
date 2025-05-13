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
	{ID: "d343b80c-ae82-41f5-ad77-7a32d1be85e2", Name: "Perico", Surname: "De los palotes", Tuppers: []Tupper{{ID: "tuper1", Model: 0}}},
	{ID: "898ceaf8-2b51-4a4b-8055-d04384620dc9", Name: "Manolo", Surname: "El del bombo"},
}

func GetUser(id string) (User, error) {
	for _, u := range testUsers {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("no such user")
}
