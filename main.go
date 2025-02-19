package main

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1920 / 2
	screenHeight = 1080 / 2
)

func main() {
	game := NewGame()
	game.GameOver = true

	rl.InitWindow(screenWidth, screenHeight, "Gridpoint Clicker")

	rl.InitAudioDevice()
	game.Load()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}

	game.Unload()
	rl.CloseAudioDevice()
	rl.CloseWindow()
	os.Exit(0)
}
