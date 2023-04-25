/*
 * Virtual machine with 1 operand
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm1

import "fmt"

// TODO: Make this configurable
const memSize = 32000

type VM1 struct {
	mem    [memSize]uint // Memory
	pc     uint          // Program Counter
	ac     uint          // 32-bit accumulator
	x      uint          // 32-bit index? register
	y      uint          // 32-bit index? register
	r      uint          // 32-bit return register
	hltVal uint          // A value returned by HLT
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

func (s *VM1) Mem() [memSize]uint {
	return s.mem
}

func (v *VM1) LoadRoutine(routine []uint) {
	copy(v.mem[:], routine)
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}

// fetch gets the next instruction from memory
// Returns: opcode, addr
// TODO: describe instruction format
func (s *VM1) fetch() (uint, uint, error) {
	if s.pc >= memSize {
		return 0, 0, fmt.Errorf("outside memory range: %d", s.pc)
	}
	ir := s.mem[s.pc]
	//opcode := (ir & 0x3F000000) >> 24
	opcode := (ir & 0xFF000000)
	//opcode := (ir & 0x3F000000)
	addr := ir & 0xFFFFFF
	/*
		baseIndex := ir & 0x40000000
		if baseIndex != 0 {
			baseIndirect := addr >> 12
			indexIndirect := addr & 0xFFF
			// TODO: Assume always at least 4096 memory to avoid check
			base := s.mem[baseIndirect]
			index := s.mem[indexIndirect]
			addr = base + index
		} else {
			indirect := ir & 0x80000000
			if indirect != 0 {
				if addr >= memSize {
					// TODO: Implement an error
					panic("outside memory range")
				}
				addr = s.mem[addr]
			}
		}
		if addr >= memSize {
			// TODO: Implement an error
			panic("outside memory range")
		}
	*/

	return opcode, addr, nil
}

// execute executes the supplied instruction
// Returns: hlt, error
func (s *VM1) execute(opcode uint, addr uint) (bool, error) {
	//	fmt.Printf("PC: %d, opcode: %d (%d), addr: %d, AC: %d\n", s.pc, opcode, (opcode&0x3F000000)>>24, addr, s.ac)
	switch opcode {
	case 0 << 24: // HLT
		s.hltVal = s.mem[addr]
		return true, nil
	case 1 << 24: // LDA
		s.ac = s.mem[addr]
		s.pc++
	case 2 << 24: // STA
		s.mem[addr] = s.ac
		s.pc++
	case 3 << 24: // ADD
		//fmt.Printf("PC: %d  ADD addr: %d\n", s.pc, addr)
		s.ac = mask32(s.ac + s.mem[addr])
		s.pc++
	case 4 << 24: // SUB
		s.ac = mask32(s.ac - s.mem[addr])
		s.pc++
	case 5 << 24: // AND
		s.ac &= s.mem[addr]
		s.pc++
	case 6 << 24: // INC
		s.mem[addr] = mask32(s.mem[addr] + 1)
		s.pc++
	case 7 << 24: // JNZ
		if s.ac != 0 {
			s.pc = addr
		} else {
			s.pc++
		}
	case 8 << 24: // LDO - Load offset
	case 9 << 24: // STA13 - Store least significant 13 bits
		s.mem[addr] = s.ac & 0o17777
		s.pc++
	case 10 << 24: // INC12 - Increment and store least significant 12 bits
		s.mem[addr] = (s.mem[addr] + 1) & 0o7777
		s.pc++
	case 11 << 24: // DSZ
		s.mem[addr] = mask32(s.mem[addr] - 1)
		if s.mem[addr] == 0 {
			s.pc += 2
		} else {
			s.pc++
		}
	case 12 << 24: // JMP
		s.pc = addr
	case 13 << 24: // SHL
		s.mem[addr] = mask32(s.mem[addr] << 1)
		s.pc++
	case 14 << 24: // STA12 - Store least significant 12 bits
		s.mem[addr] = s.ac & 0o7777
		s.pc++
	case 15 << 24: // LDX
		s.x = s.mem[addr]
		s.pc++
	case 16 << 24: // LDY
		s.y = s.mem[addr]
		s.pc++
	case 18 << 24: // DYJNZ
		//		fmt.Printf("PC: %d  DYJNZ  AC: %d, Y: %d\n", s.pc, s.ac, s.y)
		s.y = mask32(s.y - 1)
		if s.y != 0 {
			s.pc = addr
		} else {
			s.pc++
		}
	case 19 << 24: // STAD - Store A in first 24 bits of word
		//		fmt.Printf("PC: %d  STAD  preaddr: %d, ", s.pc, addr)
		s.mem[addr] = (s.mem[addr] & 0xFF000000) | (s.ac & 0xFFFFFF)
		//		fmt.Printf("postaddr: %d, ac: %d\n", addr, s.ac)
		s.pc++
	case 20 << 24: // JSR - Jump to address, store return address in RET
		//		fmt.Printf("PC: %d  JSR  AC: %d, Y: %d\n", s.pc, s.ac, s.y)
		s.r = mask32(s.pc + 1)
		s.pc = addr
	case 21 << 24: // RET - Jump to address in R
		// TODO: NOTE: This could also be used to loop on a JSR
		// TODO: NOTE: because the value won't change in R unless
		// TODO: NOTE: another RET is called
		s.pc = s.r
	case 22 << 24: // TAY - Transfer AC to Y
		s.y = s.ac
		s.pc++
	case 23 << 24: // STY - Store Y
		s.mem[addr] = s.y
		s.pc++
	case (1 | 0x80) << 24: // LDA I
		addr = s.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.ac = s.mem[addr]
		s.pc++
	case (1 | 0x40) << 24: // LDA II
		baseIndirect := addr >> 12
		indexIndirect := addr & 0xFFF
		// TODO: Assume always at least 4096 memory to avoid check
		base := s.mem[baseIndirect]
		index := s.mem[indexIndirect]
		//fmt.Printf("ADD II baseIndirect: %d, indexIndirect: %d, base: %d, index: %d\n", baseIndirect, indexIndirect, base, index)
		addr = base + index
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.ac = s.mem[addr]
		s.pc++

	case (2 | 0x80) << 24: // STA I
		addr = s.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.mem[addr] = s.ac
		s.pc++
	case (3 | 0x80) << 24: // ADD I
		// TODO: remove? experimental built-in in-direct
		addr = s.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.ac = mask32(s.ac + s.mem[addr])
		s.pc++
	case (3 | 0x40) << 24: // ADD II
		// TODO: remove? experimental built-in in-direct
		baseIndirect := addr >> 12
		indexIndirect := addr & 0xFFF
		// TODO: Assume always at least 4096 memory to avoid check
		base := s.mem[baseIndirect]
		index := s.mem[indexIndirect]
		//fmt.Printf("ADD II baseIndirect: %d, indexIndirect: %d, base: %d, index: %d\n", baseIndirect, indexIndirect, base, index)
		addr = base + index
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.ac = mask32(s.ac + s.mem[addr])
		s.pc++

	case (6 | 0x80) << 24: // INC I
		addr = s.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.mem[addr] = mask32(s.mem[addr] + 1)
		s.pc++
	case (10 | 0x40) << 24: // INC12 II - Increment and store least significant 12 bits
		baseIndirect := addr >> 12
		indexIndirect := addr & 0xFFF
		// TODO: Assume always at least 4096 memory to avoid check
		base := s.mem[baseIndirect]
		index := s.mem[indexIndirect]
		//fmt.Printf("ADD II baseIndirect: %d, indexIndirect: %d, base: %d, index: %d\n", baseIndirect, indexIndirect, base, index)
		addr = base + index
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.mem[addr] = (s.mem[addr] + 1) & 0o7777
		s.pc++
	case (10 | 0x80) << 24: // INC12 I - Increment and store least significant 12 bits
		addr = s.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.mem[addr] = (s.mem[addr] + 1) & 0o7777
		s.pc++
	case (12 | 0x40) << 24: // JMP II
		baseIndirect := addr >> 12
		indexIndirect := addr & 0xFFF
		// TODO: Assume always at least 4096 memory to avoid check
		base := s.mem[baseIndirect]
		index := s.mem[indexIndirect]
		//fmt.Printf("ADD II baseIndirect: %d, indexIndirect: %d, base: %d, index: %d\n", baseIndirect, indexIndirect, base, index)
		addr = base + index
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.pc = addr
	case (17) << 24: // ADD IX
		// TODO: Clarify encoding as could be:
		// TODO:   addr = mem[addr+x]  or addr = mem[addr]+x
		// TODO: Currently:
		// TODO:   addr = mem[addr+x]
		// TODO: Encode the X in the bitmap for general use
		// TODO: This isn't implement the same as the other II
		// TODO: instructions as it is in flux
		// TODO: Assume always at least 4096 memory to avoid check
		base := s.mem[addr]
		//		fmt.Printf("ADD  IX addr: %d, x: %d, base: %d\n", addr, s.x, base)
		addr = base + s.x
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		s.ac = mask32(s.ac + s.mem[addr])
		s.pc++
	default:
		panic(fmt.Sprintf("unknown opcode: %d (%d)", opcode, (opcode&0x3f000000)>>24))
	}
	return false, nil
}
