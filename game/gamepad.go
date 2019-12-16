package game

import "github.com/faiface/pixel/pixelgl"

type Gamepad struct {
	win   *pixelgl.Window
	input *pixelgl.Joystick
}

func NewGamepad(win *pixelgl.Window) *Gamepad {
	gp := Gamepad{win: win}
	for js := pixelgl.Joystick1; js <= pixelgl.JoystickLast; js++ {
		if win.JoystickPresent(js) {
			gp.input = &js
		}
	}
	return &gp
}

func (gp *Gamepad) MovingUp() bool {
	return gp.win.Pressed(pixelgl.KeyW)
}

func (gp *Gamepad) MovingRight() bool {
	return gp.win.Pressed(pixelgl.KeyD)
}

func (gp *Gamepad) MovingDown() bool {
	return gp.win.Pressed(pixelgl.KeyS)
}

func (gp *Gamepad) MovingLeft() bool {
	return gp.win.Pressed(pixelgl.KeyA)
}
