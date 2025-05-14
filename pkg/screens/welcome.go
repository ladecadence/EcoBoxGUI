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

type Welcome struct {
	Container    *fyne.Container
	labelWelcome *canvas.Text
	labelInfo    *canvas.Text
	state        *appstate.AppState
	langBar      *components.LangBar
}

func NewWelcome(a *appstate.AppState) *Welcome {
	w := Welcome{state: a}
	w.langBar = components.NewLangBar(w.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	w.labelWelcome = canvas.NewText(languages.GetString("welcome.welcome", w.state.Lang()), theme.Color(theme.ColorNameForeground))
	w.labelWelcome.TextSize = 20
	w.labelWelcome.Alignment = fyne.TextAlignCenter
	w.labelInfo = canvas.NewText(languages.GetString("welcome.info", w.state.Lang()), theme.Color(theme.ColorNameForeground))
	w.labelInfo.TextSize = 15
	w.labelInfo.Alignment = fyne.TextAlignCenter
	vBox := container.NewVBox(logo, spacer, w.labelWelcome, w.labelInfo)
	center := container.NewCenter(vBox)
	w.Container = container.NewBorder(nil, w.langBar.Container, nil, nil, center)

	return &w
}

func (w *Welcome) UpdateLanguage(lang string) {
	w.state.SetLang(lang)
	w.labelWelcome.Text = languages.GetString("welcome.welcome", w.state.Lang())
	w.labelInfo.Text = languages.GetString("welcome.info", w.state.Lang())
}
