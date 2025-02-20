package main

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestCheckVector2Vector2Collision_01(t *testing.T) {
	v := rl.Vector2{1, 2}
	dv := rl.Vector2{0, 1}
	w := rl.Vector2{1, 2}
	dw := rl.Vector2{0, 1}

	_, t1, t2, err := CheckVector2Vector2Collision(v, dv, w, dw)
	if t1 != -1 {
		t.Errorf("t1 = %f; want -1", t1)
	}
	if t2 != -1 {
		t.Errorf("t2 = %f; want -1", t2)
	}
	if err != nil {
		t.Errorf("err = %s; want *nil", err.Error())
	}
}

func TestCheckVector2Vector2Collision_02(t *testing.T) {
	v := rl.Vector2{1, 1}
	dv := rl.Vector2{1, 1}
	w := rl.Vector2{2, 1}
	dw := rl.Vector2{-1, 1}
	_, t1, t2, err := CheckVector2Vector2Collision(v, dv, w, dw)
	if t1 != .5 {
		t.Errorf("t1 = %f; want .5", t1)
	}
	if t2 != .5 {
		t.Errorf("t2 = %f; want .5", t2)
	}
	if err != nil {
		t.Errorf("err = %s; want *nil", err.Error())
	}
}

func TestCheckVector2Vector2Collision_03(t *testing.T) {
	v := rl.Vector2{1, 1}
	dv := rl.Vector2{1, 2}
	w := rl.Vector2{2, 1}
	dw := rl.Vector2{-1, 1}
	_, t1, t2, err := CheckVector2Vector2Collision(v, dv, w, dw)
	if t1 != 1.0/3.0 {
		t.Errorf("t1 = %f; want .33", t1)
	}
	if t2 != 2.0/3.0 {
		t.Errorf("t2 = %f; want .67", t2)
	}
	if err != nil {
		t.Errorf("err = %s; want *nil", err.Error())
	}
}

func TestCheckVector2Vector2Collision_04(t *testing.T) {
	v := rl.Vector2{2, 1}
	dv := rl.Vector2{1, 1}
	w := rl.Vector2{1, 1}
	dw := rl.Vector2{-1, 1}
	_, t1, t2, err := CheckVector2Vector2Collision(v, dv, w, dw)
	if t1 != -.5 {
		t.Errorf("t1 = %f; want -.5", t1)
	}
	if t2 != -.5 {
		t.Errorf("t2 = %f; want -.5", t2)
	}
	if err != nil {
		t.Errorf("err = %s; want *nil", err.Error())
	}
}
