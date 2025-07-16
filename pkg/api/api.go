package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ladecadence/EcoBoxGUI/pkg/models"
)

const (
	apiURL         = "https://ecobox.clienteslotura.com"
	tokenPath      = "/api/oauth/token"
	userPath       = "/api/usuario"
	openPath       = "/api/armario/apertura"
	closePath      = "/api/armario/cierre"
	adquirePath    = "/api/contenedor/adquisicion"
	returnPath     = "/api/contenedor/devolucion"
	containersPath = "/api/contenedor/"
)

var testUsers = []models.User{
	{ID: 0, Name: "Perico"},
	{ID: 1, Name: "Manolo"},
}

type Token struct {
	Type        string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type ApiRequest struct {
	User    string `json:"usuario"`
	Cabinet string `json:"armario"`
}

type ContainerRequest struct {
	Cabinet    string             `json:"armario"`
	Containers []models.Container `json:"contenedores"`
}

type InventoryRequest struct {
	User       string             `json:"usuario"`
	Cabinet    string             `json:"armario"`
	Containers []models.Container `json:"contenedores"`
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
	userRequest := ApiRequest{User: id, Cabinet: cabinet}
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

	fmt.Println(string(body))

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		return models.User{}, errors.New("Can't parse user response json")
	}

	// check response
	if user.Result == 1 {
		return user, nil
	} else {
		return models.User{}, errors.New("no such user")
	}
}

func Open(token *Token, id string, cabinet string) (models.Response, error) {
	// url
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = openPath

	// body
	userRequest := ApiRequest{User: id, Cabinet: cabinet}
	jsonBody, err := json.Marshal(userRequest)
	if err != nil {
		return models.Response{}, errors.New("Problem encoding open request")
	}
	// make request
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.Response{}, errors.New("Can't create open request")
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(r)
	if err != nil {
		return models.Response{}, errors.New("Can't execute open request")
	}

	// parse response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, errors.New("Can't parse open response body")
	}

	var user models.Response
	if err := json.Unmarshal(body, &user); err != nil {
		return models.Response{}, errors.New("Can't parse open response json")
	}

	// check response
	if user.Result == 1 {
		return user, nil
	} else {
		return models.Response{}, errors.New("Problem with open request")
	}

}

func Close(token *Token, id string, cabinet string) (models.Response, error) {
	// url
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = closePath

	// body
	userRequest := ApiRequest{User: id, Cabinet: cabinet}
	jsonBody, err := json.Marshal(userRequest)
	if err != nil {
		return models.Response{}, errors.New("Problem encoding close request")
	}
	// make request
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.Response{}, errors.New("Can't create close request")
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(r)
	if err != nil {
		return models.Response{}, errors.New("Can't execute close request")
	}

	// parse response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, errors.New("Can't parse close response body")
	}

	var user models.Response
	if err := json.Unmarshal(body, &user); err != nil {
		return models.Response{}, errors.New("Can't parse close response json")
	}

	// check response
	if user.Result == 1 {
		return user, nil
	} else {
		return models.Response{}, errors.New("Problem with open request")
	}
}

func GetContainers(token *Token, cabinet string) ([]models.Container, error) {
	// url
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = containersPath

	// body
	containerRequest := ContainerRequest{Cabinet: cabinet, Containers: []models.Container{}}
	jsonBody, err := json.Marshal(containerRequest)
	if err != nil {
		return nil, errors.New("Problem encoding container request")
	}
	// make request
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return []models.Container{}, errors.New("Can't create container request")
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(r)
	if err != nil {
		return []models.Container{}, errors.New("Can't execute container request")
	}

	// parse response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Container{}, errors.New("Can't parse container response body")
	}

	var containers models.Containers
	if err := json.Unmarshal(body, &containers); err != nil {
		return []models.Container{}, errors.New("Can't parse container response json")
	}

	// check response
	if containers.Result == 1 {
		return containers.Containers, nil
	} else {
		fmt.Println(containers)
		return []models.Container{}, errors.New("Problem with container request")
	}
}
