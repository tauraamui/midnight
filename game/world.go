package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

	fireflies []entity.Entity

	camPos                             pixel.Vec
	spriteSheet                        pixel.Picture
	worldCopyCanvas                    *pixelgl.Canvas
	maskCopyCanvas                     *pixelgl.Canvas
	ambientLightIntensity              *float32
	shaderCamPos                       *mgl32.Vec2
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

		fireflies:             []entity.Entity{},
		ambientLightIntensity: new(float32),
		shaderCamPos:          &mgl32.Vec2{},

		camPos: pixel.ZV,
		shader: shader.New("/assets/shader/nighttime.glsl"),
	}

	for i := 0; i < 3; i++ {
		pos := mgl32.Vec2{rand.Float32() * 500, rand.Float32() * 500}
		world.fireflies = append(world.fireflies, entity.NewFirefly(pos.X(), pos.Y()))
	}

	world.loadSprites()

	world.shader.StrReplaceCallbacks = append(world.shader.StrReplaceCallbacks, func(src string) string {
		return strings.Replace(
			src,
			"//FIREFLY_POSITION_UNIFORMS",
			fmt.Sprintf("uniform vec2[%d] fireflyPositions;", len(world.fireflies)),
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

func (w *World) Draw(
	win *pixelgl.Canvas,
	lights *pixelgl.Canvas,
	singleLight *pixelgl.Canvas,
	dbg *imdraw.IMDraw,
) {
	if w.maskCopyCanvas == nil || w.maskCopyCanvas.Bounds() != win.Bounds() {
		w.maskCopyCanvas = pixelgl.NewCanvas(win.Bounds())
	}

	if w.worldCopyCanvas == nil || w.worldCopyCanvas.Bounds() != win.Bounds() {
		w.worldCopyCanvas = pixelgl.NewCanvas(win.Bounds())
	}

	win.Clear(color.RGBA{R: 110, G: 201, B: 57, A: 255})
	// *w.shaderCamPos = mgl32.Vec2{float32(w.camPos.X) / float32(win.Bounds().Norm().W()/SCALE), float32(w.camPos.Y) / float32(win.Bounds().Norm().H()/SCALE)}
	w.Camera = pixel.IM.Scaled(w.camPos, SCALE).Moved(win.Bounds().Center().Sub(w.camPos))
	// fmt.Printf("CANVAS CAM POS: %v\n", w.camPos)
	// fmt.Printf("SHADER CAM POS: %v\n", w.shaderCamPos)

	win.SetMatrix(w.Camera)
	w.batch.Draw(win)

	win.SetMatrix(pixel.IM)
	w.Bunny.Draw(
		win,
		w.currentVelocity,
		w.movingL, w.movingR, w.movingU, w.movingD,
	)

	// clone world canvas unchanged
	// w.worldCopyCanvas.Clear(pixel.RGB(0, 0, 0))
	// win.Draw(w.worldCopyCanvas, pixel.IM.Moved(win.Bounds().Center()))

	w.maskCopyCanvas.Clear(pixel.RGB(0, 0, 0))
	w.maskCopyCanvas.SetColorMask(pixel.RGB(0, 0, 0).Mul(pixel.Alpha(0.41)))
	w.maskCopyCanvas.SetComposeMethod(pixel.ComposePlus)
	win.Draw(w.maskCopyCanvas, pixel.IM.Moved(win.Bounds().Center()))

	lights.Clear(pixel.Alpha(0))

	w.maskCopyCanvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	dbg.SetMatrix(w.Camera)
	dbg.Color = pixel.RGB(1, 0, 0)

	dbg.Push(pixel.V(float64(w.fireflies[0].Pos().X()), float64(w.fireflies[0].Pos().Y())))
	dbg.Push(pixel.V(float64(w.fireflies[1].Pos().X()), float64(w.fireflies[1].Pos().Y())))
	dbg.Line(1)

	dbg.Push(pixel.V(float64(w.fireflies[0].Pos().X()), float64(w.fireflies[0].Pos().Y())))
	dbg.Push(pixel.V(float64(w.fireflies[2].Pos().X()), float64(w.fireflies[2].Pos().Y())))
	dbg.Line(1)

	dbg.Push(pixel.V(float64(w.fireflies[1].Pos().X()), float64(w.fireflies[1].Pos().Y())))
	dbg.Push(pixel.V(float64(w.fireflies[2].Pos().X()), float64(w.fireflies[2].Pos().Y())))
	dbg.Line(1)
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
		*w.ambientLightIntensity = lightIntensity
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

	*w.ambientLightIntensity = lightIntensity
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
