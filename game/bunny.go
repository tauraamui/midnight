package game

import "github.com/faiface/pixel"

import "github.com/tauraamui/midnight/sprite"

type Bunny struct {
	spriteSheet        pixel.Picture
	rightMotionSprites []*pixel.Sprite
	leftMotionSprites  []*pixel.Sprite
}

func NewBunny(s pixel.Picture) *Bunny {
	bunny := Bunny{
		spriteSheet: s,
	}

	for i := 0; i < 4; i++ {
		bunny.rightMotionSprites = append(bunny.rightMotionSprites, sprite.Make(bunny.spriteSheet, i, 0))
		bunny.leftMotionSprites = append(bunny.leftMotionSprites, sprite.Make(bunny.spriteSheet, i, 1))
	}

	return &bunny
}
