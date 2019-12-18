package sprite

import (
	"errors"

	"github.com/faiface/pixel"
)

func Make(sheet pixel.Picture, x, y int) *pixel.Sprite {
	if int(sheet.Bounds().Max.X)%32 != 0 {
		spritesheetDimsIncorrectPanic()
	}

	if int(sheet.Bounds().Max.Y)%32 != 0 {
		spritesheetDimsIncorrectPanic()
	}

	realX := x * 32
	realY := y * 32

	if realX < int(sheet.Bounds().Min.X) || realX > int(sheet.Bounds().Max.X) {
		boundingErrorPanic()
	}

	if realY < int(sheet.Bounds().Min.Y) || realY > int(sheet.Bounds().Max.Y) {
		boundingErrorPanic()
	}

	return pixel.NewSprite(
		sheet, pixel.R(float64(realX), float64(realY), float64(realX)+32, float64(realY)+32),
	)
}

func spritesheetDimsIncorrectPanic() { panic(errors.New("spritesheet bounds not divisible by 32px")) }
func boundingErrorPanic()            { panic(errors.New("sprite x/y out of bounds")) }
