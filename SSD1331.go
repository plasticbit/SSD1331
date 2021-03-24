package OLED

import (
	"io"
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type ScrollStep byte

const (
	Frames6 ScrollStep = iota
	Frame10
	Frame100
	Frame200
)

type DisplayMode byte

const (
	Nomal DisplayMode = iota + 0xA4
	EntireON
	EntireOFF
	Inverse
)

const (
	high = gpio.High
	low  = gpio.Low

	width     = 96
	height    = 64
	bufferLen = (width * height) * 2
)

// SSD1331 96*64 pixels, 16bit color OLED.
type SSD1331 struct {
	Name      string
	Frequency physic.Frequency
	ResetPin  gpio.PinIO
	DCPin     gpio.PinIO
	CSPin     gpio.PinOut

	connect spi.Conn
	clorser io.Closer

	display bool
	buffer  []byte
}

// Init Initialize oled.
func (oled *SSD1331) Init() error {
	_, err := host.Init()
	if err != nil {
		return err
	}

	if _, err := driverreg.Init(); err != nil {
		return err
	}

	p, err := spireg.Open(oled.Name)
	if err != nil {
		return err
	}

	oled.clorser = p.(io.Closer)

	c, err := p.Connect(oled.Frequency, spi.Mode3, 8)
	if err != nil {
		return err
	}

	oled.connect = c

	oled.CSPin.Out(low)
	oled.ResetPin.Out(low)
	time.Sleep(50 * time.Millisecond)
	oled.ResetPin.Out(high)
	time.Sleep(50 * time.Millisecond)

	// send initial commands
	oled.sendCommand([]byte{
		0xAE,             // display off
		0x15, 0x00, 0x5F, // column addr
		0x75, 0x00, 0x3F, // column addr
		// 0x81, 0x91,
		// 0x82, 0x50,
		// 0x83, 0x7D,
		0x87, 0x0A, // master current
		// 0x8A, 0x64,
		// 0x8B, 0x78,
		// 0x8C, 0x64,

		0xA0, 0x78, // remap, color depth setting
		0xA1, 0x00,
		0xA2, 0x00,
		0xA4,
		0xA8, 0x3F,
		0xAD, 0x8E,
		0xB0, 0x0B,
		0xB1, 0x31,
		0xB3, 0xF0,
		0xBB, 0x3A,
		0xBE, 0x3E,
		0x87, 0x06,
		0xAF, // display on
	})

	oled.display = true
	oled.buffer = make([]byte, bufferLen)

	return nil
}

// Close Close the port.
func (oled *SSD1331) Close() error {
	return oled.clorser.Close()
}

// Resolution Returns the resolution of OLED.
func (oled SSD1331) Resolution() (int, int) {
	return width, height
}

// IsDisplay Whether the OLED is display.
func (oled *SSD1331) IsDisplay() bool {
	return oled.display
}

// DisplayOn Turn on the OLED.
func (oled *SSD1331) DisplayOn() {
	oled.sendCommand([]byte{0xAF})
	oled.display = true
}

// DisplayOnDim Turn on the OLED in dim mode.
func (oled *SSD1331) DisplayOnDim() {
	oled.sendCommand([]byte{0xAC})
	oled.display = true
}

// DisplayOff Turn off the OLED.
func (oled *SSD1331) DisplayOff() {
	oled.sendCommand([]byte{0xAE})
	oled.display = false
}

// SettingDimMode Configure dim mode setting
// r, g, b <= 255, preChargeVoltage <= 31
func (oled *SSD1331) SettingDimMode(r, g, b, preChargeVoltage byte) {
	oled.sendCommand(append([]byte{0xAB, 0}, r, g, b, preChargeVoltage))
}

// SetDisplayMode Change the display mode.
func (oled *SSD1331) SetDisplayMode(mode DisplayMode) {
	oled.sendCommand([]byte{byte(mode)})
}

// SetRGBContrast value <= 128. a little unstable...
func (oled *SSD1331) SetRGBContrast(r, g, b byte) {
	oled.sendCommand([]byte{0x81, r, 0x82, g, 0x83})

}

// Display Send buffer to the OLED.
func (oled *SSD1331) Display() {
	oled.sendDate(oled.buffer)
}

// Clear Clear the buffer.
func (oled *SSD1331) Clear() {
	oled.buffer = make([]byte, bufferLen)
}

// ClearDisplay Clear the display.
func (oled *SSD1331) ClearDisplay() {
	oled.Clear()
	oled.Display()
}

// DrawRect Draw a rectangle to the OLED.
func (oled *SSD1331) DrawRect(x0, y0, x1, y1, lineColorR, lineColorG, lineColorB, fillColorR, fillColorG, fillColorB int, fill bool) {
	var f byte = 0xA1
	if !fill {
		f = 0xA0
	}

	oled.sendCommand([]byte{
		0x26, f,
		0x22, byte(x0), byte(y0), byte(x1), byte(y1),
		byte(lineColorR), byte(lineColorG), byte(lineColorB),
		byte(fillColorR), byte(fillColorG), byte(fillColorB),
	})
}

func (oled *SSD1331) SetHLine(x, y, w, r, g, b int) {
	for i := 0; i < w; i++ {
		oled.SetPixel(x+i, y, r, g, b)
	}
}

func (oled *SSD1331) SetVLine(x, y, h, r, g, b int) {
	for i := 0; i < h; i++ {
		oled.SetPixel(x, y+i, r, g, b)
	}
}

// SetLine Write a line on the OLED.
func (oled *SSD1331) SetLine(x0, y0, x1, y1, r, g, b int) {
	// I referred to this(https://ja.wikipedia.org/wiki/ブレゼンハムのアルゴリズム/).
	dx := x1 - x0 // Abs
	dy := y1 - y0 // Abs
	err := dx - dy

	for {
		oled.SetPixel(x0, y0, r, g, b)
		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += -1
		}
		if e2 < dx {
			err += dx
			y0 += -1
		}
	}
}

// SetPixel Set the pixel in the buffer.
func (oled *SSD1331) SetPixel(x, y, r, g, b int) {
	startIDX := ((y * width) + x) * 2
	colorH := (r & 0b11111000) | (g >> 5)
	colorL := ((g << 3) & 0b11100000) | (b >> 3)

	oled.buffer[startIDX] = byte(colorH)
	oled.buffer[startIDX+1] = byte(colorL)
}

// ActiveScroll Scrool the display.
//
// horScrlOffset  Set number of column as horizontal scroll offset.
// startRowAddr   Define start row address.
// horRowScrl     Set number of rows to be horizontal scrolled.
// verScrlOffset  Set number of row as vertical scroll offset.
// scrlInterval   Set time interval between each scroll step.
func (oled *SSD1331) ActiveScroll(horScrlOffset, startRowAddr, horRowScrl, verScrlOffset byte, scrlInterval ScrollStep) {
	oled.sendCommand([]byte{
		0x27, horRowScrl, startRowAddr, horRowScrl, verScrlOffset, byte(scrlInterval),
		0x2F,
	})
}

// DeactiveScrool If Scrool function is Active then stop the scrool.
func (oled *SSD1331) DeactiveScrool() {
	oled.sendCommand([]byte{0x2E})
}

// LOCK IMC interface will no longer accept commands.
func (oled *SSD1331) LOCK() {
	oled.sendCommand([]byte{0xFD, 0x1B})
}

// UNLOCK Unlock the LOCK function.
func (oled *SSD1331) UNLOCK() {
	oled.sendCommand([]byte{0xFD, 0x0B})

}

func (oled *SSD1331) sendCommand(b []byte) {
	oled.CSPin.Out(low)
	oled.DCPin.Out(low)
	oled.connect.Tx(b, nil)
	oled.CSPin.Out(high)
}

func (oled *SSD1331) sendDate(b []byte) {
	oled.CSPin.Out(low)
	oled.DCPin.Out(high)
	oled.connect.Tx(b, nil)
	oled.CSPin.Out(high)
}
