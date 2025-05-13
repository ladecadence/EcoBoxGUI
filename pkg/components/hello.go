package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type Hello struct {
	Container  *fyne.Container
	labelHello *canvas.Text
	labelOpen  *canvas.Text
	name       string
	lang       string
	langBar    *LangBar
	setLang    func(string)
}

func NewHello(lang string, name string, setlang func(string)) *Hello {
	h := Hello{name: name, lang: lang, setLang: setlang}
	h.langBar = NewLangBar(h.UpdateLanguage)
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	h.labelHello = canvas.NewText(languages.GetString("hello.hello", lang)+" "+h.name, theme.Color(theme.ColorNameForeground))
	h.labelHello.TextSize = 20
	h.labelHello.Alignment = fyne.TextAlignCenter
	h.labelOpen = canvas.NewText(languages.GetString("hello.open", lang), theme.Color(theme.ColorNameForeground))
	h.labelOpen.TextSize = 15
	h.labelOpen.Alignment = fyne.TextAlignCenter
	vBox := container.NewVBox(logo, spacer, h.labelHello, h.labelOpen)
	center := container.NewCenter(vBox)
	h.Container = container.NewBorder(nil, h.langBar.Container, nil, nil, center)

	return &h
}

func (h *Hello) UpdateLanguage(lang string) {
	h.labelHello.Text = languages.GetString("hello.hello", lang) + " " + h.name
	h.labelOpen.Text = languages.GetString("hello.open", lang)
	h.setLang(lang)
}
