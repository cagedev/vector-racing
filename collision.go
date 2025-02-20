package main

import (
	"errors"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collision struct {
	Model    Ball
	Position rl.Vector2
	Color    rl.Color
}

func NewCollision(p rl.Vector2, ct string) Collision {
	CollisionColors := map[string]rl.Color{
		"Collision": rl.Red,
		"Potential": rl.Yellow,
		"Past":      rl.Gray,
	}
	col := CollisionColors[ct]
	return Collision{
		Model: Ball{
			Pos:    rl.Vector2Zero(),
			Radius: 10,
			Color:  col,
		},
		Position: p,
	}
}

func (c Collision) Draw() {
	rl.PushMatrix()
	rl.Translatef(c.Position.X, c.Position.Y, 0)
	c.Model.Draw()
	rl.PopMatrix()
}

func CheckVector2Vector2Collision(v, dv, w, dw rl.Vector2) (rl.Vector2, float32, float32, error) {
	// TODO Check Edge cases
	// - angle(dv) = angle(dw) -> parallel lines => Vector ^ point
	// -
	// - dv.X = 0 : vertical line -> exchange
	// - dv.Y = 0 : horizontal line -> works? => NaN,NaN

	if rl.Vector2Angle(dv, dw) == 0 {
		// TODO mod PI
		return rl.Vector2{}, 0, 0, errors.New("angle(dv) = angle(dw) -> parallel lines => Vector ^ point")
	}
	swapped := false
	var l, m, t float32

	if (dv.X == 0 && dv.Y == 0) || (dw.X == 0 && dw.Y == 0) {
		return rl.Vector2{}, 0, 0, errors.New("contains a zero velocity => Vector ^ point")
	}

	if dv.X == 0 && dw.Y == 0 {
		m = (v.X - w.X) / dw.X
		l = (w.Y - v.Y) / dv.Y
	} else {
		if dv.X == 0 || dv.Y == 0 {
			swapped = true
			t1 := v
			t2 := dv
			v = w
			dv = dw
			w = t1
			dw = t2
		}

		t = dv.X / dv.Y
		m = ((v.X - w.X) + t*(w.Y-v.Y)) / (dw.X - t*dw.Y)
		l = (w.X + m*dw.X - v.X) / dv.X
	}
	i := rl.Vector2{
		X: w.X + m*dw.X,
		Y: w.Y + m*dw.Y,
	}
	if swapped {
		return i, l, m, nil
	}
	return i, m, l, nil
}
