package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EcoBoxGUI/pkg/appstate"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type DebugInfo struct {
	Container *fyne.Container
	textDebug *widget.TextGrid
	okButton  *widget.Button
	langBar   *components.LangBar
	state     *appstate.AppState
}

func NewDebugInfo(a *appstate.AppState) *DebugInfo {
	c := DebugInfo{state: a}
	c.langBar = components.NewLangBar(c.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	c.textDebug = widget.NewTextGrid()
	c.textDebug.Append(a.Error())
	c.okButton = widget.NewButton(languages.GetString("debug.button", c.state.Lang()), func() {
		// clear error?
		a.SetError("")
		c.state.SetState(appstate.StateWelcome)
	})
	hBox := container.NewHBox(layout.NewSpacer(), c.okButton, layout.NewSpacer())
	vBox := container.NewVBox(logo, spacer, c.textDebug, hBox)
	center := container.NewCenter(vBox)
	c.Container = container.NewBorder(nil, c.langBar.Container, nil, nil, center)

	return &c
}

func (c *DebugInfo) UpdateLanguage(lang string) {
	c.state.SetLang(lang)
	c.okButton.SetText(languages.GetString("debug.button", lang))
}
