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

	rl.InitWindow(screenWidth, screenHeight, "Vector Racing")
	// rl.InitAudioDevice()

	// game.Load()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// 1. Gather Player Moves (only gather required player's move)
		// 2. Do Game Update(s)
		// Allows Multiple Ordering Strategies:
		//  - Sequential
		//  - Simultaneous (Turn Based) Resolving
		for game.IsGettingInput() {
			// Allow closing game
			if rl.WindowShouldClose() {
				break
			}
			// Draw game regardless
			game.Draw()
		}

		game.Update()
	}

	// game.Unload()
	// rl.CloseAudioDevice()
	rl.CloseWindow()

	os.Exit(0)
}
