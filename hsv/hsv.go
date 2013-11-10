package hsv

import (
	. "github.com/nightlifelover/GoMandelbrot/types"
	"math"
)

// TODO Rgb2hsv(..)
// Read http://en.literateprograms.org/RGB_to_HSV_color_space_conversion_%28C%29

// From http://stackoverflow.com/questions/8208905/hsv-0-255-to-rgb-0-255

func Hsv2rgb(in HSVColor) Color {
	h, s, v := in.H, in.S, in.V

	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var r, g, b float64

	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return Color{r, g, b}

}
