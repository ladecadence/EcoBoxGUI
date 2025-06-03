package api

import (
	"bytes"
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
	{ID: 0, Name: "Perico"},
	{ID: 1, Name: "Manolo"},
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

type UserRequest struct {
	User    string `json:"usuario"`
	Cabinet string `json:"armario"`
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

func GetUser(token *Token, id string, cabinet string) (models.User, error) {
	// url
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = userPath

	// body
	userRequest := UserRequest{User: id, Cabinet: cabinet}
	jsonBody, err := json.Marshal(userRequest)
	if err != nil {
		return models.User{}, errors.New("Problem encoding user request")
	}

	// make request
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.User{}, errors.New("Can't create user request")
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(r)
	if err != nil {
		return models.User{}, errors.New("Can't execute user request")
	}

	// parse response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.User{}, errors.New("Can't parse user response body")
	}
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		return models.User{}, errors.New("Can't parse user response json")
	}

	// check response
	if user.Code == 1 {
		return user, nil
	} else {
		return models.User{}, errors.New("no such user")
	}
}

func GetTupper(id string) (models.Tupper, error) {
	for _, t := range testTuppers {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Tupper{}, errors.New("no such tupper")
}
