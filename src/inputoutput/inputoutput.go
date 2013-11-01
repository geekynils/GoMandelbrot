package inputoutput

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"github.com/go-gl/gltext"
	"log"
	"os"
	"runtime"
	"time"
	. "types"
)

var ScreenOutputChan = make(chan *ScreenData)

var InputChan = make(chan State, 100)

var running = false

var state State = Play

func Init() {

	runtime.LockOSThread()

	// Initialize GLFW

	var err error

	if err = glfw.Init(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	err = glfw.OpenWindow(SCREEN_WIDTH, SCREEN_HEIGHT,
		0, 0, 0, 0, 0, 0, glfw.Windowed)

	if err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	glfw.SetWindowTitle("Mandelbrot")
	glfw.SetSwapInterval(1)

	glfw.SetWindowSizeCallback(onResize)
	glfw.SetWindowCloseCallback(onClose)
	glfw.SetMouseButtonCallback(onMouseBtn)
	glfw.SetMouseWheelCallback(onMouseWheel)
	glfw.SetKeyCallback(onKey)
	glfw.SetCharCallback(onChar)

	// Initialize OpenGL
	gl.Disable(gl.DEPTH_TEST)
	gl.ClearColor(0, 0, 0, 0)
}

// loadFont loads the specified font at the given scale.
func loadFont(file string, scale int32) (*gltext.Font, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	return gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
}

func drawInfoBox(x, y float64, fontSize, iteration, nThreads int, execTime string) {

	// TODO
	font, err := loadFont("/Users/nils/tmp/DejaVuSansMono.ttf", int32(fontSize))

	if err != nil {
		log.Printf("LoadFont: %v", err)
		return
	}

	firstLine := fmt.Sprintf("Iteration nr %d", iteration)

	secondLine := fmt.Sprintf("Took %s using %d threads.", execTime, nThreads)

	w1, h1 := font.Metrics(firstLine)

	w2, h2 := font.Metrics(secondLine)

	var w int

	if w1 > w2 {
		w = w1
	} else {
		w = w2
	}

	w += 3

	h := h1 + h2

	gl.Color4f(1, 1, 1, 0.7)
	gl.Rectd(x, y, x+float64(w), y+float64(h))
	gl.Color4f(0, 0, 0, 1)

	err = font.Printf(float32(x)+3, float32(y), firstLine)

	if err != nil {
		log.Printf("Something went wrong when drawing fonts: %v", err)
		return
	}

	err = font.Printf(float32(x)+3, float32(y+float64(h1)), secondLine)

	if err != nil {
		log.Printf("Something went wrong when drawing fonts: %v", err)
		return
	}

}

func Run() {

	running = true

	for running {

		start := time.Now()

		escPressed := glfw.Key(glfw.KeyEsc)
		windowOpen := glfw.WindowParam(glfw.Opened)

		running = escPressed == 0 && windowOpen == 1

		select {
		case data := <-ScreenOutputChan:

			drawFrame(data.Pixels)
			// TODO nThreads
			drawInfoBox(3, 3, 12, data.IterNr, 1, data.ExecTime)
			glfw.SwapBuffers()

		default: // Non blocking!
			glfw.PollEvents()

		}

		// Limit to 60hz

		elapsed := time.Since(start)

		timeToSleep := 16*time.Millisecond - elapsed

		if timeToSleep > 0 {
			time.Sleep(timeToSleep)
		}

	}
}

func drawFrame(screenData *Screen) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Begin(gl.POINTS)

	for i := 0; i < SCREEN_WIDTH; i++ {
		for j := 0; j < SCREEN_HEIGHT; j++ {
			var pixel Color = screenData[i][j]
			gl.Color3d(pixel.Red, pixel.Green, pixel.Blue)
			gl.Vertex2i(i, j)
		}
	}
	gl.End()
}

// --- Callbacks ---------------------------------------------------------------

// TODO Dynamically fetch size and render accordingly.
func onResize(w, h int) {
	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
	gl.ClearColor(0.255, 0.255, 0.255, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	log.Printf("resized: %dx%d\n", w, h)
}

func onClose() int {
	glfw.CloseWindow()
	glfw.Terminate()
	log.Println("closed")
	return 1 // return 0 to keep window open.
}

func onMouseBtn(button, state int) {
	log.Printf("mouse button: %d, %d\n", button, state)
}

func onMouseWheel(delta int) {
	log.Printf("mouse wheel: %d\n", delta)
}

func onKey(key, state int) {
	log.Printf("key: %d, %d\n", key, state)
}

func onChar(key, keyState int) {
	log.Printf("char: %d, %d\n", key, keyState)

	if keyState == glfw.KeyPress {

		switch key {
		case glfw.KeySpace: // space
			if state == Stop {
				state = Play
			} else {
				state = Stop
			}
		case 102: // f
			state = StepFwd
		case 98: // b
			state = StepBack
		default:
			return
		}

		log.Printf("Sending state %d", state)

		InputChan <- state
	}
}

// -----------------------------------------------------------------------------
