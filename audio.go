package main

import rl "github.com/gen2brain/raylib-go/raylib"
import "math"

func WaveFromFreq(sampleRate uint32, frequency float64) (rl.Wave) {
  var data = make([]byte, sampleRate * 4) // 4 is for sizeof(float32)

  for i := 0; i < len(data); i++ {
    s := math.Float32bits(float32(math.Sin((2.0 * math.Pi * frequency * float64(i)) / float64(sampleRate))))
    data[0] = byte(s >> 24)
    data[1] = byte(s >> 16)
    data[2] = byte(s >> 8)
    data[3] = byte(s)
  }

  return rl.NewWave(sampleRate, sampleRate, 32, 1, data)
}

