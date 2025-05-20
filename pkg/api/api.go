package api

import (
	"errors"

	"github.com/ladecadence/EcoBoxGUI/pkg/models"
)

var testUsers = []models.User{
	{ID: "d343b80c-ae82-41f5-ad77-7a32d1be85e2", Name: "Perico", Surname: "De los palotes", Tuppers: []models.Tupper{{ID: "tuper1", Number: 5}}},
	{ID: "898ceaf8-2b51-4a4b-8055-d04384620dc9", Name: "Manolo", Surname: "El del bombo"},
}

var testTuppers = []models.Tupper{
	{ID: "300833B2DDD9014000000001", Number: 1},
	{ID: "300833B2DDD9014000000002", Number: 2},
	{ID: "300833B2DDD9014000000003", Number: 3},
	{ID: "300833B2DDD9014000000004", Number: 4},
	{ID: "300833B2DDD9014000000005", Number: 5},
}

func GetUser(id string) (models.User, error) {
	for _, u := range testUsers {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, errors.New("no such user")
}

func GetTupper(id string) (models.Tupper, error) {
	for _, t := range testTuppers {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Tupper{}, errors.New("no such tupper")
}
