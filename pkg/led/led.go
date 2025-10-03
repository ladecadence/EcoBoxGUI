package led

import (
	pixarray "github.com/SimonWaldherr/ws2812/pixarray"
)

type Led struct {
	pin      int
	strip    pixarray.LEDStrip
	ledArray *pixarray.PixArray
}

func NewLed(pin int) (*Led, error) {
	led := Led{}

	var err error
	led.strip, err = pixarray.NewWS281x(10, 3, pixarray.StringOrders["GRB"], 800000, 10, []int{pin})
	if err != nil {
		return nil, err
	}
	led.ledArray = pixarray.NewPixArray(10, 3, led.strip)
	var p pixarray.Pixel
	p.R = 200
	p.G = 0
	p.B = 255
	led.ledArray.SetAll(p)
	led.ledArray.Write()

	return &led, nil
}

func (l *Led) Normal() {
	var p pixarray.Pixel
	p.R = 200
	p.G = 0
	p.B = 255
	l.ledArray.SetAll(p)
	l.ledArray.Write()
}

func (l *Led) DoorOpen() {
	var p pixarray.Pixel
	p.R = 0
	p.G = 255
	p.B = 0
	l.ledArray.SetAll(p)
	l.ledArray.Write()
}

func (l *Led) Error() {
	var p pixarray.Pixel
	p.R = 255
	p.G = 0
	p.B = 0
	l.ledArray.SetAll(p)
	l.ledArray.Write()
}
