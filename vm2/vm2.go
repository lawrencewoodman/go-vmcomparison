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
	mem     [memSize]int64   // Memory
	pc      int64            // Program Counter
	hltVal  int64            // A value returned by HLT
	symbols map[string]int64 // The symbols table from the assembler - to aid debugging
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

func (v *VM2) Mem() [memSize]int64 {
	return v.mem
}

func (v *VM2) LoadRoutine(routine []int64, symbols map[string]int64) {
	copy(v.mem[:], routine)
	v.symbols = symbols
}

// fetch gets the next instruction from memory
// Returns: opcode, operandA, operandB
// TODO: describe instruction format
func (v *VM2) fetch() (int64, int64, int64, error) {
	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, v.pc+2)
	}
	opcode := v.mem[v.pc]
	operandA := v.mem[v.pc+1]
	operandB := v.mem[v.pc+2]

	// If addressing mode: operand A indirect
	if operandA < 0 {
		operandA = 0 - operandA
		if operandA >= memSize {
			return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
		}
		operandA = v.mem[operandA]
	}

	// If addressing mode: operand B indirect
	if operandB < 0 {
		operandB = 0 - operandB
		if operandB >= memSize {
			return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
		}
		operandB = v.mem[operandB]
	}

	if operandA >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
	}
	if operandB >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
	}

	// TODO: Decide if should increment PC here
	return opcode, operandA, operandB, nil
}

func (v *VM2) addr2symbol(addr int64) string {
	for k, v := range v.symbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}

func (v *VM2) opcode2mnemonic(opcode int64) string {
	for m, o := range instructions {
		if o == opcode {
			return m
		}
	}
	panic(fmt.Sprintf("opcode not found: %d", opcode))
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *VM2) execute(opcode int64, operandA int64, operandB int64) (bool, error) {
	//fmt.Printf("%7s:    %s   %s, %s\n", v.addr2symbol(v.pc), v.opcode2mnemonic(opcode), v.addr2symbol(operandA), v.addr2symbol(operandB))
	//fmt.Printf("            pre:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])

	switch opcode {
	case 0: // HLT
		v.hltVal = v.mem[operandA]
		// TODO: this wastes the following memory location, should it?
		return true, nil
	case 1: // MOV
		v.mem[operandB] = v.mem[operandA]
		v.pc += 3
	case 2: // JSR
		v.mem[operandB] = v.pc + 3
		v.pc = operandA
	case 3: // ADD
		v.mem[operandB] = v.mem[operandA] + v.mem[operandB]
		v.pc += 3
	case 4: // DJNZ
		v.mem[operandA] = v.mem[operandA] - 1
		if v.mem[operandA] != 0 {
			v.pc = operandB
		} else {
			v.pc += 3
		}
	case 5: // JMP
		v.pc = operandA + operandB
	case 6: // AND
		v.mem[operandB] = v.mem[operandA] & v.mem[operandB]
		v.pc += 3
	case 7: // OR
		v.mem[operandB] = v.mem[operandA] | v.mem[operandB]
		v.pc += 3
	case 8: // SHL
		v.mem[operandB] = v.mem[operandB] << v.mem[operandA]
		v.pc += 3
	case 9: // JNZ
		if v.mem[operandA] != 0 {
			v.pc = operandB
		} else {
			v.pc += 3
		}
	case 10: // SNE
		if v.mem[operandA] != v.mem[operandB] {
			v.pc += 6
		} else {
			v.pc += 3
		}
	case 11: // SLE
		if v.mem[operandA] <= v.mem[operandB] {
			v.pc += 6
		} else {
			v.pc += 3
		}
	case 12: // SUB
		v.mem[operandB] = v.mem[operandB] - v.mem[operandA]
		v.pc += 3
	case 13: // JGT
		if v.mem[operandA] > 0 {
			v.pc = operandB
		} else {
			v.pc += 3
		}

	default:
		return false, fmt.Errorf("PC: %d, unknown opcode: %d (%d)", v.pc, opcode, (opcode&0x3f000000)>>24)
	}

	//	fmt.Printf("            post:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])
	return false, nil
}
