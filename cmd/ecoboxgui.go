package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	ep9000 "github.com/ladecadence/EP9000"
	"github.com/ladecadence/EcoBoxGUI/pkg/api"
	"github.com/ladecadence/EcoBoxGUI/pkg/appstate"
	"github.com/ladecadence/EcoBoxGUI/pkg/components"
)

func main() {
	appState := appstate.New("es")

	// QR Scanner
	qrData := make(chan []uint8)
	scanner, err := ep9000.New("/dev/ttyACM0", 115200)
	if err != nil {
		panic(err)
	}

	// GUI
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	mainWindow := app.NewWindow("EcoBox")
	mainWindow.SetFullScreen(true)
	ctrlQ := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
	mainWindow.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		mainWindow.Close()
	})

	mainContainer := components.NewWelcome(appState.Lang(), appState.SetLang).Container
	mainWindow.SetContent(mainContainer)

	// listen to scanner
	go func() {
		for {
			err := scanner.Listen(qrData)
			if err != nil {
				panic(err)
			}
		}
	}()

	// check GUI state and hardware events
	go func() {
		for {
			// AppState
			if appState.StateChanged() {
				fmt.Printf("New state: %d\n", appState.State())
				switch appState.State() {
				case appstate.StateWelcome:
					appState.StateClean() // aknowledged
					mainContainer := components.NewWelcome(appState.Lang(), appState.SetLang).Container
					fyne.Do(func() { mainWindow.SetContent(mainContainer) })
				}
			}

			// QR Scanner data
			select {
			case recv := <-qrData:
				fmt.Printf("QR Data: %s\n", recv)
				if appState.State() == appstate.StateWelcome {
					user, err := api.GetUser(strings.TrimSpace(string(recv)))
					if err != nil {
						mainContainer := components.NewNoUser(appState.Lang(), appState.SetLang, appState).Container
						fyne.Do(func() { mainWindow.SetContent(mainContainer) })
					} else {
						mainContainer := components.NewHello(appState.Lang(), user.Name, appState.SetLang).Container
						fyne.Do(func() { mainWindow.SetContent(mainContainer) })
					}
				}
			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	mainWindow.ShowAndRun()
}
