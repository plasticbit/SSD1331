package OLED

import (
	"testing"
)

var (
	x0 int = 0
	y0 int = 0
	x1 int = 95
	y1 int = 63
)

func BenchmarkSet(b *testing.B) {
	// 471636              2574 ns/op
	// 481854              2594 ns/op
	// 181260              6544 ns/op
	// 181514              6586 ns/op
	// 461000              2634 ns/op
	for i := 0; i < b.N; i++ {
		dx := x1 - x0
		dy := y1 - y0

		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		if dy > dx {
			dx, dy = dy, dx
		}

		derr := float64(dy) / float64(dx)

		for x, y, err := 0, y0, .0; x < maxint(dx, dy); x++ {
			err += derr
			if err >= .5 {
				err--
				y++
			}
		}
	}
}

func BenchmarkSetf(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var sx, sy int
		dx := x1 - x0 // Abs
		dy := y1 - y0 // Abs

		if x0 < x1 {
			sx = 1
		} else {
			sx = -1
		}

		if y0 < y1 {
			sy = 1
		} else {
			sy = -1
		}

		err := dx - dy

		for {
			// setPixel(x, y)
			if x0 == x1 && y0 == y1 {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x0 += sx
			}
			if e2 < dx {
				err += dx
				y0 += sy
			}
		}
	}
}
