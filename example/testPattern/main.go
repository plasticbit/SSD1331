package main

import (
    "log"

	OLED "github.com/BinaryDolphin29/SSD1331"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3/rpi"
)

type pattColor struct {
	R, G, B int
}

const max75 int = (2<<7)*0.75 - 1

var testPattern = [6]pattColor{
	{max75, max75, max75}, // white
	{max75, max75, 0},     // yellow
	{0, max75, max75},     // cyan
	{0, max75, 0},         // green
	{max75, 0, max75},     // magenta
	{0, 0, max75},         // blue
}

func main() {
	display := &OLED.SSD1331{
		Name:      "/dev/spidev0.0",
		Frequency: 10 * physic.MegaHertz,
		ResetPin:  rpi.P1_22,
		DCPin:     rpi.P1_18,
		CSPin:     rpi.P1_24,
	}

	if err := display.Init(); err != nil {
		log.Fatalln(err.Error())
	}

    width, height := display.Resolution()

	defer display.Close()
	display.ClearDisplay()

	var (
		rowSize = width / len(testPattern)
		maxX    = rowSize
		x       = 0
	)

	for _, c := range testPattern {
		for ; x < maxX; x++ {
			for y := 0; y < height; y++ {
				display.SetPixel(x, y, c.R, c.G, c.B)
			}
		}

		maxX += rowSize
	}

    display.Display()
}