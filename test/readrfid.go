package main

import (
	"encoding/hex"
	"fmt"

	r200 "github.com/ladecadence/GoR200"
)

func main() {
	rfid, err := r200.New("/dev/ttyUSB0", 115200, 10, false)
	if err != nil {
		panic(err)
	}
	data, err := rfid.ReadTags()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, d := range data {
		fmt.Printf("\tPC: 0x%0x\n", d.PC)
		fmt.Printf("\tEPC: %s\n", hex.EncodeToString(d.EPC))
		fmt.Printf("\tCRC: 0x%0x\n", d.CRC)
	}
}
