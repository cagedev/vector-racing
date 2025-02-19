package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CarVector struct {
	Start rl.Vector2
	End   rl.Vector2
	Color rl.Color
}

func (cv CarVector) Draw() {
	rl.DrawLine(
		int32(cv.Start.X),
		int32(cv.Start.Y),
		int32(cv.End.X),
		int32(cv.End.Y),
		cv.Color,
	)
}

type Car struct {
	Next     Ball
	Previous Ball
	Velocity rl.Vector2
	Vectors  []CarVector
}

func (c Car) Draw() {
	for _, v := range c.Vectors {
		v.Draw()
	}
	c.Next.Draw()
	c.Previous.Draw()
}

type Ball struct {
	X      int32
	Y      int32
	Radius int32
	Color  rl.Color
}

func (b Ball) Draw() {
	rl.DrawCircle(b.X, b.Y, float32(b.Radius), b.Color)
}

// Game type
type Game struct {
	// FxFlap  rl.Sound
	// FxSlap  rl.Sound
	// FxPoint rl.Sound
	// FxClick rl.Sound

	// TxSprites rl.Texture2D
	// TxSmoke   rl.Texture2D
	// TxClouds  rl.Texture2D

	// CloudRec rl.Rectangle
	// FrameRec rl.Rectangle

	GameOver bool
	Dead     bool
	Pause    bool
	SuperFX  bool

	Score         int
	HiScore       int
	FramesCounter int32

	WindowShouldClose bool

	MousePos rl.Vector2

	Camera   rl.Camera2D
	ZoomMode int

	MaxVelocityDelta rl.Vector2

	Balls []Ball

	HighlightOn bool
	Highlight   Car

	GridStep int
	GridSize int
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.Init()
	return
}

func (g *Game) Init() {
	// Gopher
	// g.Floppy = Floppy{rl.NewVector2(80, float32(screenHeight)/2-spriteSize/2)}

	// Sprite rectangle
	// g.FrameRec = rl.NewRectangle(0, 0, spriteSize, spriteSize)

	// Cloud rectangle
	// g.CloudRec = rl.NewRectangle(0, 0, float32(screenWidth), float32(g.TxClouds.Height))

	// Initialize particles
	// g.Particles = make([]Particle, maxParticles)
	// for i := 0; i < maxParticles; i++ {
	// 	g.Particles[i].Position = rl.NewVector2(0, 0)
	// 	g.Particles[i].Color = rl.RayWhite
	// 	g.Particles[i].Alpha = 1.0
	// 	g.Particles[i].Size = float32(rl.GetRandomValue(1, 30)) / 20.0
	// 	g.Particles[i].Rotation = float32(rl.GetRandomValue(0, 360))
	// 	g.Particles[i].Active = false
	// }

	// Pipes positions
	// g.PipesPos = make([]rl.Vector2, maxPipes)
	// for i := 0; i < maxPipes; i++ {
	// 	g.PipesPos[i].X = float32(480 + 360*i)
	// 	g.PipesPos[i].Y = -float32(rl.GetRandomValue(0, 240))
	// }

	// Pipes colors
	// colors := []rl.Color{
	// 	rl.Orange, rl.Red, rl.Gold, rl.Lime,
	// 	rl.Violet, rl.Brown, rl.LightGray, rl.Blue,
	// 	rl.Yellow, rl.Green, rl.Purple, rl.Beige,
	// }

	// Pipes
	// g.Pipes = make([]Pipe, maxPipes*2)
	// for i := 0; i < maxPipes*2; i += 2 {
	// 	g.Pipes[i].Rec.X = g.PipesPos[i/2].X
	// 	g.Pipes[i].Rec.Y = g.PipesPos[i/2].Y
	// 	g.Pipes[i].Rec.Width = pipesWidth
	// 	g.Pipes[i].Rec.Height = 550
	// 	g.Pipes[i].Color = colors[rl.GetRandomValue(0, int32(len(colors)-1))]

	// 	g.Pipes[i+1].Rec.X = g.PipesPos[i/2].X
	// 	g.Pipes[i+1].Rec.Y = 1200 + g.PipesPos[i/2].Y - 550
	// 	g.Pipes[i+1].Rec.Width = pipesWidth
	// 	g.Pipes[i+1].Rec.Height = 550

	// 	g.Pipes[i/2].Active = true
	// }

	g.Score = 0
	g.FramesCounter = 0
	g.WindowShouldClose = false

	g.GameOver = false
	g.Dead = false
	g.SuperFX = false
	g.Pause = false

	g.Camera = rl.Camera2D{
		Offset:   rl.Vector2{X: 0, Y: 0},
		Target:   rl.Vector2{X: 0, Y: 0},
		Rotation: 0.0,
		Zoom:     1.0,
	}

	// g.Balls = append(g.Balls,
	// 	Ball{
	// 		X:      0,
	// 		Y:      0,
	// 		Radius: 10.0,
	// 		Color:  rl.Red,
	// 	},
	// 	Ball{
	// 		X:      400,
	// 		Y:      300,
	// 		Radius: 50.0,
	// 		Color:  rl.Green,
	// 	},
	// )

	g.GridSize = 1000
	g.GridStep = 100

	g.MaxVelocityDelta = rl.Vector2{
		X: float32(g.GridStep),
		Y: float32(g.GridStep),
	}

	g.HighlightOn = true
	g.Highlight.Velocity = rl.Vector2{X: 0.0, Y: 0.0}
	g.Highlight.Previous = Ball{
		X:      0.0,
		Y:      0.0,
		Radius: 5.0,
		Color:  rl.Blue,
	}
	g.Highlight.Next = Ball{
		Radius: 10.0,
		Color:  rl.DarkBlue,
	}
}

// Load - Load resources
func (g *Game) Load() {
	// g.FxFlap = rl.LoadSound("assets/sounds/flap.wav")
	// g.FxSlap = rl.LoadSound("assets/sounds/slap.wav")
	// g.FxPoint = rl.LoadSound("assets/sounds/point.wav")
	// g.FxClick = rl.LoadSound("assets/sounds/click.wav")
	// g.TxSprites = rl.LoadTexture("assets/images/sprite.png")
	// g.TxSmoke = rl.LoadTexture("assets/images/smoke.png")
	// g.TxClouds = rl.LoadTexture("assets/images/clouds.png")
}

// Unload - Unload resources
func (g *Game) Unload() {
	// rl.UnloadSound(g.FxFlap)
	// rl.UnloadSound(g.FxSlap)
	// rl.UnloadSound(g.FxPoint)
	// rl.UnloadSound(g.FxClick)
	// rl.UnloadTexture(g.TxSprites)
	// rl.UnloadTexture(g.TxSmoke)
	// rl.UnloadTexture(g.TxClouds)
}

// Update - Update game
func (g *Game) Update() {
	if rl.WindowShouldClose() {
		g.WindowShouldClose = true
	}

	if !g.GameOver {
		if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeyBack) {
			g.Pause = !g.Pause
		}
	}

	// Get input
	g.MousePos = rl.GetMousePosition()

	// Highlight next move
	gp := g.getNearestGridPosition(
		rl.GetScreenToWorld2D(g.MousePos, g.Camera),
	)

	g.Highlight.Next.X = int32(gp.X)
	g.Highlight.Next.Y = int32(gp.Y)

	// calculate velocity delta
	dx := g.Highlight.Next.X - g.Highlight.Previous.X
	dy := g.Highlight.Next.Y - g.Highlight.Previous.Y
	ds := rl.Vector2{
		X: float32(dx),
		Y: float32(dy),
	}

	if rl.Vector2Length(rl.Vector2Subtract(ds, g.Highlight.Velocity)) > rl.Vector2Length(g.MaxVelocityDelta) {
		g.Highlight.Next.Color = rl.Red
	} else {
		g.Highlight.Next.Color = rl.DarkBlue
	}

	// Make move (ignore collisions...) -> Helper stuff
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		g.Highlight.Vectors = append(g.Highlight.Vectors,
			CarVector{
				Start: rl.Vector2{
					X: float32(g.Highlight.Previous.X),
					Y: float32(g.Highlight.Previous.Y),
				},
				End: rl.Vector2{
					X: float32(g.Highlight.Next.X),
					Y: float32(g.Highlight.Next.Y),
				},
				Color: rl.SkyBlue, // -> Fade based on stuff in draw loop
			},
		)
		g.Highlight.Velocity =
			rl.Vector2{
				X: float32(g.Highlight.Next.X - g.Highlight.Previous.X),
				Y: float32(g.Highlight.Next.Y - g.Highlight.Previous.Y),
			}
		g.Highlight.Previous = g.Highlight.Next
		fmt.Println("v=", rl.Vector2Length(g.Highlight.Velocity))
	}

	// Camera pan
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1.0/g.Camera.Zoom)
		g.Camera.Target = rl.Vector2Add(g.Camera.Target, delta)
	}

	// Camera zoom
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), g.Camera)

		g.Camera.Offset = rl.GetMousePosition()
		g.Camera.Target = mouseWorldPos

		scaleFactor := float32(1.0 + (0.25 * math.Abs(float64(wheel))))
		if wheel < 0 {
			scaleFactor = 1.0 / scaleFactor
		}
		g.Camera.Zoom = Clamp(g.Camera.Zoom*scaleFactor, 0.125, 64.0)
	}

	// Check Collisions
	// Check End
}

// Draw - Draw game
func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.Beige)

	// Draw 2D Scene
	rl.BeginMode2D(g.Camera)
	// Draw Grid from (0,0)
	rl.PushMatrix()
	rl.Translatef(0, float32(g.GridSize), 0)
	rl.Rotatef(90, 1, 0, 0)
	rl.DrawGrid(int32(g.GridSize), float32(g.GridStep)) // Box is 50*GridSize (double in X foulded from center)
	rl.PopMatrix()

	// for _, b := range g.Balls {
	// 	b.Draw()
	// }

	if g.HighlightOn {
		g.Highlight.Draw()
	}

	rl.EndMode2D()

	// Draw UI
	mp := g.MousePos
	wp := rl.GetScreenToWorld2D(g.MousePos, g.Camera)
	rl.DrawText(fmt.Sprintf("(% 5.1f, % 5.1f) => (% 5.1f, % 5.1f)", mp.Y, mp.Y, wp.X, wp.Y), 20, 20, 20, rl.Black)

	rl.EndDrawing()
}

// Clamp returns f clamped to [low, high]
func Clamp(f, low, high float32) float32 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func (g *Game) getNearestGridPosition(wp rl.Vector2) rl.Vector2 {
	xr := float32(0.5)
	if wp.X < 0 {
		xr = -xr
	}

	yr := float32(0.5)
	if wp.Y < 0 {
		yr = -yr
	}

	return rl.Vector2{
		X: float32(g.GridStep * int((wp.X)/float32(g.GridStep)+xr)),
		Y: float32(g.GridStep * int((wp.Y)/float32(g.GridStep)+yr)),
	}
}
