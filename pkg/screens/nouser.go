package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EcoBoxGUI/pkg/appstate"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type NoUser struct {
	Container  *fyne.Container
	labelError *canvas.Text
	labelMsg   *canvas.Text
	okButton   *widget.Button
	lang       string
	langBar    *components.LangBar
	setLang    func(string)
	state      *appstate.AppState
}

func NewNoUser(lang string, setlang func(string), a *appstate.AppState) *NoUser {
	h := NoUser{lang: lang, setLang: setlang, state: a}
	h.langBar = components.NewLangBar(h.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	h.labelError = canvas.NewText(languages.GetString("nouser.error", lang), theme.Color(theme.ColorNameForeground))
	h.labelError.TextSize = 20
	h.labelError.Alignment = fyne.TextAlignCenter
	h.labelMsg = canvas.NewText(languages.GetString("nouser.msg", lang), theme.Color(theme.ColorNameForeground))
	h.labelMsg.TextSize = 15
	h.labelMsg.Alignment = fyne.TextAlignCenter
	h.okButton = widget.NewButton(languages.GetString("nouser.button", lang), func() {
		h.state.SetState(appstate.StateWelcome)
	})
	hBox := container.NewHBox(layout.NewSpacer(), h.okButton, layout.NewSpacer())
	vBox := container.NewVBox(logo, spacer, h.labelError, h.labelMsg, hBox)
	center := container.NewCenter(vBox)
	h.Container = container.NewBorder(nil, h.langBar.Container, nil, nil, center)

	return &h
}

func (h *NoUser) UpdateLanguage(lang string) {
	h.labelError.Text = languages.GetString("nouser.error", lang)
	h.labelMsg.Text = languages.GetString("nouser.msg", lang)
	h.okButton.SetText(languages.GetString("nouser.button", lang))
	h.setLang(lang)
}
