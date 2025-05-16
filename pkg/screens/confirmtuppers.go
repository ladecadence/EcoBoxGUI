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

type ConfirmTuppers struct {
	Container    *fyne.Container
	labelConfirm *canvas.Text
	labelMsg     *canvas.Text
	okButton     *widget.Button
	langBar      *components.LangBar
	state        *appstate.AppState
}

func NewConfirmTuppers(a *appstate.AppState) *ConfirmTuppers {
	c := ConfirmTuppers{state: a}
	c.langBar = components.NewLangBar(c.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	c.labelConfirm = canvas.NewText(languages.GetString("confirm.confirm", c.state.Lang()), theme.Color(theme.ColorNameForeground))
	c.labelConfirm.TextSize = 20
	c.labelConfirm.Alignment = fyne.TextAlignCenter
	switch a.NumTuppers() {
	case 0:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.none", c.state.Lang()), theme.Color(theme.ColorNameForeground))
	case 1:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.sing"+a.TupperListString(), c.state.Lang()), theme.Color(theme.ColorNameForeground))
	default:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.plur"+a.TupperListString(), c.state.Lang()), theme.Color(theme.ColorNameForeground))
	}

	c.labelMsg.TextSize = 15
	c.labelMsg.Alignment = fyne.TextAlignCenter
	c.okButton = widget.NewButton(languages.GetString("confirm.button", c.state.Lang()), func() {
		c.state.SetState(appstate.StateWelcome)
	})
	hBox := container.NewHBox(layout.NewSpacer(), c.okButton, layout.NewSpacer())
	vBox := container.NewVBox(logo, spacer, c.labelConfirm, c.labelMsg, hBox)
	center := container.NewCenter(vBox)
	c.Container = container.NewBorder(nil, c.langBar.Container, nil, nil, center)

	return &c
}

func (c *ConfirmTuppers) UpdateLanguage(lang string) {
	c.state.SetLang(lang)
	c.labelConfirm.Text = languages.GetString("confirm.confirm", lang)
	switch c.state.NumTuppers() {
	case 0:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.none", c.state.Lang()), theme.Color(theme.ColorNameForeground))
	case 1:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.sing"+c.state.TupperListString(), c.state.Lang()), theme.Color(theme.ColorNameForeground))
	default:
		c.labelMsg = canvas.NewText(languages.GetString("confirm.msg.plur"+c.state.TupperListString(), c.state.Lang()), theme.Color(theme.ColorNameForeground))
	}
	c.okButton.SetText(languages.GetString("confirm.button", lang))
}
