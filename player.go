package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Name          string
	Color         rl.Color
	Car           Car
	Moves         []*Move2
	NextMove      *Move2
	MoveRequested bool
}

func NewPlayer(sp rl.Vector2) Player {
	PlayerNames := []string{"Alice", "Bob", "Cedric", "Dave"}
	PlayerColors := []rl.Color{rl.Green, rl.Red, rl.Blue, rl.Yellow}
	clr := PlayerColors[rand.Perm(len(PlayerNames))[0]]
	return Player{
		Name:  PlayerNames[rand.Perm(len(PlayerNames))[0]],
		Color: clr,
		Car: Car{
			Model: Ball{
				Pos:    sp,
				Radius: 20,
				Color:  clr,
			},
			Velocity:        rl.Vector2{X: 0.0, Y: 100.0},
			Color:           clr,
			PositionHistory: []rl.Vector2{sp},
		},
		Moves:         nil,
		NextMove:      nil,
		MoveRequested: true,
	}

}

func (p *Player) Draw() {
	p.Car.DrawHistory()
	p.Car.Draw()
	p.Car.DrawVelocity()
	rl.DrawText(
		p.Name,
		int32(p.Car.Model.Pos.X),
		int32(p.Car.Model.Pos.Y),
		10, rl.ColorLerp(rl.Black, p.Color, 0.5))
}

// Pass in gameState
// Check vs. All Other Players
// Check vs. TrackEdges
// Check vs. TrackCheckpoints
func (p *Player) ExecuteNextMove(g Game) {
	p.Moves = append(p.Moves, p.NextMove)
	// Check for collisions

	// Update Car Position
	p.Car.Velocity = Move2ToPositionDelta(*p.NextMove, g.GridStep)
	p.Car.Model.Pos = rl.Vector2Add(
		p.Car.Model.Pos,
		p.Car.Velocity,
	)
	p.Car.PositionHistory = append(p.Car.PositionHistory, p.Car.Model.Pos)
}

// A Move is an integer VecN where the unit vectors are of size dv
type Move2 struct {
	DX       int
	DY       int
	New      bool
	Approved bool
}

func CalculateMove(sp rl.Vector2, ep rl.Vector2, res int) *Move2 {
	return &Move2{
		DX:       int(ep.X-sp.X) / res,
		DY:       int(ep.Y-sp.Y) / res,
		New:      true,
		Approved: false,
	}
}

func Move2ToPositionDelta(m Move2, sf int) rl.Vector2 {
	return rl.Vector2{
		X: float32(m.DX * sf),
		Y: float32(m.DY * sf),
	}
}

func ValidateMove(p Player, g Game) bool {
	acc := rl.Vector2Subtract(
		Move2ToPositionDelta(*p.NextMove, g.GridStep),
		p.Car.Velocity,
	)

	return rl.Vector2Length(acc) <= g.MaxAcceleration
}
