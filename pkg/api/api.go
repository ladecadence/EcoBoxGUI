package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ladecadence/EcoBoxGUI/pkg/models"
)

const (
	apiURL      = "https://ecobox.clienteslotura.com"
	tokenPath   = "/api/oauth/token"
	userPath    = "/api/usuario"
	openPath    = "/api/armario/apertura"
	closePath   = "/api/armario/cierre"
	adquirePath = "/api/contenedor/adquisicion"
	returnPath  = "/api/contenedor/devolucion"
)

var testUsers = []models.User{
	{ID: "d343b80c-ae82-41f5-ad77-7a32d1be85e2", Name: "Perico", Surname: "De los palotes", Tuppers: []models.Tupper{{ID: "tuper1", Number: 5}}},
	{ID: "898ceaf8-2b51-4a4b-8055-d04384620dc9", Name: "Manolo", Surname: "El del bombo"},
}

var testTuppers = []models.Tupper{
	{ID: "300833b2ddd9014000000001", Number: 1},
	{ID: "300833b2ddd9014000000002", Number: 2},
	{ID: "300833b2ddd9014000000003", Number: 3},
	{ID: "300833b2ddd9014000000004", Number: 4},
	{ID: "300833b2ddd9014000000005", Number: 5},
}

type Token struct {
	Type        string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func GetToken() (*Token, error) {
	// Form URL Encoded data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "aZ3bL9mN2xQ7Rt1P")
	data.Set("client_secret", "xF7pL2aMwTqR98ZKhV1dBnEC4jsGUtyoXmPbQ3vJr0Ae6WSzgNYl9cTXkfO5udhB")

	// url
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = tokenPath

	// make request
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return nil, errors.New("Can't create token request")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.New("Can't execute token request")
	}

	// parse response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Can't parse token response body")
	}
	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, errors.New("Can't parse token response json")
	}

	return &token, nil
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
