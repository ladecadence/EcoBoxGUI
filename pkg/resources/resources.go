package resources

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed es.svg
var esFlag []byte
var EsFlagRes = fyne.NewStaticResource("esflag", esFlag)

//go:embed es-pv.svg
var eusFlag []byte
var EusFlagRes = fyne.NewStaticResource("eusflag", eusFlag)

//go:embed gb.svg
var enFlag []byte
var EnFlagRes = fyne.NewStaticResource("enflag", enFlag)
