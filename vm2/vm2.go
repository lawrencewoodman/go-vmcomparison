/*
 * Simple Implementation Virtual Machine v2
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm2

import (
	"fmt"
)

// TODO: Make this configurable
const memSize = 32000

type VM2 struct {
	mem     [memSize]uint   // Memory
	pc      uint            // Program Counter
	hltVal  uint            // A value returned by HLT
	symbols map[string]uint // The symbols table from the assembler - to aid debugging
}

func New() *VM2 {
	return &VM2{}
}

func (v *VM2) Step() (bool, error) {
	opcode, operandA, operandB, err := v.fetch()
	if err != nil {
		return false, err
	}
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

func (v *VM2) LoadRoutine(routine []uint, symbols map[string]uint) {
	copy(v.mem[:], routine)
	v.symbols = symbols
}

// fetch gets the next instruction from memory
// Returns: opcode, operandA, operandB
// TODO: describe instruction format
func (v *VM2) fetch() (uint, uint, uint, error) {
	if v.pc+1 >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", v.pc+1)
	}
	ir := v.mem[v.pc]
	opcode := (ir & 0x3F000000)
	operandA := (ir & 0xFFFFFF)
	operandB := v.mem[v.pc+1]

	// If addressing mode: operand A indirect
	if ir&0x80000000 == 0x80000000 {
		if operandA >= memSize {
			return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandA)
		}
		operandA = v.mem[operandA]
	}

	// If addressing mode: operand B indirect
	if ir&0x40000000 == 0x40000000 {
		if operandB >= memSize {
			return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandB)
		}
		operandB = v.mem[operandB]
	}

	if operandA >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandA)
	}
	if operandB >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandB)
	}

	// TODO: Decide if should increment PC here
	return opcode, operandA, operandB, nil
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}

func (v *VM2) addr2symbol(addr uint) string {
	for k, v := range v.symbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}

func (v *VM2) opcode2mnemonic(opcode uint) string {
	for m, o := range instructions {
		if o == opcode&0x3F000000 {
			return m
		}
	}
	panic("opcode not found")
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *VM2) execute(opcode uint, operandA uint, operandB uint) (bool, error) {
	//	fmt.Printf("%7s:    %s   %s, %s\n", v.addr2symbol(v.pc), v.opcode2mnemonic(opcode), v.addr2symbol(operandA), v.addr2symbol(operandB))
	//	fmt.Printf("            pre:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])

	switch opcode {
	case 0 << 24: // HLT
		v.hltVal = v.mem[operandA]
		// TODO: this wastes the following memory location, should it?
		return true, nil
	case 1 << 24: // MOV
		v.mem[operandB] = v.mem[operandA]
		v.pc += 2
	case 2 << 24: // JSR
		v.mem[operandB] = mask32(v.pc + 2)
		v.pc = operandA
	case 3 << 24: // ADD
		v.mem[operandB] = mask32(v.mem[operandA] + v.mem[operandB])
		v.pc += 2
	case 4 << 24: // DJNZ
		v.mem[operandA] = mask32(v.mem[operandA] - 1)
		if v.mem[operandA] != 0 {
			v.pc = operandB
		} else {
			v.pc += 2
		}
	case 5 << 24: // JMP
		v.pc = operandA + operandB
	case 7 << 24: // AND
		v.mem[operandB] = v.mem[operandA] & v.mem[operandB]
		v.pc += 2
	case 8 << 24: // OR
		v.mem[operandB] = v.mem[operandA] | v.mem[operandB]
		v.pc += 2
	case 9 << 24: // SHL
		v.mem[operandB] = mask32(v.mem[operandB] << v.mem[operandA])
		v.pc += 2
	case 10 << 24: // JNZ
		if v.mem[operandA] != 0 {
			v.pc = operandB
		} else {
			v.pc += 2
		}
	case 11 << 24: // SNE
		if v.mem[operandA] != v.mem[operandB] {
			v.pc += 4
		} else {
			v.pc += 2
		}
	case 12 << 24: // SLE
		if v.mem[operandA] <= v.mem[operandB] {
			v.pc += 4
		} else {
			v.pc += 2
		}
	case 13 << 24: // SUB
		v.mem[operandB] = mask32(v.mem[operandB] - v.mem[operandA])
		v.pc += 2

	default:
		return false, fmt.Errorf("unknown opcode: %d (%d)", opcode, (opcode&0x3f000000)>>24)
	}

	//	fmt.Printf("            post:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])
	return false, nil
}
