package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/resources"
)

// custom theme for language icons
type LangBarTheme struct{}

func (m LangBarTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}
func (m LangBarTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (m LangBarTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m LangBarTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInlineIcon {
		return 30
	}
	return theme.DefaultTheme().Size(name)
}

type LangBar struct {
	Lang      string
	Container *container.ThemeOverride
}

func NewLangBar(callback func(string)) *LangBar {
	l := LangBar{Lang: "es"}
	butES := NewLangButton("", resources.EsFlagRes, func() {
		l.Lang = "es"
		callback(l.Lang)
	})
	butEUS := NewLangButton("", resources.EusFlagRes, func() {
		l.Lang = "eus"
		callback(l.Lang)
	})
	butEN := NewLangButton("", resources.EnFlagRes, func() {
		l.Lang = "en"
		callback(l.Lang)
	})

	hBox := container.NewHBox(layout.NewSpacer(), butEUS, butES, butEN)
	l.Container = container.NewThemeOverride(hBox, LangBarTheme{})

	return &l
}
