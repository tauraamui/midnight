package game

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/sprite"
)

type Bunny struct {
	spriteSheet           pixel.Picture
	rightMotionSprites    []*pixel.Sprite
	leftMotionSprites     []*pixel.Sprite
	currentAnimFrameIndex int
	sinceAnimFrameSwitch  time.Time
}

func NewBunny() *Bunny {
	bunny := Bunny{
		currentAnimFrameIndex: 0,
		sinceAnimFrameSwitch:  time.Now(),
	}
	bunny.loadSprites()

	return &bunny
}

func (b *Bunny) Draw(win *pixelgl.Window, matrix pixel.Matrix, animSpeed float64) {
	if animSpeed > 0 && time.Since(b.sinceAnimFrameSwitch).Milliseconds() >= calcTimeTwixtSwitchMS(100, animSpeed) {
		b.currentAnimFrameIndex++
		if b.currentAnimFrameIndex >= len(b.rightMotionSprites) {
			b.currentAnimFrameIndex = 0
		}
		b.sinceAnimFrameSwitch = time.Now()
	}
	b.rightMotionSprites[b.currentAnimFrameIndex].Draw(win, matrix)
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

func calcTimeTwixtSwitchMS(normal, speed float64) int64 {
	if speed > 0 {
		return int64(normal/speed) + 10
	}
	return int64(normal)
}
