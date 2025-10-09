package appstate

import (
	"fmt"
	"strings"

	"github.com/ladecadence/EcoBoxGUI/pkg/api"
	"github.com/ladecadence/EcoBoxGUI/pkg/models"
)

const (
	StateWelcome = iota
	StateHello
	StateUserError
	StateOpened
	StateClosed
	StateChecked
	StateFinish
	StateError
)

type AppState struct {
	lang            string
	state           int
	user            *models.User
	token           *api.Token
	stateChanged    bool
	containersTaken []string
	lastError       string
}

func New(lang string, token *api.Token) *AppState {
	a := AppState{lang: lang, token: token, state: StateWelcome, stateChanged: false, user: nil}
	a.containersTaken = []string{}
	a.lastError = ""
	return &a
}

func (a *AppState) Token() *api.Token {
	return a.token
}

func (a *AppState) SetToken(t *api.Token) {
	a.token = t
}

func (a *AppState) SetLang(lang string) {
	a.lang = lang
}

func (a *AppState) Lang() string {
	return a.lang
}

func (a *AppState) SetState(s int) {
	a.state = s
	a.stateChanged = true
	fmt.Println("State changed", a.state)
}

func (a *AppState) State() int {
	return a.state
}

func (a *AppState) StateChanged() bool {
	return a.stateChanged
}

func (a *AppState) StateClean() {
	a.stateChanged = false
}

func (a *AppState) SetUser(u *models.User) {
	a.user = u
}

func (a *AppState) User() *models.User {
	return a.user
}

func (a *AppState) Error() string {
	return a.lastError
}

func (a *AppState) SetError(e string) {
	a.lastError = e
}

func (a *AppState) ClearUser() {
	a.user = nil
}

func (a *AppState) ContainersTaken() []string {
	return a.containersTaken
}

func (a *AppState) AddContainer(code string) {
	a.containersTaken = append(a.containersTaken, code)
}

func (a *AppState) DeleteContainers() {
	a.containersTaken = []string{}
}

func (a *AppState) NumContainers() int {
	return len(a.containersTaken)
}

func (a *AppState) ContainerListString() string {
	list := ""
	for _, t := range a.containersTaken {
		list = list + t + ","
	}
	// remove last ","
	list = strings.TrimSuffix(list, ",")

	return list
}
