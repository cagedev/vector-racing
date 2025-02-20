package main

import (
	"errors"
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Name          string
	Color         rl.Color
	Car           *Car
	Moves         []*Move2
	NextMove      *Move2
	MoveRequested bool
	Status        string
	IsCrashed     bool
}

func NewPlayer(sp rl.Vector2) *Player {
	PlayerNames := []string{"Alice", "Bob", "Cedric", "Dave"}
	PlayerColors := []rl.Color{rl.Green, rl.Red, rl.Blue, rl.Yellow}
	clr := PlayerColors[rand.Perm(len(PlayerNames))[0]]
	return &Player{
		Name:  PlayerNames[rand.Perm(len(PlayerNames))[0]],
		Color: clr,
		Car: &Car{
			Model: Ball{
				Pos:    rl.Vector2Zero(),
				Radius: 20,
				Color:  clr,
			},
			Velocity:        rl.Vector2{X: 0.0, Y: 100.0},
			Color:           clr,
			Position:        sp,
			PositionHistory: []rl.Vector2{sp},
			Animation:       [2]int32{0, 1},
		},
		Moves:         nil,
		NextMove:      nil,
		MoveRequested: true,
	}
}

func (p *Player) Draw(t int32) {
	// Draw Tail
	p.Car.DrawHistory()

	// Draw Car
	p.Car.Draw(t)

	// Draw Velocity vector
	p.Car.DrawVelocity()

	// Draw Player label
	rl.PushMatrix()
	rl.Translatef(p.Car.Position.X, p.Car.Position.Y, 0)
	rl.DrawText(
		p.Name+p.Status,
		int32(p.Car.Model.Pos.X),
		int32(p.Car.Model.Pos.Y),
		80, rl.ColorLerp(rl.Black, p.Color, 0.5))
	rl.PopMatrix()
}

func (p *Player) ExecuteNextMove(g Game) {
	if p.NextMove == nil {
		return
	}
	// Check for collisions? -> in second loop? depends on strategy
	p.Moves = append(p.Moves, p.NextMove)

	// Update Car Position
	p.Car.Velocity = Move2ToPositionDelta(*p.NextMove, g.GridStep)

	// Don't directly update position, let Animation finish this
	p.Car.PositionHistory = append(
		p.Car.PositionHistory,
		rl.Vector2Add(
			p.Car.Position,
			p.Car.Velocity,
		),
	)
	p.Car.Animation[0] = g.FramesCounter + 0
	p.Car.Animation[1] = g.FramesCounter + 60
}

func (p *Player) Crash() {
	p.IsCrashed = true
	p.Status = "(on fire)"
	p.Car.Velocity = rl.Vector2{X: 0, Y: 0}
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

func ValidateMove(p *Player, g Game) bool {
	if p.IsCrashed {
		if (p.NextMove.DX != 0) || (p.NextMove.DY != 0) {
			return false
		}
		return true
	}
	acc := rl.Vector2Subtract(
		Move2ToPositionDelta(*p.NextMove, g.GridStep),
		p.Car.Velocity,
	)
	return rl.Vector2Length(acc) <= g.MaxAcceleration
}

func CheckPlayerPlayerCollision(p1 *Player, p2 *Player) (i rl.Vector2, t1 float32, t2 float32, err error) {

	// Only check for collisions if at least one move has been made
	if len(p2.Car.PositionHistory) < 2 || len(p1.Car.PositionHistory) < 2 {
		return rl.Vector2{}, 0, 0, errors.New("no vector available to collide")
	}

	v1 := p1.Car.PositionHistory[len(p1.Car.PositionHistory)-2]
	v2 := p1.Car.PositionHistory[len(p1.Car.PositionHistory)-1]
	dv := rl.Vector2Subtract(v2, v1)
	w1 := p2.Car.PositionHistory[len(p2.Car.PositionHistory)-2]
	w2 := p2.Car.PositionHistory[len(p2.Car.PositionHistory)-1]
	dw := rl.Vector2Subtract(w2, w1)

	fmt.Println("v:", v1, v2, dv)
	fmt.Println("w:", w1, w2, dw)

	return CheckVector2Vector2Collision(
		v1,
		dv,
		w1,
		dw,
	)
}

// func CheckPlayerTrackCollision(p Player, tr Track) (t float32) {}

// func CheckPlayerCheckpointCollision(p Player, cp CheckPoint) (t float32) {}

// A Move is an integer VecN where the unit vectors are of size dv
type Move2 struct {
	DX       int
	DY       int
	New      bool
	Approved bool
}
