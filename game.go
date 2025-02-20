package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game type
type Game struct {
	GameOver bool
	Pause    bool

	FramesCounter int32
	TurnCounter   int

	BgColor   rl.Color
	GridColor rl.Color

	SuccessColor rl.Color
	InfoColor    rl.Color
	WarningColor rl.Color
	ErrorColor   rl.Color

	Message        string
	MessageTimeout int32
	MessageColor   rl.Color

	WindowShouldClose bool
	AvailableForInput bool
	MousePos          rl.Vector2

	Camera   rl.Camera2D
	ZoomMode int

	MaxVelocityDelta rl.Vector2
	MaxAcceleration  float32

	NumPlayers    int
	CurrentPlayer int
	Players       []*Player
	Balls         []Ball
	Collisions    []*Collision

	HighlightOn bool
	Highlight   Ball

	GridStep int
	GridSize int

	Track          Track
	InputAvailable bool
}

// NewGame - Start new game
func NewGame() (g Game) {
	g.FramesCounter = 0
	g.TurnCounter = 0
	g.WindowShouldClose = false
	g.InputAvailable = true

	g.GameOver = false
	g.Pause = false
	g.AvailableForInput = true

	// Theme
	g.BgColor = rl.White
	// GridColor rl.Color

	g.SuccessColor = rl.Green
	g.InfoColor = rl.Blue
	g.WarningColor = rl.Yellow
	g.ErrorColor = rl.Red

	g.Camera = rl.Camera2D{
		Offset:   rl.Vector2{X: 0, Y: 0},
		Target:   rl.Vector2{X: 0, Y: 0},
		Rotation: 0.0,
		Zoom:     1.0,
	}

	g.GridSize = 1000
	g.GridStep = 100

	g.MaxVelocityDelta = rl.Vector2{
		X: float32(g.GridStep),
		Y: float32(g.GridStep),
	}
	g.MaxAcceleration = rl.Vector2Length(g.MaxVelocityDelta)

	g.HighlightOn = true
	g.Highlight = Ball{
		Radius: 10,
		Color:  rl.Blue,
	}

	g.NumPlayers = 2
	for range g.NumPlayers {
		g.Players = append(g.Players, NewPlayer(
			// TODO Get from Track.Start
			rl.Vector2{
				X: float32(rand.Intn(10) * g.GridStep),
				Y: 0,
			},
		),
		)
	}
	g.CurrentPlayer = 0

	g.Track = Track{
		Start: [2]rl.Vector2{
			rl.Vector2{X: 0, Y: 0},
			rl.Vector2{X: 1000, Y: 0},
		},
		End: [2]rl.Vector2{
			rl.Vector2{X: 0, Y: 1000},
			rl.Vector2{X: 1000, Y: 1000},
		},
	}

	g.Collisions = make([]*Collision, 10)

	return
}

// Load - Load resources
func (g *Game) Load() {
}

// Unload - Unload resources
func (g *Game) Unload() {
}

func (g *Game) IsGettingInput() bool {
	g.GetInput()

	for i, p := range g.Players {
		if p.MoveRequested && p.NextMove == nil {
			g.CurrentPlayer = i
			// fmt.Println("Setting current player to ", i)
			return true
		}
	}

	// if all moves received -> Validate
	for i, p := range g.Players {
		if p.MoveRequested {
			if ValidateMove(p, *g) {
				g.Players[i].MoveRequested = false
				g.Players[i].NextMove.New = false
				g.Players[i].NextMove.Approved = true
				g.Players[i].IsCrashed = false
				g.Players[i].Status = ""
			} else {
				g.SetMessage(fmt.Sprintf("illegal move %s!", p.Name), 45, g.ErrorColor)
				g.Players[i].MoveRequested = true
				g.Players[i].NextMove = nil
				return true
			}
		}
	}

	g.SetMessage("All moves accepted", 30, g.InfoColor)
	return false
}

func (g *Game) GetInput() {
	// DEBUG Restore control during animation etc.
	if rl.IsKeyPressed(rl.KeyQ) {
		g.InputAvailable = !g.InputAvailable
	}
	if !g.InputAvailable {
		return
	}

	// Get keyboard input
	// Pause
	if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeyBack) {
		g.Pause = !g.Pause
	}
	// Highlight Next Move
	if rl.IsKeyPressed(rl.KeyH) {
		g.HighlightOn = !g.HighlightOn
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
		g.Camera.Zoom = rl.Clamp(g.Camera.Zoom*scaleFactor, 0.125, 64.0)
	}

	// Get Player Input
	// Get mouse input
	g.MousePos = rl.GetMousePosition()
	gp := g.getNearestGridPosition(
		rl.GetScreenToWorld2D(g.MousePos, g.Camera),
	)
	g.Highlight.Pos = gp

	// Prevent double press
	if rl.IsMouseButtonUp(rl.MouseButtonLeft) {
		g.AvailableForInput = true
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && g.AvailableForInput {
		g.Players[g.CurrentPlayer].NextMove = CalculateMove(
			g.Players[g.CurrentPlayer].Car.Position,
			gp,
			g.GridStep,
		)
		g.AvailableForInput = false
	}
}

// Update - Update game
func (g *Game) Update() {
	if rl.WindowShouldClose() {
		g.WindowShouldClose = true
		return
	}

	for i, _ := range g.Players {
		// Execute all the moves
		g.Players[i].ExecuteNextMove(*g)
		// Reset the user input request
		g.Players[i].MoveRequested = true
		g.Players[i].NextMove = nil
	}

	// Check Player-Player Collisions
	for j := 0; j < len(g.Players); j++ {
		for i := j; i < len(g.Players); i++ {
			if i != j {
				cp, ti, tj, err := CheckPlayerPlayerCollision(g.Players[j], g.Players[i])
				if err != nil {
					continue
				}
				if ((ti > 0) && (ti <= 1)) && ((tj > 0) && (tj <= 1)) {

					tImpact := max(ti, tj)
					fmt.Println(ti, tj, "->", tImpact)
					ftImpact := int32(tImpact * 60)
					explDuration := int32(60)
					if ti > tj {
						g.Players[i].IsCrashed = true
						g.Players[i].Status = "(on fire)"
						g.Players[i].MoveRequested = true
						g.Players[i].Car.PositionHistory[len(g.Players[i].Car.PositionHistory)-1] = g.getNearestGridPosition(cp)
						g.Players[i].Car.Velocity = rl.Vector2{0, 0}
					}
					if tj > ti {
						g.Players[j].IsCrashed = true
						g.Players[j].Status = "(on fire)"
						g.Players[j].MoveRequested = true
						g.Players[j].Car.PositionHistory[len(g.Players[j].Car.PositionHistory)-1] = g.getNearestGridPosition(cp)
						g.Players[j].Car.Velocity = rl.Vector2{0, 0}
					}
					if tj == ti {
						g.Players[i].IsCrashed = true
						g.Players[i].Status = "(on fire)"
						g.Players[i].MoveRequested = true
						g.Players[i].Car.PositionHistory[len(g.Players[i].Car.PositionHistory)-1] = g.getNearestGridPosition(cp)
						g.Players[i].Car.Velocity = rl.Vector2{0, 0}

						g.Players[j].IsCrashed = true
						g.Players[j].Status = "(on fire)"
						g.Players[j].MoveRequested = true
						g.Players[j].Car.PositionHistory[len(g.Players[j].Car.PositionHistory)-1] = g.getNearestGridPosition(cp)
						g.Players[j].Car.Velocity = rl.Vector2{0, 0}
					}

					g.Collisions = append(g.Collisions,
						NewCollision(
							cp, "Explosion",
							g.FramesCounter+ftImpact,
							explDuration,
						),
					)
					g.Collisions = append(g.Collisions,
						NewCollision(
							cp, "Burning",
							g.FramesCounter+ftImpact+explDuration,
							300,
						),
					)
				}
				if (ti <= 0) || (tj <= 0) {
					g.Collisions = append(g.Collisions, NewCollision(cp, "Past", g.FramesCounter, 180))
				}
				if (ti > 1) && (tj > 1) {
					g.Collisions = append(g.Collisions, NewCollision(cp, "Potential", g.FramesCounter, 120))
				}
				// fmt.Println(cp, ti, tj)
			}
		}
	}

	// Reset current Player
	g.CurrentPlayer = 0
	g.TurnCounter++
	fmt.Println("TURN=", g.TurnCounter)
	g.SetMessage(fmt.Sprintf("Turn %d!", g.TurnCounter), 60, g.InfoColor)
}

// Draw - Draw game
func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(g.BgColor)

	// Draw 2D Scene
	rl.BeginMode2D(g.Camera)

	// Draw Grid from (0,0)
	// TODO -> g.Grid.Draw()
	rl.PushMatrix()
	rl.Translatef(0, float32(g.GridSize), 0)
	rl.Rotatef(90, 1, 0, 0)
	rl.DrawGrid(int32(g.GridSize), float32(g.GridStep)) // Box is 50*GridSize (double in X foulded from center)
	rl.PopMatrix()

	// Draw Track
	g.Track.Draw()

	// Draw Players
	for _, p := range g.Players {
		p.Draw(g.FramesCounter)
	}

	// Draw Collisions
	for i := 0; i < len(g.Collisions); i++ {
		if g.Collisions[i] != nil {
			g.Collisions[i].Draw(g.FramesCounter)
		}
	}

	// Draw Highlights
	if g.HighlightOn {
		g.Highlight.Color = g.Players[g.CurrentPlayer].Color
		g.Highlight.Draw()
	}

	rl.EndMode2D()

	g.DrawPlayerStatusBox()
	g.DrawMessage()

	rl.EndDrawing()
	g.FramesCounter++
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

func (g *Game) DrawPlayerStatusBox() {
	if g.CurrentPlayer != -1 {
		rl.DrawText(
			fmt.Sprintf("Round %d: %s's Turn", g.TurnCounter, g.Players[g.CurrentPlayer].Name),
			screenWidth/2-200, 20,
			20, rl.ColorLerp(rl.Black, g.Players[g.CurrentPlayer].Color, 0.5))
	}

	// Draw Player List
	rl.DrawText(
		fmt.Sprintf(
			"Round %d",
			g.TurnCounter,
		),
		40, 40,
		20, rl.Black,
	)
	for i, _ := range g.Players {
		act := " "
		if g.CurrentPlayer == i {
			act = "*"
		}

		req := "-"
		if g.Players[i].MoveRequested {
			req = "*"
		}

		app := "-"
		mv := "(-,-)"

		if g.Players[i].NextMove != nil {
			if g.Players[i].NextMove.Approved {
				app = "+"
			}
			mv = fmt.Sprintf("(%d,%d)", g.Players[i].NextMove.DX, g.Players[i].NextMove.DY)
		}
		rl.DrawText(
			fmt.Sprintf("%s %d:%s req:%s app:%s mv:%s",
				act,
				i,
				g.Players[i].Name,
				req,
				app,
				mv,
			),
			40, int32(40+(i+1)*40),
			20, rl.ColorLerp(rl.Black, g.Players[i].Color, 0.5))
	}
}

func (g *Game) DrawMessage() {
	// fmt.Println(g.FramesCounter, g.MessageTimeout)
	if g.FramesCounter < g.MessageTimeout {
		rl.DrawText(
			g.Message, screenWidth/2-200, screenHeight/2, 40, g.MessageColor)
	}
}

func (g *Game) SetMessage(msg string, to int, col rl.Color) {
	g.Message = msg
	g.MessageColor = col
	g.MessageTimeout = g.FramesCounter + int32(to)
}
