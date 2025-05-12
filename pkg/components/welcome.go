package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/languages"
)

type Welcome struct {
	Container *fyne.Container
}

func NewWelcome(lang string) *Welcome {
	w := Welcome{}
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.NewSize(1000, 100))
	labelWelcome := canvas.NewText(languages.GetString("welcome.welcome", lang), theme.Color(theme.ColorNameForeground))
	labelWelcome.TextSize = 50
	labelWelcome.Alignment = fyne.TextAlignCenter
	labelInfo := canvas.NewText(languages.GetString("welcome.info", lang), theme.Color(theme.ColorNameForeground))
	labelInfo.TextSize = 30
	labelInfo.Alignment = fyne.TextAlignCenter
	w.Container = container.NewVBox(logo, spacer, labelWelcome, labelInfo)

	return &w
}
