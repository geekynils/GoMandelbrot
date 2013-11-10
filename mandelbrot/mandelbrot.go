package mandelbrot

import (
	"math"
	"math/cmplx"
	// "os"
	"fmt"
	"hsv"
	// "log"
	"sync"
	"time"
	. "types"
)

const (
	min = -2 - 1i
	max = 1 + 2i
)

var screenOutputChan chan *ScreenData

func LinkScreenOutput(channel chan *ScreenData) {
	screenOutputChan = channel
}

var inputChan chan State

func LinkInput(channel chan State) {
	inputChan = channel
}

/* Calculates a color value from a complex number. Brightness is inverse to the
   distance from 0. Color value is dependent on the angle from the real axis,
   starting with red. */
func calculateColor(value complex128) Color {

	// See Wikipedia
	// http://en.wikipedia.org/wiki/Complex_number#Absolute_value_and_argument

	saturation := cmplx.Abs(value)

	// Assume that 2.2 is the max abs we can get, because sqrt(2*2 + 1) ~ 2.2

	saturation /= 2.2

	phase := cmplx.Phase(value) // Phase in [-Pi, Pi]

	if phase < 0 {
		phase = 2*math.Pi + phase
	}

	hue := phase / (2 * math.Pi)

	return hsv.Hsv2rgb(HSVColor{hue, saturation, 1})
}

func pointIteration(num complex128, maxIter int) complex128 {

	z := 0 + 0i

	for i := 0; i < maxIter; i++ {
		z = z*z + num
	}

	return z
}

func computePart(from, to int, screenData *Screen) {
	for i := from; i < to; i++ {
		for j := 0; j < height; j++ {

			real := (float64(i)/float64(width-1))*3 - 2
			imag := (float64(j)/float64(height-1))*2 - 1
			num := complex(real, imag)

			// z_0 = 0
			// z_(n+1) = z_n^2 + c

			value := pointIteration(num, it)

			colorValue := calculateColor(value)

			screenData[i][j] = colorValue
		}
	}
}

func DrawMandelbrot() {

	screenData := new(Screen)
	width := len(screenData)
	height := len(screenData[0])
	it := 1
	state := Play
	nThreads := 4
	var wg sync.WaitGroup

	for {

		select {
		case state = <-inputChan:
		default:
		}

		if state == StepBack && it <= 1 {
			state = Stop
		}

		if state == Play || state == StepFwd || state == StepBack {

			if state == StepBack {
				it--
			} else {
				it++
			}

			start := time.Now()

			partSize = width / nThreads

			for i = 0; i < width; i += partSize {
				wg.Add(1)
				go computePart(i, i+partSize+1, screenData)
			}

			wg.Wait()

			elapsed := time.Since(start)

			iterationData := new(ScreenData)

			iterationData.Pixels = screenData
			iterationData.IterNr = it
			iterationData.ExecTime = fmt.Sprintf("%s", elapsed)
			iterationData.NThreads = nThreads

			screenOutputChan <- iterationData

			if state == StepFwd || state == StepBack {
				state = Stop
			}

		}
	}

}
