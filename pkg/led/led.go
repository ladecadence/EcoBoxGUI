package led

import (
	"go.bug.st/serial"
)

type Led struct {
	port serial.Port
}

func NewLed(portFile string) (*Led, error) {
	led := Led{}

	// prepare port
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	// open port
	var err error
	led.port, err = serial.Open(portFile, mode)
	if err != nil {
		return nil, err
	}

	return &led, nil
}

func (l *Led) Normal() {
	l.port.Write([]byte("N"))
}

func (l *Led) DoorOpen() {
	l.port.Write([]byte("O"))
}

func (l *Led) Error() {
	l.port.Write([]byte("E"))
}
