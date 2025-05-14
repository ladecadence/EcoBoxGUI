package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/appstate"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type DoorOpen struct {
	Container *fyne.Container
	labelOpen *canvas.Text
	labelMsg  *canvas.Text
	state     *appstate.AppState
	langBar   *components.LangBar
}

func NewDoorOpen(a *appstate.AppState) *DoorOpen {
	h := DoorOpen{state: a}
	h.langBar = components.NewLangBar(h.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	h.labelOpen = canvas.NewText(languages.GetString("open.open", a.Lang())+" "+h.state.User().Name, theme.Color(theme.ColorNameForeground))
	h.labelOpen.TextSize = 20
	h.labelOpen.Alignment = fyne.TextAlignCenter
	h.labelMsg = canvas.NewText(languages.GetString("open.msg", a.Lang()), theme.Color(theme.ColorNameForeground))
	h.labelMsg.TextSize = 15
	h.labelMsg.Alignment = fyne.TextAlignCenter
	vBox := container.NewVBox(logo, spacer, h.labelOpen, h.labelMsg)
	center := container.NewCenter(vBox)
	h.Container = container.NewBorder(nil, h.langBar.Container, nil, nil, center)

	return &h
}

func (h *DoorOpen) UpdateLanguage(lang string) {
	h.state.SetLang(lang)
	h.labelOpen.Text = languages.GetString("open.open", h.state.Lang()) + " " + h.state.User().Name
	h.labelMsg.Text = languages.GetString("open.msg", lang)
}
