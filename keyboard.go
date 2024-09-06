package main

import rl "github.com/gen2brain/raylib-go/raylib"

var keyboard = map[uint8]int32 {
  0x1: rl.KeyOne,
  0x2: rl.KeyTwo,
  0x3: rl.KeyThree,
  0xC: rl.KeyFour,

  0x4: rl.KeyQ,
  0x5: rl.KeyW,
  0x6: rl.KeyE,
  0xD: rl.KeyR,

  0x7: rl.KeyA,
  0x8: rl.KeyS,
  0x9: rl.KeyD,
  0xE: rl.KeyF,

  0xA: rl.KeyZ,
  0x0: rl.KeyX,
  0xB: rl.KeyC,
  0xF: rl.KeyV,
}

func IsKeyPressed(key uint8) bool {
  return rl.IsKeyDown(keyboard[key]);
}
