//go:build return

package main

import (
	"encoding/hex"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EcoBoxGUI/pkg/api"
	"github.com/ladecadence/EcoBoxGUI/pkg/config"
	r200 "github.com/ladecadence/GoR200"
)

func ReadAllTags(rfid r200.R200) ([]string, error) {
	responses, err := rfid.ReadTags()
	// we can have an error in one of the multiple reads but still get some data
	if err != nil && responses == nil {
		return nil, err
	}
	var tags []string
	for _, r := range responses {
		fmt.Println("Tag: ", hex.EncodeToString(r.EPC))
		tags = append(tags, hex.EncodeToString(r.EPC))
	}
	return tags, nil
}

func main() {
	tagsRFID := binding.BindStringList(
		&[]string{},
	)

	// read configuration
	config := config.Config{ConfFile: "config.toml"}
	config.GetConfig()

	// get auth token
	token, err := api.GetToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	// RFID reader
	rfid, err := r200.New(config.RFIDPort, 115200, 10, false)
	if err != nil {
		panic(err)
	}

	// app
	app := app.New()
	app.Settings().SetTheme(theme.LightTheme())
	mainWindow := app.NewWindow("EcoBox Return Containers")

	ctrlQ := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
	mainWindow.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		mainWindow.Close()
	})

	// GUI
	logo := canvas.NewImageFromFile("res/ecobox.svg")
	logo.FillMode = canvas.ImageFillOriginal
	tagList := widget.NewListWithData(tagsRFID,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	readButton := widget.NewButton("Read RFID", func() {

		fmt.Println("Reading")
		tagsRFID.Set([]string{})
		tags, err := ReadAllTags(rfid)
		if err != nil {

		} else {
			tagsRFID.Set(tags)
		}
	})

	returnButton := widget.NewButton("Return Containers", func() {
		fmt.Println("Returning")
		tags, _ := tagsRFID.Get()
		fmt.Printf("%v,\n", tags)
	})
	vBox := container.NewVBox(logo, readButton)
	border := container.NewBorder(vBox, returnButton, nil, nil, tagList)
	mainWindow.SetContent(border)

	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.ShowAndRun()
}
