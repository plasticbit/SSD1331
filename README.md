[![Go Reference](https://pkg.go.dev/badge/github.com/BinaryDolphin29/SSD1331.svg)](https://pkg.go.dev/github.com/BinaryDolphin29/SSD1331)
# SSD1331
This library is a SSD1331 driver for Raspberry Pi 4.
See a Document at https://pkg.go.dev/github.com/BinaryDolphin29/SSD1331  
[![Test Pattern](https://github.com/BinaryDolphin29/SSD1331/blob/master/image/image.JPG)](https://github.com/BinaryDolphin29/SSD1331/tree/master/example/testPattern)

# Stable??
no...this library is wrote by begginer. **So I can't take any responsibility.**

# Hardware
I use this OLED module https://akizukidenshi.com/catalog/g/gP-14435/, and Raspberry Pi 4(Raspbian Lite 32bit). I tested only these device.

# Install
```
go get github.com/BinaryDolphin29/SSD1331
```

# Example
```go
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
		Frequency: 8 * physic.MegaHertz,
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
```

# About Library used
```
periph.io/x/conn/v3/driver/driverreg  
periph.io/x/conn/v3/gpio  
periph.io/x/conn/v3/physic  
periph.io/x/conn/v3/spi  
periph.io/x/conn/v3/spi/spireg  
periph.io/x/host/v3
```
Thanks!  
Library for more information [periph.io](https://periph.io/) and GitHub repository [github.com/google/periph](https://github.com/google/periph).

# Discord
If you can speak Japanese, I created this [Discord server](https://discord.gg/r2q4q8R5b8), join us!
