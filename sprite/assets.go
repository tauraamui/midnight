package sprite

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/faiface/pixel"
)

func Make(sheet pixel.Picture, x, y, dim int) *pixel.Sprite {
	if int(sheet.Bounds().Max.X)%dim != 0 {
		spritesheetDimsIncorrectPanic(dim)
	}

	if int(sheet.Bounds().Max.Y)%dim != 0 {
		spritesheetDimsIncorrectPanic(dim)
	}

	realX := x * dim
	realY := y * dim

	if realX < int(sheet.Bounds().Min.X) || realX > int(sheet.Bounds().Max.X) {
		boundingErrorPanic()
	}

	if realY < int(sheet.Bounds().Min.Y) || realY > int(sheet.Bounds().Max.Y) {
		boundingErrorPanic()
	}

	return pixel.NewSprite(
		sheet, pixel.R(float64(realX), float64(realY), float64(realX+dim), float64(realY+dim)),
	)
}

func LoadSpritesheet(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func spritesheetDimsIncorrectPanic(dim int) {
	panic(fmt.Errorf("spritesheet bounds not divisible by %dpx", dim))
}

func boundingErrorPanic() { panic(errors.New("sprite x/y out of bounds")) }
