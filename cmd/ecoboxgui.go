//go:build !return

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	ep9000 "github.com/ladecadence/EP9000"
	"github.com/ladecadence/EcoBoxGUI/pkg/api"
	"github.com/ladecadence/EcoBoxGUI/pkg/appstate"
	"github.com/ladecadence/EcoBoxGUI/pkg/config"
	"github.com/ladecadence/EcoBoxGUI/pkg/door"
	"github.com/ladecadence/EcoBoxGUI/pkg/inventory"
	"github.com/ladecadence/EcoBoxGUI/pkg/led"
	"github.com/ladecadence/EcoBoxGUI/pkg/logging"
	"github.com/ladecadence/EcoBoxGUI/pkg/screens"
	"github.com/ladecadence/EcoBoxGUI/pkg/sound"
	r200 "github.com/ladecadence/GoR200"
)

const (
	QR_INIT_CABINET   = "****INIT CABINET****"
	QR_OPEN_DOOR      = "****OPEN DOOR****"
	ALARM_START_TIME  = 10000
	ALARM_REPEAT_TIME = 3000
)

func ReadAllTags(rfids []r200.R200) ([]string, error) {
	var tags []string
	for _, rfid := range rfids {
		responses, err := rfid.ReadTags()
		// we can have an error in one of the multiple reads but still get some data
		if err != nil && responses == nil {
			return nil, err
		}
		for _, r := range responses {
			if !slices.Contains(tags, hex.EncodeToString(r.EPC)) {
				fmt.Println("Tag: ", hex.EncodeToString(r.EPC))
				tags = append(tags, hex.EncodeToString(r.EPC))
			}
		}

		// do it for each reader after some delay to add new ones
		time.Sleep(1 * time.Second)
	}
	return tags, nil
}

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
	case appstate.StateChecked:
		mainContainer := screens.NewConfirmTuppers(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	case appstate.StateError:
		mainContainer := screens.NewError(a).Container
		fyne.Do(func() { main.SetContent(mainContainer) })
	}
}

func TestRead(appState *appstate.AppState, rfids []r200.R200, log logging.Logging) error {
	// RFID
	tags, err := ReadAllTags(rfids)
	if err != nil {
		fmt.Println("Error reading tags: ", err.Error())
		log.Log(logging.LogError, fmt.Sprintf("Error reading tags: %s", err))
		appState.SetState(appstate.StateError)
		return err
	}
	sort.Strings(tags)
	fmt.Printf("Tag (%d): ", len(tags))
	fmt.Println(tags)

	return nil
}

func main() {
	// start log
	log, err := logging.New("log")
	if err != nil {
		panic(err)
	}
	log.Log(logging.LogInfo, "Starting...")

	// read configuration
	config := config.Config{ConfFile: "config.toml"}
	config.GetConfig()

	// sound
	sound := sound.New()
	sound.Play()

	// get auth token
	token, err := api.GetToken()
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	fmt.Println(token)

	// appState
	appState := appstate.New("es", token)

	// database
	invent, err := inventory.New(config.Database)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}

	// QR Scanner
	qrData := make(chan []uint8)
	scanner, err := ep9000.New(config.QRPort, 115200)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}

	// RFID readers
	rfid, err := r200.New(config.RFIDPort, 115200, 25, false)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	defer rfid.Close()
	// configure RFID demodulator
	err = rfid.SendCommand(r200.CMD_SetReceiverDemodulatorParameters, []uint8{r200.MIX_Gain_3dB, r200.IF_AMP_Gain_40dB, 0x00, 0xB0})
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	rcv, err := rfid.Receive()
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	fmt.Printf("%v\n", rcv)

	// second reader
	rfid2, err := r200.New(config.RFIDPort2, 115200, 25, false)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	defer rfid.Close()
	// configure RFID demodulator
	err = rfid2.SendCommand(r200.CMD_SetReceiverDemodulatorParameters, []uint8{r200.MIX_Gain_3dB, r200.IF_AMP_Gain_40dB, 0x00, 0xB0})
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	rcv, err = rfid2.Receive()
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	fmt.Printf("%v\n", rcv)

	// Test reading all tags
	err = TestRead(appState, []r200.R200{rfid, rfid2}, log)
	if err != nil {
		log.Log(logging.LogError, err.Error())
	}

	// init state, clear local DB
	invent.DeleteAll()
	// get all containers from API for this cabinet
	containers, err := api.GetContainers(appState.Token(), config.Cabinet)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	// and store them in local DB
	for _, c := range containers {
		invent.InsertContainer(c)
	}

	// Door
	door, err := door.NewDoor(17, 27)
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}

	// Leds
	leds, err := led.NewLed("/dev/serial0")
	if err != nil {
		log.Log(logging.LogError, err.Error())
		panic(err)
	}
	leds.Normal()

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
				log.Log(logging.LogError, err.Error())
				panic(err)
			}
		}
	}()

	// check appState changes and hardware events
	go func() {
		for {
			if appState.StateChanged() {
				fmt.Printf("New state: %d\n", appState.State())
				log.Log(logging.LogInfo, fmt.Sprintf("New state: %d", appState.State()))
				appState.StateClean() // aknowledged
				switch appState.State() {
				case appstate.StateWelcome:
					// start, welcome screen and listen for QR code
					leds.Normal()
					appState.ClearUser()
					appState.DeleteContainers()
					ChangeScreen(appState, mainWindow)
					recv := <-qrData
					recv = bytes.Trim(recv, "\n\r")

					// check for special codes
					if bytes.Equal(recv, []byte(QR_INIT_CABINET+config.QRPass)) {
						// ok, init cabinet
						err := TestRead(appState, []r200.R200{rfid, rfid2}, log)
						if err != nil {

						}
						appState.SetState(appstate.StateWelcome)
						break
					}
					if bytes.Equal(recv, []byte(QR_OPEN_DOOR+config.QRPass)) {
						// ok, open door
						fmt.Println("Open door with special QR.")
						log.Log(logging.LogInfo, "Open door with special QR.")
						door.Open()
						appState.SetState(appstate.StateWelcome)
						break
					}

					// ok, no special codes, check user
					fmt.Printf("QR Data: %s\n", recv)
					log.Log(logging.LogInfo, fmt.Sprintf("QR Data: %s", recv))
					// get auth token
					token, err := api.GetToken()
					if err != nil {
						panic(err)
					}
					appState.SetToken(token)

					user, err := api.GetUser(appState.Token(), strings.TrimSpace(string(recv)), config.Cabinet)
					if err != nil {
						appState.SetState(appstate.StateUserError)
					} else {
						appState.SetUser(&user)
						appState.SetState(appstate.StateHello)
					}
				case appstate.StateHello:
					// hello and open door
					ChangeScreen(appState, mainWindow)
					time.Sleep(1 * time.Second)
					door.Open()
					// wait until the door is open
					for !door.IsOpen() {
						time.Sleep(10 * time.Millisecond)
					}
					leds.DoorOpen()
					// ok, change state
					appState.SetState(appstate.StateOpened)
				case appstate.StateUserError:
					ChangeScreen(appState, mainWindow)
				case appstate.StateOpened:
					api.Open(appState.Token(), appState.User().Name, config.Cabinet)
					ChangeScreen(appState, mainWindow)
					// ok wait for door to close, check if we need to play the alarm
					alarmTime := 0
					alarmStarted := false
					for door.IsOpen() {
						time.Sleep(10 * time.Millisecond)
						alarmTime += 10
						if (!alarmStarted) && alarmTime > ALARM_START_TIME {
							alarmTime = 0
							sound.Play()
							alarmStarted = true
							leds.Error()
						}
						if (alarmStarted) && alarmTime > ALARM_REPEAT_TIME {
							alarmTime = 0
							sound.Play()
						}
					}
					leds.Normal()
					// ok, now we need to make the inventory
					appState.SetState(appstate.StateClosed)
				case appstate.StateClosed:
					api.Close(appState.Token(), appState.User().Name, config.Cabinet)
					ChangeScreen(appState, mainWindow)
					// read tags
					tags, err := rfid.ReadTags()
					if err != nil {
						// RFID ERROR SCREEN?
					}
					// check database
					dbContainers, err := invent.GetContainers()
					if err != nil {
						// DB ERROR?
					}
					// remove the present containers so only removed tuppers remain
					for _, t := range tags {
						tag := hex.EncodeToString(t.EPC)
						fmt.Println("Tag:", tag)
						for i, container := range dbContainers {
							if container.Code == tag {
								dbContainers = slices.Delete(dbContainers, i, i+1)
							}
						}
					}
					fmt.Println(dbContainers)
					// add to state
					for _, t := range dbContainers {
						appState.AddContainer(t.Code)
					}
					// change state
					appState.SetState(appstate.StateChecked)
				case appstate.StateChecked:
					// ok, remove containers from inventory if neccesary
					if len(appState.ContainersTaken()) > 0 {
						for _, t := range appState.ContainersTaken() {
							invent.DeleteContainerByCode(t)
						}
						// and from API
						err = api.AdquireContainers(token, appState.User().Code, config.Cabinet, appState.ContainersTaken())
						if err != nil {
							log.Log(logging.LogError, fmt.Sprintf("Error with adquire API: %s", err))
							leds.Error()
							appState.SetState(appstate.StateError)
						}
					}
					ChangeScreen(appState, mainWindow)
				case appstate.StateError:
					ChangeScreen(appState, mainWindow)
				}
			}

		}
	}()

	mainWindow.ShowAndRun()
}
