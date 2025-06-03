package appstate

import (
	"fmt"
	"strconv"
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
	lang         string
	state        int
	user         *models.User
	token        *api.Token
	stateChanged bool
	tuppersTaken []int
}

func New(lang string) *AppState {
	a := AppState{lang: lang, state: StateWelcome, stateChanged: false, user: nil}
	a.tuppersTaken = []int{}
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

func (a *AppState) TuppersTaken() []int {
	return a.tuppersTaken
}

func (a *AppState) AddTupper(number int) {
	a.tuppersTaken = append(a.tuppersTaken, number)
}

func (a *AppState) DeleteTuppers() {
	a.tuppersTaken = []int{}
}

func (a *AppState) NumTuppers() int {
	return len(a.tuppersTaken)
}

func (a *AppState) TupperListString() string {
	list := ""
	for _, t := range a.tuppersTaken {
		list = list + strconv.Itoa(t) + ","
	}
	// remove last ","
	list = strings.TrimSuffix(list, ",")

	return list
}
