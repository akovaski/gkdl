package gkdl

type Event interface {
}

type mouseState struct {
	X int
	Y int
	State MouseButtons
}

type MouseMotion struct {
	mouseState
}

type MouseButtons struct {
	L bool // left
	R bool // right
	M bool // middle
	X1 bool // X1
	X2 bool // X2
}

type MouseButtonUp struct {
	Button MouseButtons // Which button was released
	mouseState
}

type MouseButtonDown struct {
	Button MouseButtons // Which button was pressed
	mouseState
}

type KeyDown struct {
	Scancode byte
	Keycode byte
}

type KeyUp struct {
	Scancode byte
	Keycode byte
}

type MouseWheel struct {
	WheelDelta int
	mouseState
}

type Quit struct {
}