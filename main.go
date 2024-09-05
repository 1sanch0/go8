package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
  var chip8 = Chip8_create()
  Chip8_reset(&chip8)

  // NOTES:
  // - The delay timer counts down perpetually @ 60Hz
  // - The sound timer counts down perpetually @ 60Hz
  //   - The minimum value the timer will respond is 02 (TODO)

	rl.InitWindow(800, 450, "go8 - a simple chip8 interperter in go")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
    dt := rl.GetFrameTime()

    Chip8_tick(&chip8, dt)

		rl.BeginDrawing()
      rl.ClearBackground(rl.RayWhite)
		  rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}
}
