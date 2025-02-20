package main

import (
	"encoding/json"
	"io"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Track struct {
	Start    [2]rl.Vector2 `json:"start"`
	End      [2]rl.Vector2 `json:"end"`
	Vertices []rl.Vector2  `json:"vertices"`
}

func LoadTrack(fn string) *Track {
	f, _ := os.Open(fn)
	defer f.Close()

	bs, _ := io.ReadAll(f)
	var t Track
	json.Unmarshal(bs, &t)
	return &t
}

func (t Track) SaveTrack(fn string) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	f, _ := os.Open(fn)
	os.WriteFile(fn, b, 0777)
	defer f.Close()
	return nil
}

func (t Track) Draw() {
	// Draw Edge(s)
	for i := 1; i < len(t.Vertices); i++ {
		rl.DrawLineV(t.Vertices[i-1], t.Vertices[i], rl.Black)
	}
	// Start Start
	rl.DrawLineV(t.Start[0], t.Start[1], rl.Green)
	// Start End
	rl.DrawLineV(t.End[0], t.End[1], rl.Red)
}
