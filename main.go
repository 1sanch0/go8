package main

//import rl "github.com/gen2brain/raylib-go/raylib"
import "fmt"

func main() {
  var chip8 = Chip8_init()

  _6XNN(&chip8, 0, 69)
  fmt.Println(chip8.V)

  // NOTES:
  // - The delay timer counts down perpetually @ 60Hz
  // - The sound timer counts down perpetually @ 60Hz
  //   - The minimum value the timer will respond is 02 (TODO)


	//rl.InitWindow(800, 450, "raylib [core] example - basic window")
	//defer rl.CloseWindow()

	//rl.SetTargetFPS(60)

	//for !rl.WindowShouldClose() {
	//	rl.BeginDrawing()

	//	rl.ClearBackground(rl.RayWhite)
	//	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

	//	rl.EndDrawing()
	//}
}
