/*
 * Virtual machine with 1 operand
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm1

import (
	"fmt"
)

// TODO: Make this configurable
const memSize = 32000

type VM1 struct {
	mem    [memSize]int64 // Memory
	pc     int64          // Program Counter
	ac     int64          // 32-bit accumulator
	x      int64          // 32-bit index? register
	y      int64          // 32-bit index? register
	r      int64          // 32-bit return register
	hltVal int64          // A value returned by HLT
}

func New() *VM1 {
	return &VM1{}
}

func (s *VM1) Step() (bool, error) {
	opcode, addr, err := s.fetch()
	if err != nil {
		return false, err
	}
	return s.execute(opcode, addr)
}

func (s *VM1) Run() (bool, error) {
	var err error
	hlt := false
	for !hlt {
		hlt, err = s.Step()
		if err != nil {
			return hlt, err
		}
	}
	return hlt, err
}

func (s *VM1) Mem() [memSize]int64 {
	return s.mem
}

func (v *VM1) LoadRoutine(routine []int64) {
	copy(v.mem[:], routine)
}

// fetch gets the next instruction from memory
// Returns: opcode, addr
// TODO: describe instruction format
func (s *VM1) fetch() (int64, int64, error) {
	if s.pc+1 >= memSize {
		return 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", s.pc, s.pc+1)
	}
	opcode := s.mem[s.pc]
	operand := s.mem[s.pc+1]

	// if addressing mode: indirect
	if operand < 0 {
		operand = -operand
		if operand >= memSize {
			return 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", s.pc, operand)
		}
		operand = s.mem[operand]
	}
	if operand >= memSize {
		return 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", s.pc, operand)
	}

	return opcode, operand, nil
}

func (v *VM1) opcode2mnemonic(opcode int64) string {
	for m, o := range instructions {
		if o == opcode {
			return m
		}
	}
	panic("opcode not found")
}

// execute executes the supplied instruction
// Returns: hlt, error
func (s *VM1) execute(opcode, addr int64) (bool, error) {
	//fmt.Printf("PC: %3d, instruction: %s, addr: %d, pre AC: %d, ", s.pc, s.opcode2mnemonic(opcode), addr, s.ac)
	switch opcode {
	case 0: // HLT
		s.hltVal = s.mem[addr]
		return true, nil
	case 1: // LDA
		s.ac = s.mem[addr]
		s.pc += 2
	case 2: // STA
		s.mem[addr] = s.ac
		s.pc += 2
	case 3: // ADD
		s.ac = s.ac + s.mem[addr]
		s.pc += 2
	case 4: // SUB
		s.ac = s.ac - s.mem[addr]
		s.pc += 2
	case 5: // AND
		s.ac &= s.mem[addr]
		s.pc += 2
	case 6: // INC
		s.mem[addr] = s.mem[addr] + 1
		s.pc += 2
	case 7: // JNZ
		// TODO: Rename to JNE?
		if s.ac != 0 {
			s.pc = addr
		} else {
			s.pc += 2
		}
	case 8: // DSZ
		s.mem[addr] = s.mem[addr] - 1
		if s.mem[addr] == 0 {
			s.pc += 4
		} else {
			s.pc += 2
		}
	case 9: // JMP
		s.pc = addr
	case 10: // SHL
		s.mem[addr] = s.mem[addr] << 1
		s.pc += 2
	case 11: // LDX
		s.x = s.mem[addr]
		s.pc += 2
	case 12: // LDY
		s.y = s.mem[addr]
		s.pc += 2
	case 13: // DYJNZ
		s.y = s.y - 1
		if s.y != 0 {
			s.pc = addr
		} else {
			s.pc += 2
		}
	case 14: // JSR - Jump to address, store return address in RET
		s.r = s.pc + 2
		s.pc = addr
	case 15: // RET - Jump to address in R
		// TODO: NOTE: This could also be used to loop on a JSR
		// TODO: NOTE: because the value won't change in R unless
		// TODO: NOTE: another RET is called
		s.pc = s.r
	case 16: // TAY - Transfer AC to Y
		s.y = s.ac
		s.pc += 2
	case 17: // STY - Store Y
		s.mem[addr] = s.y
		s.pc += 2
	case 18: // OR
		s.ac |= s.mem[addr]
		s.pc += 2
	case 19: // JEQ
		if s.ac == 0 {
			s.pc = addr
		} else {
			s.pc += 2
		}
	case 20: // JGT
		if s.ac > 0 {
			s.pc = addr
		} else {
			s.pc += 2
		}
	default:
		panic(fmt.Sprintf("unknown opcode: %d (%d)", opcode, (opcode&0x3f000000)>>24))
	}
	//fmt.Printf("post AC: %d\n", s.ac)
	return false, nil
}
