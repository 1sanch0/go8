package main

import rl "github.com/gen2brain/raylib-go/raylib"

import "image/color"

func main() {
  var chip8 = Chip8_create()
  Chip8_reset(&chip8)
  Chip8_load(&chip8, "test_opcode.ch8");

  // NOTES:
  // - The delay timer counts down perpetually @ 60Hz
  // - The sound timer counts down perpetually @ 60Hz
  //   - The minimum value the timer will respond is 02 (TODO)

  rl.SetTraceLogLevel(rl.LogNone);
  rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(WIDTH * 10, HEIGHT * 10, "go8 - a simple chip8 interperter in go")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

  var pxBuffer = make([]color.RGBA, WIDTH*HEIGHT);
  var chip8Display = rl.LoadRenderTexture(WIDTH, HEIGHT)
  defer rl.UnloadRenderTexture(chip8Display);

	var srcRec = rl.NewRectangle(0, 0, WIDTH, HEIGHT)
	var dstRec = rl.NewRectangle(0, 0, WIDTH, HEIGHT)
  var origin = rl.NewVector2(0, 0)

	for !rl.WindowShouldClose() {
    var dt = rl.GetFrameTime()
    var width, height = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())

    Chip8_tick(&chip8, dt)

    for i := 0; i < WIDTH * HEIGHT; i++ {
      color := rl.Black
      if chip8.display[i] == 1 {
        color = rl.White
      }
      pxBuffer[i] = color;
    }
    rl.UpdateTexture(chip8Display.Texture, pxBuffer)

    // TODO: Copy chip8 display buffer to raylib texture
    // TODO REMOVE THIS LATER
    // rl.BeginTextureMode(chip8Display)
    //   rl.ClearBackground(rl.DarkGray)
    //   rl.DrawText("  go8", 10, 10, 3, rl.RayWhite)
    // rl.EndTextureMode()

    // Math for keeping texture centered and with the same aspect ratio
    origin = rl.NewVector2(0, 0)
    if (width / (WIDTH / HEIGHT)) > height {
      dstRec.Width = height * (WIDTH / HEIGHT)
      dstRec.Height = height

      origin.X = -(width - dstRec.Width) / 2
      if origin.X > 0 { origin.X = 0 }
    } else {
      dstRec.Width = width;
      dstRec.Height = width / (WIDTH / HEIGHT)

      origin.Y = -(height - dstRec.Height) / 2
      if origin.Y > 0 { origin.Y = 0 }
    }

		rl.BeginDrawing()
      rl.ClearBackground(rl.Green)
      rl.DrawTexturePro(chip8Display.Texture, srcRec, dstRec, origin, 0, rl.White)
		rl.EndDrawing()
	}
}
