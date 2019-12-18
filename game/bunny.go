package game

import "github.com/faiface/pixel"

type Bunny struct {
	spriteSheet       pixel.Picture
	leftMotionSprites []*pixel.Sprite
}

func NewBunny(s pixel.Picture, tt []pixel.Rect) *Bunny {
	bunny := Bunny{
		spriteSheet: s,
	}

	for i := 0; i < 3; i++ {
		bunny.leftMotionSprites = append(bunny.leftMotionSprites, pixel.NewSprite(bunny.spriteSheet, tt[i]))
	}

	return &bunny
}
