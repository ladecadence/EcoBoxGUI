package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
)

func main() {
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	mainWindow := app.NewWindow("EcoBox")
	mainWindow.SetFullScreen(true)

	centerContainer := container.NewCenter(components.NewWelcome("es").Container)
	mainWindow.SetContent(centerContainer)
	mainWindow.ShowAndRun()
}
