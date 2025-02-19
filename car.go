package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Car struct {
	Model           Ball
	Velocity        rl.Vector2
	PositionHistory []rl.Vector2
	Color           rl.Color
}

func (c Car) Draw() {
	c.Model.Draw()
}

func (c Car) DrawHistory() {
	for i := 1; i < len(c.PositionHistory); i++ {
		b := float32(i) / float32(len(c.PositionHistory))
		col := rl.ColorAlpha(
			c.Color,
			b*b*b*b,
		)
		rl.DrawLineEx(c.PositionHistory[i-1], c.PositionHistory[i], 5*b, col)
	}
	c.Model.Draw()
}

func (c Car) DrawVelocity() {
	vc := rl.ColorLerp(rl.Black, c.Model.Color, 0.5)
	rl.DrawLineEx(
		c.Model.Pos,
		rl.Vector2Add(c.Model.Pos, c.Velocity),
		5.0,
		vc,
	)
}

type Ball struct {
	Pos    rl.Vector2
	Radius float32
	Color  rl.Color
}

func (b Ball) Draw() {
	rl.DrawCircleV(b.Pos, b.Radius, b.Color)
}
