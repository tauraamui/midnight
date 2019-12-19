package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/sprite"
)

type Bunny struct {
	spriteSheet        pixel.Picture
	rightMotionSprites []*pixel.Sprite
	leftMotionSprites  []*pixel.Sprite
}

func NewBunny() *Bunny {
	bunny := Bunny{}
	bunny.loadSprites()

	return &bunny
}

func (b *Bunny) Draw(win *pixelgl.Window, matrix pixel.Matrix) {
	b.rightMotionSprites[0].Draw(win, matrix)
}

func (b *Bunny) loadSprites() {
	s, err := sprite.LoadSpritesheet("./assets/bunny.png")
	if err != nil {
		panic(err)
	}

	b.spriteSheet = s
	for i := 0; i < 4; i++ {
		b.rightMotionSprites = append(b.rightMotionSprites, sprite.Make(b.spriteSheet, i, 0, 48))
		b.leftMotionSprites = append(b.leftMotionSprites, sprite.Make(b.spriteSheet, i, 1, 48))
	}
}
