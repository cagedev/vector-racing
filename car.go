package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Car struct {
	Model           Ball
	Position        rl.Vector2
	Velocity        rl.Vector2
	PositionHistory []rl.Vector2
	Color           rl.Color
	Animation       [2]int32
}

func (c *Car) Draw(t int32) {
	if len(c.PositionHistory) >= 2 {
		if t <= c.Animation[0] {
			c.Position = c.PositionHistory[len(c.PositionHistory)-2]
		} else if t > c.Animation[0] && t <= c.Animation[1] {
			f := float32(t-c.Animation[0]) / float32(c.Animation[1]-c.Animation[0])
			// fmt.Println(f)
			c.Position = rl.Vector2Lerp(
				c.PositionHistory[len(c.PositionHistory)-2],
				c.PositionHistory[len(c.PositionHistory)-1],
				f,
			)
		} else {
			c.Position = c.PositionHistory[len(c.PositionHistory)-1]
		}
	}
	rl.PushMatrix()
	rl.Translatef(c.Position.X, c.Position.Y, 0)
	c.Model.Draw()
	rl.PopMatrix()
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
	// c.Model.Draw()
}

func (c Car) DrawVelocity() {
	vc := rl.ColorLerp(rl.Black, c.Model.Color, 0.5)
	rl.DrawLineEx(
		c.Position,
		rl.Vector2Add(c.Position, c.Velocity),
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
