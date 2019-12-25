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

func (js *JS) Update() bool {
	if js.win.JoystickPresent(*js.input) {
		for i := 0; i < len(js.axii); i++ {
			js.axii[i] = js.win.JoystickAxis(*js.input, i)
		}
		return true
	}
	js.axii = make([]float64, len(js.axii))
	return false
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

func (gp *Gamepad) Update() bool {
	if gp.joystick != nil {
		present := gp.joystick.Update()
		if present {
			return true
		}
		gp.joystick = nil
	}
	return false
}

func (gp *Gamepad) GetPressedButtons() []int {
	pressedButtons := []int{}
	jsAttached := gp.Update()
	if jsAttached {
		for b := 0; b < gp.joystick.win.JoystickButtonCount(*gp.joystick.input); b++ {
			if gp.joystick.win.JoystickPressed(*gp.joystick.input, b) {
				pressedButtons = append(pressedButtons, b)
			}
		}
	}
	return pressedButtons
}

func (gp *Gamepad) LeftJoystickPressed() bool {
	pressedButtons := gp.GetPressedButtons()
	for _, b := range pressedButtons {
		if b == 13 {
			return true
		}
	}
	return false
}

func (gp *Gamepad) RightJoystickPressed() bool {
	pressedButtons := gp.GetPressedButtons()
	for _, b := range pressedButtons {
		if b == 14 {
			return true
		}
	}
	return false
}

func (gp *Gamepad) MovingUp() (float64, bool) {
	jsAttached := gp.Update()
	if jsAttached {
		return gp.joystick.axii[1] * -1, gp.joystick.axii[1] <= -0.20
	}
	return 1, gp.win.Pressed(pixelgl.KeyW)
}

func (gp *Gamepad) MovingRight() (float64, bool) {
	jsAttached := gp.Update()
	if jsAttached {
		return gp.joystick.axii[0], gp.joystick.axii[0] >= 0.20
	}
	return 1, gp.win.Pressed(pixelgl.KeyD)
}

func (gp *Gamepad) MovingDown() (float64, bool) {
	jsAttached := gp.Update()
	if jsAttached {
		return gp.joystick.axii[1], gp.joystick.axii[1] >= 0.20
	}
	return 1, gp.win.Pressed(pixelgl.KeyS)
}

func (gp *Gamepad) MovingLeft() (float64, bool) {
	jsAttached := gp.Update()
	if jsAttached {
		return gp.joystick.axii[0] * -1, gp.joystick.axii[0] <= -0.20
	}
	return 1, gp.win.Pressed(pixelgl.KeyA)
}

func (gp *Gamepad) Debug() string {
	if gp.joystick != nil {
		return gp.joystick.Debug()
	}
	return "NO GAMEPAD ATTACHED"
}
