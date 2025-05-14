package door

import (
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Door struct {
	openPin   gpio.PinOut
	closedPin gpio.PinIn
}

func NewDoor(pinOpen, pinClosed int) (*Door, error) {
	door := Door{}
	door.openPin = gpioreg.ByName(fmt.Sprintf("%d", pinOpen))
	if door.openPin == nil {
		return nil, errors.New("Can't use open pin")
	}
	err := door.openPin.Out(gpio.Low)
	if err != nil {
		return nil, err
	}

	door.closedPin = gpioreg.ByName(fmt.Sprintf("%d", pinClosed))
	if door.closedPin == nil {
		return nil, errors.New("Can't use closed pin")
	}
	err = door.closedPin.In(gpio.PullUp, gpio.FallingEdge)
	if err != nil {
		return nil, err
	}

	return &door, nil
}

func (d *Door) IsOpen() bool {
	return bool(d.closedPin.Read())
}

func (d *Door) Open() {
	d.openPin.Out(gpio.High)
	time.Sleep(100 * time.Millisecond)
	d.openPin.Out(gpio.Low)
}
