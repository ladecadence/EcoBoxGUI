package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type Welcome struct {
	Container    *fyne.Container
	labelWelcome *canvas.Text
	labelInfo    *canvas.Text
	lang         string
	langBar      *components.LangBar
	setLang      func(string)
}

func NewWelcome(lang string, setlang func(string)) *Welcome {
	w := Welcome{lang: lang, setLang: setlang}
	w.langBar = components.NewLangBar(w.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	w.labelWelcome = canvas.NewText(languages.GetString("welcome.welcome", lang), theme.Color(theme.ColorNameForeground))
	w.labelWelcome.TextSize = 20
	w.labelWelcome.Alignment = fyne.TextAlignCenter
	w.labelInfo = canvas.NewText(languages.GetString("welcome.info", lang), theme.Color(theme.ColorNameForeground))
	w.labelInfo.TextSize = 15
	w.labelInfo.Alignment = fyne.TextAlignCenter
	vBox := container.NewVBox(logo, spacer, w.labelWelcome, w.labelInfo)
	center := container.NewCenter(vBox)
	w.Container = container.NewBorder(nil, w.langBar.Container, nil, nil, center)

	return &w
}

func (w *Welcome) UpdateLanguage(lang string) {
	w.labelWelcome.Text = languages.GetString("welcome.welcome", lang)
	w.labelInfo.Text = languages.GetString("welcome.info", lang)
	w.setLang(lang)
}
