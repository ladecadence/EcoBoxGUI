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
	"github.com/ladecadence/EcoBoxGUI/pkg/door"
	"github.com/ladecadence/EcoBoxGUI/pkg/screens"
)

func ChangeScreen(a *appstate.AppState, main fyne.Window) {
	switch a.State() {
	case appstate.StateWelcome:
		mainContainer := screens.NewWelcome(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	case appstate.StateHello:
		mainContainer := screens.NewHello(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	case appstate.StateUserError:
		mainContainer := screens.NewNoUser(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	case appstate.StateOpened:
		mainContainer := screens.NewDoorOpen(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	case appstate.StateClosed:
		mainContainer := screens.NewDoorClosed(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	}
}

func main() {
	appState := appstate.New("es")

	// QR Scanner
	qrData := make(chan []uint8)
	scanner, err := ep9000.New("/dev/ttyACM0", 115200)
	if err != nil {
		panic(err)
	}

	// Door
	door, err := door.NewDoor(17, 27)
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

	// init state
	appState.SetState(appstate.StateWelcome)

	// listen to scanner thread
	go func() {
		for {
			err := scanner.Listen(qrData)
			if err != nil {
				panic(err)
			}
		}
	}()

	// check appState changes and hardware events
	go func() {
		for {
			if appState.StateChanged() {
				fmt.Printf("New state: %d\n", appState.State())
				appState.StateClean() // aknowledged
				switch appState.State() {
				case appstate.StateWelcome:
					// start, welcome screen and listen for QR code
					ChangeScreen(appState, mainWindow)
					recv := <-qrData
					fmt.Printf("QR Data: %s\n", recv)
					user, err := api.GetUser(strings.TrimSpace(string(recv)))
					if err != nil {
						appState.SetState(appstate.StateUserError)
					} else {
						appState.SetUser(&user)
						appState.SetState(appstate.StateHello)
					}
				case appstate.StateHello:
					// hello and open door
					ChangeScreen(appState, mainWindow)
					time.Sleep(3 * time.Second)
					door.Open()
					// wait until the door is open
					for !door.IsOpen() {
						time.Sleep(10 * time.Millisecond)
					}
					// ok, change state
					appState.SetState(appstate.StateOpened)
				case appstate.StateUserError:
					ChangeScreen(appState, mainWindow)
				case appstate.StateOpened:
					ChangeScreen(appState, mainWindow)
					// ok wait for door to close
					for door.IsOpen() {
						time.Sleep(10 * time.Millisecond)
					}
					// ok, now we need to make the inventory
					appState.SetState(appstate.StateClosed)
				case appstate.StateClosed:
					ChangeScreen(appState, mainWindow)
				}
			}

		}
	}()

	mainWindow.ShowAndRun()
}
