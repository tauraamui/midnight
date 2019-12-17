package game

import (
	"fmt"

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
	return fmt.Sprintf("AXII: %#v", js.axii)
}

type Gamepad struct {
	win      *pixelgl.Window
	joystick *JS
}

func NewGamepad(win *pixelgl.Window) *Gamepad {
	gp := Gamepad{win: win}
	for js := pixelgl.Joystick1; js < pixelgl.JoystickLast; js++ {
		if win.JoystickPresent(js) {
			gp.joystick = NewJS(win, &js)
			break
		}
	}
	return &gp
}

func (gp *Gamepad) MovingUp() bool {
	gp.joystick.Update()
	return gp.win.Pressed(pixelgl.KeyW)
}

func (gp *Gamepad) MovingRight() bool {
	gp.joystick.Update()
	return gp.win.Pressed(pixelgl.KeyD)
}

func (gp *Gamepad) MovingDown() bool {
	gp.joystick.Update()
	return gp.win.Pressed(pixelgl.KeyS)
}

func (gp *Gamepad) MovingLeft() bool {
	gp.joystick.Update()
	return gp.win.Pressed(pixelgl.KeyA)
}

func (gp *Gamepad) Debug() string {
	return gp.joystick.Debug()
}
