package game

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tauraamui/midnight/sprite"
)

type Bunny struct {
	spriteSheet           pixel.Picture
	animSprites           []*pixel.Sprite
	rightMotionSprites    []*pixel.Sprite
	leftMotionSprites     []*pixel.Sprite
	upMotionSprites       []*pixel.Sprite
	downMotionSprites     []*pixel.Sprite
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

func (b *Bunny) Draw(
	win *pixelgl.Window,
	matrix pixel.Matrix,
	animSpeed float64,
	movingL, movingR, movingU, movingD bool,
) {

	if movingR {
		b.animSprites = b.rightMotionSprites
	}

	if movingL {
		b.animSprites = b.leftMotionSprites
	}

	if movingU {
		b.animSprites = b.upMotionSprites
	}

	if movingD {
		b.animSprites = b.downMotionSprites
	}

	if animSpeed > 0 && time.Since(b.sinceAnimFrameSwitch).Milliseconds() >= calcTimeTwixtSwitchMS(300, animSpeed) {
		b.currentAnimFrameIndex++
		if b.currentAnimFrameIndex >= len(b.animSprites) {
			b.currentAnimFrameIndex = 0
		}
		b.sinceAnimFrameSwitch = time.Now()
	}

	b.animSprites[b.currentAnimFrameIndex].Draw(win, matrix)
}

func (b *Bunny) loadSprites() {
	s, err := sprite.LoadSpritesheet("./assets/bunnysheet5.png")
	if err != nil {
		panic(err)
	}

	b.spriteSheet = s
	for i := 0; i < 7; i++ {
		b.upMotionSprites = append(b.upMotionSprites, sprite.Make(b.spriteSheet, i, 0, 48))
		b.downMotionSprites = append(b.downMotionSprites, sprite.Make(b.spriteSheet, i, 1, 48))
		b.rightMotionSprites = append(b.rightMotionSprites, sprite.Make(b.spriteSheet, i, 2, 48))
		b.leftMotionSprites = append(b.leftMotionSprites, sprite.Make(b.spriteSheet, i, 3, 48))
	}

	b.animSprites = b.downMotionSprites
}

func calcTimeTwixtSwitchMS(normal, speed float64) int64 {
	if speed > 0 {
		return int64(normal/speed) + 10
	}
	return int64(normal)
}
