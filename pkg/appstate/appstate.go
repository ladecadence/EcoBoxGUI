package appstate

import "fmt"

const (
	StateWelcome = iota
	StateHello
	StateOpened
	StateClosed
	StateFinish
	StateError
)

type AppState struct {
	lang         string
	state        int
	stateChanged bool
}

func New(lang string) *AppState {
	a := AppState{lang: lang, state: StateWelcome, stateChanged: false}
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
