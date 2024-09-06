package main

import rl "github.com/gen2brain/raylib-go/raylib"

import "image/color"

var BG_COLOR = rl.Black
var SPRITE_COLOR = rl.White

func main() {
  var chip8 = Chip8_create()
  Chip8_reset(&chip8)
  Chip8_load(&chip8, "3-corax+.ch8");

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

    // Copy Chip8's display buffer to a pixel buffer and then upload it to the texture
    for i := 0; i < WIDTH * HEIGHT; i++ {
      color := BG_COLOR
      if chip8.display[i] == 1 {
        color = SPRITE_COLOR
      }
      pxBuffer[i] = color;
    }
    rl.UpdateTexture(chip8Display.Texture, pxBuffer)

    // Math for keeping texture centered, fitted to window and with the same aspect ratio
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

    // Draw
		rl.BeginDrawing()
      rl.ClearBackground(rl.Green)
      rl.DrawTexturePro(chip8Display.Texture, srcRec, dstRec, origin, 0, rl.White)
		rl.EndDrawing()
	}
}
