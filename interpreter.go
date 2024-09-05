package main

import "math/rand/v2"

const MEM_SIZE = 0x1000;
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
	V [16]uint8    // Registers
	I    uint16    // Address register

  PC   uint16    // Program counter
  SP   uint16    // Stack pointer

  dtimer uint8    // Delay timer register
  stimer uint8    // Sound timer register

	memory []uint8 // Memory
}

func Chip8_init() (chip8 Chip8) {
	chip8 = Chip8{}
	chip8.memory = make([]uint8, MEM_SIZE)

  chip8.PC = 0x200;
  chip8.SP = MEM_SIZE - 1;

  for i := 0; i < len(FONT_SPRITES); i++ {
    chip8.memory[FONTS_ADDR + uint16(i)] = FONT_SPRITES[i]
  }

	return
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

  chip8.V[0xF] = borrow;
}
func _8XY7(chip8 *Chip8, x,  y uint8) {
  var vx, vy = uint16(chip8.V[x]), uint16(chip8.V[y])
  var borrow uint8 = 1;
  if vy > vx { borrow = 0; }

  chip8.V[x] = uint8(vx - vy);
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
  var msb = uint8((chip8.PC >> 8) & 0xFF)
  var lsb = uint8(chip8.PC & 0xFF)

  chip8.memory[chip8.SP] = lsb; chip8.SP++;
  chip8.memory[chip8.SP] = msb; chip8.SP++;
}
func _00EE(chip8 *Chip8) {
  var lsb = uint16(chip8.memory[chip8.SP]); chip8.SP--;
  var msb = uint16(chip8.memory[chip8.SP]); chip8.SP--;

  chip8.PC = ((msb << 8) & 0xFF00) | (lsb & 0x00FF)
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
func _FX1E(chip8 *Chip8, x uint8) { chip8.I += uint16(x) }
func _DXYN(chip8 *Chip8, x, y, n uint8) {
  // TODO: Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I 
  // Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
}
func _FX29(chip8 *Chip8, x uint8) { chip8.I = 5 * uint16(chip8.V[x]) + FONTS_ADDR; }
func _FX33(chip8 *Chip8, x uint8) {
  var vx = chip8.V[x]
  chip8.memory[chip8.I] = vx / 100;       chip8.I++;
  chip8.memory[chip8.I] = (vx / 10) % 10; chip8.I++;
  chip8.memory[chip8.I] = vx % 10;        chip8.I++;
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
