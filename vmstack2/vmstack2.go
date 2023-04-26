/*
 * Simple Implementation Stack Virtual Machine using a stack
 * This version uses the bottom 24 bits of the word as if it
 * is the TOS if > 0
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package vmstack2

import "fmt"

// TODO: Make this configurable
const memSize = 32000

type VMStack2 struct {
	mem    [memSize]uint // Memory
	pc     uint          // Program Counter
	dstack *LStack       // 8 element limited data stack
	// stack  *CStack // 8 element circular data stack
	rstack *LStack // 8 element limited return
	hltVal uint    // A value returned by HLT
}

func New() *VMStack2 {
	return &VMStack2{dstack: NewLStack(), rstack: NewLStack()}
	// return &VMStack2{stack: NewCStack()}
}

func (v *VMStack2) Run() (bool, error) {
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

func (v *VMStack2) Mem() [memSize]uint {
	return v.mem
}

func (v *VMStack2) LoadRoutine(routine []uint) {
	copy(v.mem[:], routine)
}

// Returns: hlt, error
func (v *VMStack2) Step() (bool, error) {
	if v.pc >= memSize {
		return false, fmt.Errorf("outside memory range: %d", v.pc)
	}
	ir := v.mem[v.pc]
	opcode := (ir & 0xFF000000)
	operand := (ir & 0x00FFFFFF)
	if operand > 0 {
		v.dstack.push(operand)
	}
	//	fmt.Printf("PC: %d, opcode: %d (%d)\n", v.pc, opcode, opcode>>24)
	// TODO: do something with operand for LIT, STORE, FETCH, ADD?
	switch opcode {
	case 0 << 24: // HLT
		//		if operand > 0 {
		//			v.hltVal = operand
		//		} else {
		v.hltVal = v.dstack.pop()
		//		}
		return true, nil
	case 1 << 24: // FETCH
		//		if operand > 0 {
		//			if operand >= memSize {
		//				return false, fmt.Errorf("outside memory range: %d", operand)
		//			}
		//			v.dstack.push(v.mem[operand])
		//		} else {
		addr := v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		v.dstack.replace(v.mem[addr])
		//		}
		v.pc++
	case 2 << 24: // STORE (n addr --)
		//		if operand > 0 {
		//			addr = operand
		//		} else {
		addr := v.dstack.pop()
		//		}
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}

		v.mem[addr] = v.dstack.pop()
		//		fmt.Printf("PC: %d  STORE n:%d addr:%d\n", v.pc, v.mem[addr], addr)
		v.pc++
	case 3 << 24: // ADD
		//		a := operand
		//		if operand == 0 {
		a := v.dstack.pop()
		//		}
		b := v.dstack.peek()
		c := mask32(a + b)
		v.dstack.replace(c)
		//fmt.Printf("PC: %d  ADD %d + %d = %d\n", v.pc, a, b, c)
		v.pc++
	case 4 << 24: // SUB
		v.dstack.replace(mask32(v.dstack.pop() - v.dstack.peek()))
		v.pc++
	case 5 << 24: // AND
		//		if operand > 0 {
		//			v.dstack.replace(operand & v.dstack.peek())
		//		} else {
		v.dstack.replace(v.dstack.pop() & v.dstack.peek())
		//		}
		v.pc++
	case 6 << 24: // INC
		v.dstack.replace(mask32(v.dstack.peek() + 1))
		v.pc++
	case 7 << 24: // JNZ (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val != 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 9 << 24: // STORE13 - Store least significant 13 bits
		//		if operand > 0 {
		//			addr = operand
		//		} else {
		addr := v.dstack.pop()
		//		}
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		v.mem[addr] = v.dstack.pop() & 0o17777
		v.pc++
	case 10 << 24: // INC12 - Increment and store least significant 12 bits
		v.dstack.replace((v.dstack.peek() + 1) & 0o7777)
		v.pc++
	case 11 << 24: // DJNZ - (val addr -- val) - Decrement and Jump if not Zero
		//		if operand > 0 {
		//			addr = operand
		//		} else {
		addr := v.dstack.pop()
		//		}
		val := v.dstack.peek()
		val = mask32(val - 1)
		v.dstack.replace(val)
		if val != 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 12 << 24: // JMP
		//		if operand > 0 {
		//			v.pc = operand
		//		} else {
		v.pc = v.dstack.pop()
		//		}
	case 13 << 24: // SHL
		v.dstack.replace(mask32(v.dstack.peek() << 1))
		v.pc++
	case 14 << 24: // STORE12 - Store least significant 12 bits
		addr := v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		val := v.dstack.pop()
		v.mem[addr] = val & 0o7777
		//		fmt.Printf("PC: %d  STORE12 mem[%d] = %d\n", v.pc, addr, val)
		v.pc++
	case 15 << 24: // LIT - Put the 24-bit operand on the stack
		if operand == 0 {
			v.dstack.push(0)
		}
		// else operand is pushed to TOS at start
		// of this function
		v.pc++
	case 16 << 24:
		panic("TODO: remove this")
	case 17 << 24:
		panic("TODO: remove this")
	case 18 << 24: // DROP - (n --)
		v.dstack.pop()
		//fmt.Printf("PC: %d  DROP %d\n", v.pc, a)
		v.pc++
	case 19 << 24: // SWAP - (a b -- b a)
		a := v.dstack.pop()
		b := v.dstack.peek()
		v.dstack.replace(a)
		v.dstack.push(b)
		// TODO: see if we can use peek/replace here
		//fmt.Printf("PC: %d  SWAP (%d %d -- %d %d\n", v.pc, a, b, b, a)
		v.pc++
	case 20 << 24: // FETCHBI - (base index -- n)
		addr := v.dstack.pop() + v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		v.dstack.replace(v.mem[addr])
		v.pc++
	case 21 << 24:
		panic("TODO: Remove this")
	case 22 << 24:
		panic("TODO: Remove this")
	case 23 << 24:
		panic("TODO: Remove this")
	case 24 << 24: // ADDBI - (n base index -- n)
		addr := v.dstack.pop() + v.dstack.pop()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		val := mask32(v.mem[addr] + v.dstack.peek())
		//		fmt.Printf("PC: %d  ADDBI addr: %d, newVal: %d\n", v.pc, addr, val)
		v.dstack.replace(val)
		v.pc++
	case 25 << 24:
		panic("TODO: remove this")
	case 26 << 24:
		panic("TODO: remove this")
	case 27 << 24: // FETCHI
		addr := v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		addr = v.mem[addr]
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		v.dstack.replace(v.mem[addr])
		v.pc++
	case 28 << 24: // JSR
		v.rstack.push(mask32(v.pc + 1))
		//		if operand > 0 {
		//			v.pc = operand
		//		} else {
		v.pc = v.dstack.pop()
		//		}
	case 29 << 24: // RET
		v.pc = v.rstack.pop()

	case 30 << 24: // DUP
		v.dstack.push(v.dstack.peek())
		v.pc++
	}
	return false, nil
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}
