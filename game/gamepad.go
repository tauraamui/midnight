package game

import (
	"strconv"

	"github.com/faiface/pixel/pixelgl"
)

type JS struct {
	win   *pixelgl.Window
	input *pixelgl.Joystick
	axii  []float64
}

func NewJS(win *pixelgl.Window, js *pixelgl.Joystick) *JS {
	joystick := JS{win: win, input: js}
	for i := 0; i < win.JoystickAxisCount(*js); i++ {
		joystick.axii = append(joystick.axii, win.JoystickAxis(*js, i))
	}
	return &joystick
}

func (js *JS) Update() {
	if js.win.JoystickPresent(*js.input) {
		for i := 0; i < len(js.axii); i++ {
			js.axii[i] = js.win.JoystickAxis(*js.input, i)
		}
	} else {
		js.axii = nil
	}
}

func (js *JS) Debug() string {
	debugString := "["
	for i := 0; i < len(js.axii); i++ {
		debugString += strconv.FormatFloat(js.axii[i], 'f', 3, 64)
		if i+1 < len(js.axii) {
			debugString += ", "
		}
	}
	debugString += "]"
	return debugString
}

type Gamepad struct {
	win      *pixelgl.Window
	joystick *JS
}

func NewGamepad(win *pixelgl.Window) *Gamepad {
	gp := Gamepad{win: win}
	gp.AttachToJoystick()
	return &gp
}

func (gp *Gamepad) JoystickAttached() bool { return gp.joystick != nil }

func (gp *Gamepad) AttachToJoystick() {
	for js := pixelgl.Joystick1; js < pixelgl.JoystickLast; js++ {
		if gp.win.JoystickPresent(js) {
			gp.joystick = NewJS(gp.win, &js)
			break
		}
	}
}

func (gp *Gamepad) Update() {
	if gp.joystick != nil {
		gp.joystick.Update()
	}
}

func (gp *Gamepad) MovingUp() bool {
	gp.Update()
	return gp.win.Pressed(pixelgl.KeyW)
}

func (gp *Gamepad) MovingRight() bool {
	gp.Update()
	return gp.win.Pressed(pixelgl.KeyD)
}

func (gp *Gamepad) MovingDown() bool {
	gp.Update()
	return gp.win.Pressed(pixelgl.KeyS)
}

func (gp *Gamepad) MovingLeft() bool {
	gp.Update()
	return gp.win.Pressed(pixelgl.KeyA)
}

func (gp *Gamepad) Debug() string {
	if gp.joystick != nil {
		return gp.joystick.Debug()
	}
	return "NO GAMEPAD ATTACHED"
}
