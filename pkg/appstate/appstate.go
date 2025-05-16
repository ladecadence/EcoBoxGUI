package appstate

import (
	"fmt"

	"github.com/ladecadence/EcoBoxGUI/pkg/models"
)

const (
	StateWelcome = iota
	StateHello
	StateUserError
	StateOpened
	StateClosed
	StateFinish
	StateError
)

type AppState struct {
	lang         string
	state        int
	user         *models.User
	stateChanged bool
}

func New(lang string) *AppState {
	a := AppState{lang: lang, state: StateWelcome, stateChanged: false, user: nil}
	return &a
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

func (a *AppState) ClearUser() {
	a.user = nil
}
