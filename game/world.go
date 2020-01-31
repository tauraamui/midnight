package game

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/markbates/pkger"
	"github.com/tauraamui/midnight/entity"
	"github.com/tauraamui/midnight/sprite"
	"github.com/tauraamui/midnight/ui/input"
	"github.com/tauraamui/midnight/ui/shader"
)

const (
	TILE_SIZE         = 32
	SCALE             = 2
	CAM_SPEED         = 160.0
	GRASS_TILES_START = 9
	GRASS_TILES_END   = 11
)

type World struct {
	Camera pixel.Matrix
	Bunny  *entity.Bunny
	Clock  *WorldClock

	fireflies             []entity.Entity
	ambientLightIntensity *float32
	shaderCamPos          *mgl32.Vec2

	camPos                             pixel.Vec
	spriteSheet                        pixel.Picture
	shader                             *shader.Shader
	batch                              *pixel.Batch
	grassTiles                         []*pixel.Sprite
	movingL, movingR, movingU, movingD bool
	currentVelocity                    float64
}

func NewWorld() *World {
	world := World{
		Camera: pixel.IM,
		Bunny:  entity.NewBunny(SCALE),
		Clock:  NewWorldClock(),

		fireflies: []entity.Entity{
			entity.NewFirefly(1, 1),
		},
		ambientLightIntensity: new(float32),

		camPos: pixel.ZV,
		shader: shader.New("/assets/shader/nighttime.glsl"),
	}
	world.loadSprites()

	world.shader.StrReplaceCallbacks = append(world.shader.StrReplaceCallbacks, func(src string) string {
		return strings.Replace(
			src,
			"//FIREFLY_POSITION_UNIFORMS",
			fmt.Sprintf("uniform vec2[%d] fireflyPositions", len(world.fireflies)),
			-1,
		)
	})

	world.shader.Uniforms["camPos"] = world.shaderCamPos
	world.shader.Uniforms["ambientLightIntensity"] = world.ambientLightIntensity

	for i, firefly := range world.fireflies {
		world.shader.Uniforms[fmt.Sprintf("fireflyPositions[%d]", i)] = firefly.Pos()
	}

	return &world
}

func (w *World) Update(gp *input.Gamepad, dt float64) *shader.Shader {
	for _, entity := range w.fireflies {
		entity.Update()
	}
	w.updateCamPos(gp, dt)
	// w.Clock.Update()
	w.updateShader()

	return w.shader
}

func (w *World) Draw(win *pixelgl.Canvas) {
	win.Clear(color.RGBA{R: 110, G: 201, B: 57, A: 255})
	w.shaderCamPos = &mgl32.Vec2{float32(w.camPos.X) / float32(win.Bounds().Norm().W()/SCALE), float32(w.camPos.Y) / float32(win.Bounds().Norm().H()/SCALE)}
	w.Camera = pixel.IM.Scaled(w.camPos, SCALE).Moved(win.Bounds().Center().Sub(w.camPos))
	win.SetMatrix(w.Camera)

	w.batch.Draw(win)

	w.Camera = pixel.IM
	win.SetMatrix(w.Camera)
	w.Bunny.Draw(
		win,
		w.currentVelocity,
		w.movingL, w.movingR, w.movingU, w.movingD,
	)
}

func (w *World) updateCamPos(gp *input.Gamepad, dt float64) {
	speedMultiplier := 1.0
	if gp.LeftJoystickPressed() {
		speedMultiplier = 2.5
	}
	w.currentVelocity = 0
	w.movingL, w.movingR, w.movingU, w.movingD = false, false, false, false

	if speed, movingL := gp.MovingLeft(); movingL {
		w.movingL = movingL
		w.currentVelocity = (CAM_SPEED * speed * speedMultiplier) * dt
		w.camPos.X -= w.currentVelocity
	}

	if speed, movingR := gp.MovingRight(); movingR {
		w.movingR = movingR
		w.currentVelocity = (CAM_SPEED * speed * speedMultiplier) * dt
		w.camPos.X += w.currentVelocity
	}

	if speed, movingU := gp.MovingUp(); movingU {
		w.movingU = movingU
		w.currentVelocity = (CAM_SPEED * speed * speedMultiplier) * dt
		w.camPos.Y += w.currentVelocity
	}

	if speed, movingD := gp.MovingDown(); movingD {
		w.movingD = movingD
		w.currentVelocity = (CAM_SPEED * speed * speedMultiplier) * dt
		w.camPos.Y -= w.currentVelocity
	}

}

func (w *World) updateShader() {
	var lightIntensity float32 = MINIMUM_LIGHT_INTENSITY
	if w.Clock.Current.Hour() >= 8 && w.Clock.Current.Hour() <= 17 {
		lightIntensity = 1
		return
	}

	currentHour := w.Clock.Current.Hour()
	currentMinutes := w.Clock.Current.Minute()

	isMorning := currentHour > 4 && currentHour < 8
	isEvening := currentHour > 17 && currentHour < 21

	if isMorning || isEvening {
		minuteOffset := func() int {
			if isMorning {
				return 300
			}
			return 1080
		}()

		minutesBetween := (currentHour*60 + currentMinutes - minuteOffset) + 1
		lightIntensity = (float32(.049999996) * float32(minutesBetween) * float32(.11111112))

		if isEvening {
			lightIntensity = 1 - lightIntensity
		}

		if lightIntensity < .11111112 {
			lightIntensity = .11111112
		}
	}

	w.ambientLightIntensity = &lightIntensity
}

func (w *World) loadSprites() {
	spritesheetFile, err := pkger.Open("/assets/img/terrain.png")
	if err != nil {
		panic(err)
	}

	s, err := sprite.LoadSpritesheet(spritesheetFile)
	if err != nil {
		panic(err)
	}

	w.spriteSheet = s
	w.grassTiles = append(w.grassTiles, sprite.Make(w.spriteSheet, 3, 0, 32))
	w.batch = pixel.NewBatch(&pixel.TrianglesData{}, w.spriteSheet)
	w.batch.Clear()
	for x := 0; x < 500; x++ {
		for y := 0; y < 500; y++ {
			grass := w.grassTiles[0]
			grass.Draw(w.batch, pixel.IM.Moved(pixel.V(float64(x*(32)), float64(y*(32)))))
		}
	}
}
