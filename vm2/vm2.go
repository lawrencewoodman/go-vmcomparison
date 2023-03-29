/*
 * Simple Implementation Virtual Machine v2
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm2

import "fmt"

// TODO: video about refining instruction set benchmarks and testing

// TODO: Make this configurable
const memSize = 32000

type VM2 struct {
	mem    [memSize]uint // Memory
	pc     uint          // Program Counter
	hltVal uint          // A value returned by HLT
}

func New() *VM2 {
	return &VM2{}
}

func (v *VM2) Step() (bool, error) {
	opcode, operandA, operandB := v.fetch()
	return v.execute(opcode, operandA, operandB)
}

func (v *VM2) Run() (bool, error) {
	var err error
	hlt := false
	for !hlt {
		hlt, err = v.Step()
		if err != nil {
			return hlt, err
		}
	}
	return hlt, err
}

func (v *VM2) Mem() [memSize]uint {
	return v.mem
}

func (v *VM2) LoadRoutine(routine []uint) {
	copy(v.mem[:], routine)
}

// fetch gets the next instruction from memory
// Returns: opcode, operandA, operandB
// TODO: describe instruction format
func (v *VM2) fetch() (uint, uint, uint) {
	ir := v.mem[v.pc]
	// fmt.Printf("fetch PC: %d, ir: %d\n", v.pc, ir)
	opcode := (ir & 0xFF000000)
	operandA := (ir & 0xFFFFFF)
	operandB := v.mem[mask32(v.pc+1)]
	// TODO: Decide if should increment PC here
	return opcode, operandA, operandB
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *VM2) execute(opcode uint, operandA uint, operandB uint) (bool, error) {
	// fmt.Printf("PC: %d, opcode: %d (%d), A: %d, B: %d\n", v.pc, opcode, (opcode&0x3F000000)>>24, operandA, operandB)
	switch opcode {
	case 0 << 24: // HLT
		v.hltVal = v.mem[operandA]
		// TODO: this wastes the following memory location, should it?
		return true, nil
	case 1 << 24: // MOV
		v.mem[operandB] = v.mem[operandA]
		v.pc = mask32(v.pc + 2)
	case 2 << 24: // JSR
		v.mem[operandB] = mask32(v.pc + 2)
		v.pc = operandA
	case 3 << 24: // ADD
		v.mem[operandB] = mask32(v.mem[operandA] + v.mem[operandB])
		v.pc = mask32(v.pc + 2)
	case 4 << 24: // DJNZ
		v.mem[operandA] = mask32(v.mem[operandA] - 1)
		if v.mem[operandA] != 0 {
			v.pc = operandB
		} else {
			v.pc = mask32(v.pc + 2)
		}
	case 5 << 24: // JMP
		v.pc = operandA
	case 7 << 24: // AND
		v.mem[operandB] = v.mem[operandA] & v.mem[operandB]
		v.pc = mask32(v.pc + 2)
	case 8 << 24: // SHL
		v.mem[operandB] = mask32(v.mem[operandB] << v.mem[operandA])
		v.pc = mask32(v.pc + 2)
	case (3 | 0x80) << 24: // ADD I
		// TODO: Decide if I is first operand indirect and II is both, what if just 2nd?
		addr := v.mem[operandA]
		if addr >= memSize {
			// TODO: Implement an error
			panic("outside memory range")
		}
		v.mem[operandB] = mask32(v.mem[addr] + v.mem[operandB])
		v.pc = mask32(v.pc + 2)
	case (5 | 0x80) << 24: // JMP I
		v.pc = v.mem[operandA]
	case 6 << 24: // LIT
		//fmt.Printf("PC: %d  LIT  A: %d, B: %d\n", v.pc, operandA, operandB)
		v.mem[operandB] = operandA
		v.pc = mask32(v.pc + 2)
	default:
		panic(fmt.Sprintf("unknown opcode: %d (%d)", opcode, (opcode&0x3f000000)>>24))
	}
	return false, nil
}
