/*
 * Simple Generated Code Virtual Machine
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package codegen

import "errors"

// TODO: Make this configurable
const memSize = 32000

type CGVM struct {
	mem    [memSize]int64 // Memory
	ac     int64          // Accumulator
	pc     int64          // Program Counter
	hltNow bool           // Whether to halt
	hltVal int64          // A value returned by HLT
}

func New() *CGVM {
	return &CGVM{}
}

func calcBaseIndexAddr(v *CGVM, baseIndirect int64, indexIndirect int64) int64 {
	base := calcIndirectAddr(v, baseIndirect)
	index := calcIndirectAddr(v, indexIndirect)
	return base + index
}

func calcIndirectAddr(v *CGVM, addr int64) int64 {
	if addr >= memSize {
		// TODO: Implement an error
		panic("outside memory range")
	}
	return v.mem[addr]
}

func op_HLT(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.hltNow = true
	v.hltVal = v.mem[addr]
	v.pc++
	return nil
}

func op_ADD(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.ac += v.mem[addr]
	v.pc++
	return nil
}

func op_SUB(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.ac -= v.mem[addr]
	v.pc++
	return nil
}

func op_AND(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.ac &= v.mem[addr]
	v.pc++
	return nil
}

func op_STA(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.mem[addr] = v.ac
	v.pc++
	// TODO: what to do about PC being too high
	return nil
}

func op_LDA(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.ac = v.mem[addr]
	v.pc++
	return nil
}

func op_JMP(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	// TODO: swap calcAddr and range check for program variant
	v.pc = addr
	return nil
}

func op_JEQ(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	if v.ac == 0 {
		v.pc = addr
	} else {
		v.pc++
	}
	return nil
}

func op_JGT(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	if v.ac > 0 {
		v.pc = addr
	} else {
		v.pc++
	}
	return nil
}

func op_DSZ(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.mem[addr]--
	if v.mem[addr] == 0 {
		v.pc += 2
	} else {
		v.pc++
		// TODO: this mask isn't quite right because of separate program
		// TODO: array
	}
	return nil
}

func op_INC(v *CGVM, addr int64) error {
	if addr >= memSize {
		return errors.New("outside memory range")
	}
	v.mem[addr]++
	v.pc++
	return nil
}

func (v *CGVM) LoadMem(mem []int64) {
	copy(v.mem[:], mem)
}

func (v *CGVM) Run(program []func(*CGVM) error) error {
	for !v.hltNow {
		if err := program[v.pc](v); err != nil {
			return err
		}
	}
	return nil
}
