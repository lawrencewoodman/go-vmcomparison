/*
 * Simple Generated Code Virtual Machine
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package codegen

import (
	"math"
)

// TODO: Make this configurable
const memSize = 32000

type CGVM struct {
	mem    [memSize]uint // Memory
	ac     uint          // 32-bit accumulator
	pc     uint          // Program Counter
	hltNow bool          // Whether to halt
	hltVal uint          // A value returned by HLT
}

func New() *CGVM {
	return &CGVM{}
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}

func calcBaseIndexAddr(v *CGVM, baseIndirect uint, indexIndirect uint) uint {
	// TODO: Assume always at least 4096 memory to avoid check
	base := v.mem[baseIndirect]
	index := v.mem[indexIndirect]
	return base + index
}

func calcIndirectAddr(v *CGVM, addr uint) uint {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	return v.mem[addr]
}

func op_HLT(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.hltNow = true
	v.hltVal = v.mem[addr]
	v.pc = mask32(v.pc + 1)
}

func op_ADD(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.ac = mask32(v.ac + v.mem[addr])
	v.pc = mask32(v.pc + 1)
}

func op_SUB(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.ac = mask32(v.ac - v.mem[addr])
	v.pc = mask32(v.pc + 1)
}

func op_AND(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.ac = v.ac & v.mem[addr]
	v.pc = mask32(v.pc + 1)
}

func op_STA(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.mem[addr] = v.ac
	v.pc = mask32(v.pc + 1)
	// TODO: what to do about PC being too high
}

func op_LDA(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.ac = v.mem[addr]
	v.pc = mask32(v.pc + 1)
}

func op_JMP(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	// TODO: swap calcAddr and range check for program variant
	v.pc = addr
}

func op_JEQ(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	if v.ac == 0 {
		v.pc = addr
	} else {
		v.pc = mask32(v.pc + 1)
	}
}

func op_JGT(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	if v.ac != 0 && v.ac <= math.MaxInt32 {
		v.pc = addr
	} else {
		v.pc = mask32(v.pc + 1)
	}
}

func op_DSZ(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.mem[addr] = mask32(v.mem[addr] - 1)
	if v.mem[addr] == 0 {
		v.pc = mask32(v.pc + 2)
	} else {
		v.pc = mask32(v.pc + 1)
		// TODO: this mask isn't quite right because of separate program
		// TODO: array
	}
}

func op_INC(v *CGVM, addr uint) {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	v.mem[addr] = mask32(v.mem[addr] + 1)
	v.pc = mask32(v.pc + 1)
}

func (v *CGVM) LoadMem(mem []uint) {
	copy(v.mem[:], mem)
}

func (v *CGVM) Run(program []func(*CGVM)) {
	for !v.hltNow {
		program[v.pc](v)
	}
}
