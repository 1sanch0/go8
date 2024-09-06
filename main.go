package main

import rl "github.com/gen2brain/raylib-go/raylib"

import (
  "image/color"
  "flag"
)

var LightBlue = rl.NewColor(93, 201, 226, 255)

var BG_COLOR = rl.Black
var SPRITE_COLOR = rl.White

func main() {
  // CLI
  var ticksPtr = flag.Uint("ticks", 24, "Number of ticks per frame")

  flag.Parse()

  var tail = flag.Args()
  var ticks = *ticksPtr
  var hasRom = len(tail) > 0

  // Main
  var chip8 = Chip8_create()

  Chip8_reset(&chip8)

  if hasRom {
    Chip8_load(&chip8, tail[0]);
  }

  rl.SetTraceLogLevel(rl.LogNone);
  rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(WIDTH * 10, HEIGHT * 10, "go8 - a simple chip8 interperter in go")
	defer rl.CloseWindow()
  rl.InitAudioDevice()
  defer rl.CloseAudioDevice()

  rl.SetTargetFPS(60)

  var pxBuffer = make([]color.RGBA, WIDTH*HEIGHT);
  var chip8Display = rl.LoadRenderTexture(WIDTH, HEIGHT)
  defer rl.UnloadRenderTexture(chip8Display);

  var wave = WaveFromFreq(48000, 400.0)
  var sound = rl.LoadSoundFromWave(wave)
//  defer rl.UnloadWave(wave) Dont unload bc WaveFromFreq uses Go's make data
  defer rl.UnloadSound(sound)
  // TODO: audio doesn't work

	var srcRec = rl.NewRectangle(0, 0, WIDTH, HEIGHT)
	var dstRec = rl.NewRectangle(0, 0, WIDTH, HEIGHT)
  var origin = rl.NewVector2(0, 0)

  var d0, d1, dt float64
  d0 = rl.GetTime()

	for !rl.WindowShouldClose() {
    // var dt = rl.GetFrameTime()
    var width, height = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())

    for i := uint(0); (i < ticks) && hasRom; i++ {
      d1 = rl.GetTime()
      dt = d1 - d0
      Chip8_tick(&chip8, dt)
      d0 = rl.GetTime()

      // Copy Chip8's display buffer to a pixel buffer and then upload it to the texture
      for i := 0; i < WIDTH * HEIGHT; i++ {
        color := BG_COLOR
        if chip8.display[i] == 1 {
          color = SPRITE_COLOR
        }
        pxBuffer[i] = color;
      }
      rl.UpdateTexture(chip8Display.Texture, pxBuffer)
    }

    // Drop in files
    if rl.IsFileDropped() {
      droppedFiles := rl.LoadDroppedFiles()
      defer rl.UnloadDroppedFiles()

      Chip8_reset(&chip8)
      Chip8_load(&chip8, droppedFiles[0])
      hasRom = true
    }

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
      rl.ClearBackground(BG_COLOR)

      if !hasRom {
        rl.DrawText("Drop in a ROM to start playing", int32(-origin.X+dstRec.Width/40), int32(-origin.Y+dstRec.Height/2), int32(dstRec.Height/8), LightBlue)
      }

      rl.DrawTexturePro(chip8Display.Texture, srcRec, dstRec, origin, 0, rl.White)
		rl.EndDrawing()
	}
}
