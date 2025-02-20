package main

import (
	"errors"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Collision struct {
	Model     *Ball
	Position  rl.Vector2
	Color     rl.Color
	Animation [2]int32
	Type      string
	Active    bool
}

func NewCollision(p rl.Vector2, ct string, t int32, dt int32) *Collision {
	CollisionColors := map[string]rl.Color{
		"Explosion": rl.Red,
		"Potential": rl.Yellow,
		"Past":      rl.Gray,
		"Burning":   rl.Red,
	}
	col := CollisionColors[ct]
	fmt.Println("New collisions at", p, "(", t, "->", t+dt, ")")
	return &Collision{
		Model: &Ball{
			Pos:    rl.Vector2Zero(),
			Radius: 10,
			Color:  col,
		},
		Position:  p,
		Animation: [2]int32{t, t + dt},
		Type:      ct,
		Active:    false,
	}
}

func (c *Collision) Draw(t int32) {
	if (t < c.Animation[0]) || (t > c.Animation[1]) {
		c.Active = false
	}
	if t > c.Animation[0] && t <= c.Animation[1] {
		c.Active = true
		f := float32(t-c.Animation[0]) / float32(c.Animation[1]-c.Animation[0])
		switch c.Type {
		case "Explosion":
			c.Model.Color = rl.Red
			er := InterpolateExplosionRadius(f)
			c.Model.Radius = er * 100
		case "Burning":
			c.Model.Radius = InterpolateFireRadius(f) * 50
			c.Model.Color = rl.ColorLerp(rl.Red, rl.Yellow, f*4)
		}
	}

	if c.Active {
		rl.PushMatrix()
		rl.Translatef(c.Position.X, c.Position.Y, 0)
		c.Model.Draw()
		rl.PopMatrix()
	}
}

func CheckVector2Point2Collision(v, dv, p rl.Vector2) (t float32, err error) {
	if (dv.X == 0) && (dv.Y == 0) {
		return 0, errors.New("=> point ^ point")
	}

	// test for collision
	if (p.X-v.X)*dv.Y != (p.Y-v.Y)*dv.X {
		return 0, errors.New("no collision")
	}
	if dv.X == 0 {
		t = (p.Y - v.Y) / dv.Y
	} else {
		t = (p.X - v.X) / dv.X
	}
	return
}

func CheckVector2Vector2Collision(v, dv, w, dw rl.Vector2) (rl.Vector2, float32, float32, error) {
	// TODO Check Edge cases
	// - angle(dv) = angle(dw) -> parallel lines => Vector ^ point
	// -
	// - dv.X = 0 : vertical line -> exchange
	// - dv.Y = 0 : horizontal line -> works? => NaN,NaN

	if rl.Vector2Angle(dv, dw) == 0 {
		// TODO mod PI
		fmt.Println("parallel lines")
		return rl.Vector2{}, 0, 0, errors.New("angle(dv) = angle(dw) -> parallel lines => Vector ^ point")
	}
	if rl.Vector2Angle(dv, dw) == math.Pi {
		// TODO mod PI
		fmt.Println("antiparallel lines")
		return rl.Vector2{}, 0, 0, errors.New("- angle(dv) = angle(dw) -> parallel lines => Vector ^ point")
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

func InterpolateExplosionAlpha(t float32) float32 {
	if t < 0 {
		return 0
	}
	if t < 0.25 {
		return 4 * t
	}
	if t < .5 {
		return 1
	}
	if t < 1 {
		return 1 - 2*(t-.5)
	}
	return 0
}

func InterpolateExplosionRadius(t float32) float32 {
	if t < 0 {
		return 0
	}
	if t < 1 {
		return t * t
	}
	return 0
}

func InterpolateFireRadius(t float32) float32 {
	if t < 0 {
		return 0
	}
	if t < 1 {
		return float32(
			(1-float64(t))*math.Sin(float64(2*math.Pi*t*3*10))*.25 + .75,
		)
	}
	return 0
}
