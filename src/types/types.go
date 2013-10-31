package types

const SCREEN_WIDTH int = 800
const SCREEN_HEIGHT int = 600

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}

type HSVColor struct {
	H float64
	S float64
	V float64
}

type Screen [SCREEN_WIDTH][SCREEN_HEIGHT]Color

type ScreenData struct {
	Pixels   *Screen
	ExecTime string
	IterNr   int
}

type State int

const (
	Play State = iota
	Stop
	StepFwd
	StepBack
)

// TODO Not a type, but still keep it here!

var CurrentState State
