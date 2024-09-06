package main

import (
  "math/rand/v2"
  "os"
)

const WIDTH = 64
const HEIGHT = 32
const MEM_SIZE = 0x1000
const FONTS_ADDR uint16 = 0x0 // The memory address at which this data resides is unspecified but it's usually lower that 0x200
var FONT_SPRITES = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
	V [16]uint8     // Registers
	I    uint16     // Address register

  PC   uint16     // Program counter
  SP   uint16     // Stack pointer

  dtimer uint8    // Delay timer register
  stimer uint8    // Sound timer register

	memory []uint8  // Memory

  display []uint8 // Display
}

func Chip8_create() (chip8 Chip8) {
	chip8 = Chip8{}
	chip8.memory = make([]uint8, MEM_SIZE)
  chip8.display = make([]uint8, WIDTH * HEIGHT)

  for i := 0; i < len(FONT_SPRITES); i++ {
    chip8.memory[FONTS_ADDR + uint16(i)] = FONT_SPRITES[i]
  }

	return
}

func Chip8_reset(chip8 *Chip8) {
  for i := 0; i < len(chip8.V); i++ {
    chip8.V[i] = 0
  }
  chip8.I = 0

  chip8.PC = 0x200;
  chip8.SP = MEM_SIZE - 1;

  chip8.dtimer = 0
  chip8.stimer = 0

  for i := 0x200; i < len(chip8.memory); i++ {
    chip8.memory[i] = 0
  }

  for i := 0; i < len(chip8.display); i++ {
    chip8.display[i] = 0
  }
}

func Chip8_load(chip8 *Chip8, filename string) {
  data, err := os.ReadFile(filename)
  if err != nil {
    panic(err)
  }

  for i := 0; i < len(data); i++ {
    chip8.memory[0x200 + uint16(i)] = data[i]
  }
}

func Chip8_tick(chip8 *Chip8, dt float32) {
  // Update timers
  var increment = dt * 60.0
  chip8.dtimer -= uint8(increment)
  chip8.stimer -= uint8(increment)

  // Fetch
  var lsb = uint16(chip8.memory[chip8.PC]); chip8.PC++;
  var msb = uint16(chip8.memory[chip8.PC]); chip8.PC++;
  var opcode = int(((lsb << 8) & 0xFF00) | (msb & 0x00FF))

  var x = uint8((opcode >> 8) & 0xF)
  var y = uint8((opcode >> 4) & 0xF)
  var n = uint8(opcode & 0xF)
  var nn = uint8(opcode & 0xFF)
  var nnn = uint16(opcode & 0x0FFF)

  // Decode & Execute
  switch {
    case opcode == 0x00E0: _00E0(chip8)
    case opcode == 0x00EE: _00EE(chip8)
    case (0x1000 <= opcode) && (opcode < 0x2000): _1NNN(chip8, nnn)
    case (0x2000 <= opcode) && (opcode < 0x3000): _2NNN(chip8, nnn)
    case (0x3000 <= opcode) && (opcode < 0x4000): _3XNN(chip8, x, nn)
    case (0x4000 <= opcode) && (opcode < 0x5000): _4XNN(chip8, x, nn)
    case (0x5000 <= opcode) && (opcode < 0x6000): _5XY0(chip8, x, y)
    case (0x6000 <= opcode) && (opcode < 0x7000): _6XNN(chip8, x, nn)
    case (0x7000 <= opcode) && (opcode < 0x8000): _7XNN(chip8, x, nn)
    case (0x8000 <= opcode) && (opcode < 0x9000):
      switch n {
        case 0x0: _8XY0(chip8, x, y)
        case 0x1: _8XY1(chip8, x, y)
        case 0x2: _8XY2(chip8, x, y)
        case 0x3: _8XY3(chip8, x, y)
        case 0x4: _8XY4(chip8, x, y)
        case 0x5: _8XY5(chip8, x, y)
        case 0x6: _8XY6(chip8, x, y)
        case 0x7: _8XY7(chip8, x, y)
        case 0xE: _8XYE(chip8, x, y)
        default: panic("Invalid opcode")
      }
    case (0x9000 <= opcode) && (opcode < 0xA000): _9XY0(chip8, x, y)
    case (0xA000 <= opcode) && (opcode < 0xB000): _ANNN(chip8, nnn)
    case (0xB000 <= opcode) && (opcode < 0xC000): _BNNN(chip8, nnn)
    case (0xC000 <= opcode) && (opcode < 0xD000): _CXNN(chip8, x, nn)
    case (0xD000 <= opcode) && (opcode < 0xE000): _DXYN(chip8, x, y, n)
    case (0xE000 <= opcode) && (opcode < 0xF000):
      switch nn {
        case 0x9E: _EX9E(chip8, x)
        case 0xA1: _EXA1(chip8, x)
        default: panic("Invalid opcode")
      }
    case (0xF000 <= opcode):
      switch nn {
        case 0x07: _FX07(chip8, x);
        case 0x0A: _FX0A(chip8, x);
        case 0x15: _FX15(chip8, x);
        case 0x18: _FX18(chip8, x);
        case 0x1E: _FX1E(chip8, x);
        case 0x29: _FX29(chip8, x);
        case 0x33: _FX33(chip8, x);
        case 0x55: _FX55(chip8, x);
        case 0x65: _FX65(chip8, x);
        default: panic("Invalid opcode")
      }
    default: panic("Invalid opcode")
  }

  // Sound
  // TODO
}

// Instructions
// ---------------
func _6XNN(chip8 *Chip8, x, nn uint8) { chip8.V[x] = nn }
func _8XY0(chip8 *Chip8, x,  y uint8) { chip8.V[x] = chip8.V[y] }
func _7XNN(chip8 *Chip8, x, nn uint8) { chip8.V[x] += nn }
func _8XY4(chip8 *Chip8, x,  y uint8) {
  var vx, vy = uint16(chip8.V[x]), uint16(chip8.V[y])
  var carry uint8 = 0;
  if (vx + vy) > 0xFF { carry = 1 }

  chip8.V[x] = uint8(vx + vy);
  chip8.V[0xF] = carry;
}
func _8XY5(chip8 *Chip8, x,  y uint8) {
  var vx, vy = uint16(chip8.V[x]), uint16(chip8.V[y])
  var borrow uint8 = 1;
  if vy > vx { borrow = 0; }

  chip8.V[x] -= chip8.V[y];
  chip8.V[0xF] = borrow;
}
func _8XY7(chip8 *Chip8, x,  y uint8) {
  var vx, vy = uint16(chip8.V[x]), uint16(chip8.V[y])
  var borrow uint8 = 1;
  if vy < vx { borrow = 0; }

  chip8.V[x] = uint8(vy - vx);
  chip8.V[0xF] = borrow;
}
func _8XY2(chip8 *Chip8, x,  y uint8) { chip8.V[x] &= chip8.V[y] }
func _8XY1(chip8 *Chip8, x,  y uint8) { chip8.V[x] |= chip8.V[y] }
func _8XY3(chip8 *Chip8, x,  y uint8) { chip8.V[x] ^= chip8.V[y] }
func _8XY6(chip8 *Chip8, x,  y uint8) {
  var lsb uint8 = (chip8.V[y] & 0x80) >> 7
  chip8.V[x] = chip8.V[y] >> 1
  chip8.V[0xF] = lsb
}
func _8XYE(chip8 *Chip8, x,  y uint8) {
  var msb uint8 = chip8.V[y] & 0x01
  chip8.V[x] = chip8.V[y] << 1
  chip8.V[0xF] = msb
}
func _CXNN(chip8 *Chip8, x, nn uint8) { chip8.V[x] = uint8(rand.IntN(0x100)) & nn }
func _1NNN(chip8 *Chip8, nnn uint16) { chip8.PC = nnn; }
func _BNNN(chip8 *Chip8, nnn uint16) { chip8.PC = nnn + uint16(chip8.V[0]); }
func _2NNN(chip8 *Chip8, nnn uint16) {
  var lsb = uint8((chip8.PC >> 8) & 0xFF)
  var msb = uint8(chip8.PC & 0xFF)

  chip8.memory[chip8.SP] = lsb; chip8.SP--;
  chip8.memory[chip8.SP] = msb; chip8.SP--;

  chip8.PC = nnn;
}
func _00EE(chip8 *Chip8) {
  chip8.SP++;
  var msb = uint16(chip8.memory[chip8.SP]);
  chip8.SP++;
  var lsb = uint16(chip8.memory[chip8.SP]);

  chip8.PC = ((lsb << 8) & 0xFF00) | (msb & 0x00FF)
}
func _3XNN(chip8 *Chip8, x, nn uint8) { if chip8.V[x] == nn { chip8.PC += 2 } }
func _4XNN(chip8 *Chip8, x, nn uint8) { if chip8.V[x] != nn { chip8.PC += 2 } }
func _5XY0(chip8 *Chip8, x,  y uint8) { if chip8.V[x] == chip8.V[y] { chip8.PC += 2 } }
func _9XY0(chip8 *Chip8, x,  y uint8) { if chip8.V[x] != chip8.V[y] { chip8.PC += 2 } }
func _FX15(chip8 *Chip8, x uint8) { chip8.dtimer = chip8.V[x] }
func _FX07(chip8 *Chip8, x uint8) { chip8.V[x] = chip8.dtimer }
func _FX18(chip8 *Chip8, x uint8) { chip8.stimer = chip8.V[x] }
func _FX0A(chip8 *Chip8, x uint8) {
  var key uint8 = 0; // TODO: wait for a keypress
  chip8.V[x] = key;
}
func _EX9E(chip8 *Chip8, x uint8) {
  // TODO: Skip the following instruction if the key corresponding to the hex value currently stored in register VX is pressed
}
func _EXA1(chip8 *Chip8, x uint8) {
  // TODO: Skip the following instruction if the key corresponding to the hex value currently stored in register VX is not pressed
}
func _ANNN(chip8 *Chip8, nnn uint16) { chip8.I = nnn }
func _FX1E(chip8 *Chip8, x uint8) { chip8.I += uint16(chip8.V[x]) }
func _DXYN(chip8 *Chip8, x, y, n uint8) {
  var vx, vy = chip8.V[x], chip8.V[y]
  var xPos, yPos uint
  var data uint8

  chip8.V[0xF] = 0
  for i := 0; uint8(i) < n; i++ {
    data = chip8.memory[chip8.I + uint16(i)]
    yPos = (uint(vy) + uint(i)) % HEIGHT
    for j := 0; j < 8; j++ {
      xPos = (uint(vx) + uint(7-j)) % WIDTH
      idx := yPos * WIDTH + xPos      // Display buffer index

      Spx := (data >> j) & 0x1        // Sprite pixel
      Dpx := chip8.display[idx] & 0x1 // Display pixel

      px := Spx ^ Dpx                 // New pixel

      if (Dpx == 1 && Dpx == 0) { chip8.V[0xF] = 1 }

      chip8.display[idx] = px;
    }
  }
}
func _00E0(chip8 *Chip8) {
  for i := 0; i < WIDTH * HEIGHT; i++ {
    chip8.display[i] = 0;
  }
}
func _FX29(chip8 *Chip8, x uint8) { chip8.I = 5 * uint16(chip8.V[x]) + FONTS_ADDR; }
func _FX33(chip8 *Chip8, x uint8) {
  var vx = chip8.V[x]
  chip8.memory[chip8.I+0] = vx / 100;
  chip8.memory[chip8.I+1] = (vx / 10) % 10;
  chip8.memory[chip8.I+2] = vx % 10;
}
func _FX55(chip8 *Chip8, x uint8) {
  for i := uint8(0); i <= x; i++ {
    chip8.memory[chip8.I] = chip8.V[i]
    chip8.I++;
  }
}
func _FX65(chip8 *Chip8, x uint8) {
  for i := uint8(0); i <= x; i++ {
    chip8.V[i] = chip8.memory[chip8.I]
    chip8.I++;
  }
}
