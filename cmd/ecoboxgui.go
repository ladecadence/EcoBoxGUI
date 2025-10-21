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
	QR_INIT_CABINET         = "****INIT CABINET****"
	QR_OPEN_DOOR            = "****OPEN DOOR****"
	ALARM_START_TIME        = 10000
	ALARM_REPEAT_TIME       = 3000
	ALARM_NOTIFICATION_TIME = 60000

	APP_ERROR_RFID = "0000"
	APP_ERROR_API  = "0001"
	APP_ERROR_QR   = "0002"
	APP_ERROR_DB   = "0003"
)

var (
	rfidGainConfig = []uint8{r200.MIX_Gain_3dB, r200.IF_AMP_Gain_36dB, 0x00, 0xA0}
)

func ReadAllTags(rfids []r200.R200) ([]string, error) {
	var tags []string
	for i, rfid := range rfids {
		fmt.Printf("Reading tags from RFID reader %d\n", i)
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
	log.Log(logging.LogInfo, "Obtained API token...")

	// appState
	appState := appstate.New("es", token)

	// database
	invent, err := inventory.New(config.Database)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error opening database: %s", err))
		panic(err)
	}

	// QR Scanner
	qrData := make(chan []uint8)
	scanner, err := ep9000.New(config.QRPort, 115200)
	if err != nil {
		// try to open next port
		scanner, err = ep9000.New("/dev/ttyACM1", 115200)
		if err != nil {
			log.Log(logging.LogError, fmt.Sprintf("Error opening QR scanner: %s", err))
			panic(err)
		}
	}

	// RFID readers
	rfid, err := r200.New(config.RFIDPort, 115200, 35, false)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error opening RFID reader 1: %s", err))
		panic(err)
	}
	defer rfid.Close()
	// query params
	err = rfid.SendCommand(r200.CMD_SetQueryParameters, []uint8{0x10, 0x18})
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error sending command to RFID 1: %s", err))
		panic(err)
	}
	rcv, err := rfid.Receive()
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error reading answer from RFID 1: %s", err))
		panic(err)
	}
	fmt.Printf("Query params 1: %v\n", rcv)
	// configure RFID demodulator
	err = rfid.SendCommand(r200.CMD_SetReceiverDemodulatorParameters, rfidGainConfig)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error sending command to RFID 1: %s", err))
		panic(err)
	}
	rcv, err = rfid.Receive()
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error reading answer from RFID 1: %s", err))
		panic(err)
	}
	fmt.Printf("Demodulator 1: %v\n", rcv)

	// second reader
	rfid2, err := r200.New(config.RFIDPort2, 115200, 35, false)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error opening RFID reader 2: %s", err))
		panic(err)
	}
	defer rfid.Close()
	// query params
	err = rfid2.SendCommand(r200.CMD_GetQueryParameters, []uint8{})
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error sending command to RFID 2: %s", err))
		panic(err)
	}
	rcv, err = rfid2.Receive()
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error reading answer from RFID 2: %s", err))
		panic(err)
	}
	fmt.Printf("Query params 2: %v\n", rcv)
	// configure RFID demodulator
	err = rfid2.SendCommand(r200.CMD_SetReceiverDemodulatorParameters, rfidGainConfig)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error sending command to RFID 2: %s", err))
		panic(err)
	}
	rcv, err = rfid2.Receive()
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error reading answer from RFID 2: %s", err))
		panic(err)
	}
	fmt.Printf("Demodulator 2: %v\n", rcv)

	// Test reading all tags
	err = TestRead(appState, []r200.R200{rfid, rfid2}, log)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error reading tags: %s", err))
	}

	// init state, clear local DB
	invent.DeleteAll()
	// get all containers from API for this cabinet
	containers, err := api.GetContainers(appState.Token(), config.Cabinet)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error getting containers from API: %s", err))
		panic(err)
	}
	// and store them in local DB
	for _, c := range containers {
		if c.Active == 1 && c.Available {
			invent.InsertContainer(c)
		}
	}
	// debug
	containers, err = invent.GetContainers()
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error getting containers from database: %s", err))
		panic(err)
	}
	list := "Initial containers: "
	for _, c := range containers {
		list += c.Code
		list += " "
	}
	log.Log(logging.LogData, list)

	// Door
	door, err := door.NewDoor(17, 27)
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error with door GPIO: %s", err))
		panic(err)
	}

	// Leds
	leds, err := led.NewLed("/dev/serial0")
	if err != nil {
		log.Log(logging.LogError, fmt.Sprintf("Error opening LEDs serial port: %s", err))
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
				log.Log(logging.LogError, fmt.Sprintf("Error launching QR reader thread: %s", err))
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
							log.Log(logging.LogError, "Error reading tags on init cabinet")
						}
						// ok, reinit DB
						invent.DeleteAll()
						// get current API container status
						// get auth token
						token, err = api.GetToken()
						if err != nil {
							log.Log(logging.LogError, err.Error())
							panic(err)
						}
						appState.SetToken(token)
						containers, err := api.GetContainers(appState.Token(), config.Cabinet)
						if err != nil {
							log.Log(logging.LogError, fmt.Sprintf("Error getting containers from API: %s", err))
							panic(err)
						}
						// and store them in local DB
						for _, c := range containers {
							if c.Active == 1 && c.Available {
								invent.InsertContainer(c)
							}
						}
						// ok, start again
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
					time.Sleep(1500 * time.Millisecond)
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
					api.Open(appState.Token(), appState.User().Code, config.Cabinet)
					ChangeScreen(appState, mainWindow)
					// ok wait for door to close, check if we need to play the alarm or send
					// the notification of open door
					alarmTime := 0
					notificationTime := 0
					alarmStarted := false
					notificationSent := false
					for door.IsOpen() {
						time.Sleep(10 * time.Millisecond)
						alarmTime += 10
						notificationTime += 10
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
						if (!notificationSent) && notificationTime > ALARM_NOTIFICATION_TIME {
							notificationSent = true
							api.DoorAlarm(appState.Token(), config.Cabinet)
						}
					}
					leds.Normal()
					// ok, now we need to make the inventory
					appState.SetState(appstate.StateClosed)
				case appstate.StateClosed:
					api.Close(appState.Token(), appState.User().Code, config.Cabinet)
					ChangeScreen(appState, mainWindow)
					// read tags
					tags, err := ReadAllTags([]r200.R200{rfid, rfid2})
					if err != nil {
						// RFID ERROR SCREEN
						log.Log(logging.LogError, fmt.Sprintf("Error with RFID: %s", err))
						leds.Error()
						appState.SetError(APP_ERROR_RFID)
						appState.SetState(appstate.StateError)
						break
					}
					// check database
					dbContainers, err := invent.GetContainers()
					if err != nil {
						log.Log(logging.LogError, fmt.Sprintf("Error with database: %s", err))
						leds.Error()
						appState.SetError(APP_ERROR_DB)
						appState.SetState(appstate.StateError)
						break
					}
					// remove the present containers so only removed tuppers remain
					for _, t := range tags {
						fmt.Println("Tag:", t)
						for i, container := range dbContainers {
							// ignore caps on hex digits
							if strings.EqualFold(container.Code, t) {
								dbContainers = slices.Delete(dbContainers, i, i+1)
							}
						}
					}
					fmt.Printf("Contenedores retirados: %v", dbContainers)
					log.Log(logging.LogInfo, fmt.Sprintf("Contenedores retirados: %v\n", dbContainers))
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
						// get auth token
						token, err := api.GetToken()
						if err != nil {
							log.Log(logging.LogError, fmt.Sprintf("Error with adquire API: %s", err))
							leds.Error()
							appState.SetError(APP_ERROR_API)
							appState.SetState(appstate.StateError)
						}
						appState.SetToken(token)
						err = api.AdquireContainers(token, appState.User().Code, config.Cabinet, appState.ContainersTaken())
						if err != nil {
							log.Log(logging.LogError, fmt.Sprintf("Error with adquire API: %s", err))
							leds.Error()
							appState.SetError(APP_ERROR_API)
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
