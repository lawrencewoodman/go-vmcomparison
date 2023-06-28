/*
 * Simple Implementation Stack Virtual Machine using a stack
 * This version uses the bottom 24 bits of the word as if it
 * is the TOS if > 0
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package vmstack

import (
	"fmt"
)

// TODO: Make this configurable
const memSize = 32000

type VMStack struct {
	mem    [memSize]int64 // Memory
	pc     int64          // Program Counter
	dstack *LStack        // 8 element limited data stack
	// stack  *CStack // 8 element circular data stack
	rstack *LStack // 8 element limited return
	hltVal int64   // A value returned by HLT
}

func New() *VMStack {
	return &VMStack{dstack: NewLStack(), rstack: NewLStack()}
	// return &VMStack2{stack: NewCStack()}
}

func (v *VMStack) Run() (bool, error) {
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

func (v *VMStack) Mem() [memSize]int64 {
	return v.mem
}

func (v *VMStack) LoadRoutine(routine []int64) {
	copy(v.mem[:], routine)
}

func (v *VMStack) opcode2mnemonic(opcode int64) string {
	for m, o := range instructions {
		if o == opcode {
			return m
		}
	}
	panic("opcode not found")
}

// Returns: hlt, error
func (v *VMStack) Step() (bool, error) {
	if v.pc >= memSize {
		return false, fmt.Errorf("outside memory range: %d", v.pc)
	}
	ir := v.mem[v.pc]
	opcode := (ir & 0xFF000000)
	operand := (ir & 0x00FFFFFF)

	//fmt.Printf("%6d: %5s %-7d   (%5s %5d -- ", v.pc, v.opcode2mnemonic(opcode), operand, v.dstack.nos(), v.dstack.peek())

	if operand > 0 {
		v.dstack.push(operand)
	}
	switch opcode {
	case 0 << 24: // HLT
		v.hltVal = v.dstack.pop()
		return true, nil
	case 1 << 24: // FETCH
		addr := v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		v.dstack.replace(v.mem[addr])
		v.pc++
	case 2 << 24: // STORE (n addr --)
		addr := v.dstack.pop()
		if addr >= memSize {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}

		v.mem[addr] = v.dstack.pop()
		v.pc++
	case 3 << 24: // ADD
		a := v.dstack.pop()
		b := v.dstack.peek()
		v.dstack.replace(a + b)
		v.pc++
	case 4 << 24: // SUB (a b -- a-b)
		b := v.dstack.pop()
		a := v.dstack.peek()
		v.dstack.replace(a - b)
		v.pc++
	case 5 << 24: // AND
		v.dstack.replace(v.dstack.pop() & v.dstack.peek())
		v.pc++
	case 6 << 24: // INC
		v.dstack.replace(v.dstack.peek() + 1)
		v.pc++
	case 7 << 24: // JNZ (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val != 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 8 << 24: // DJNZ - (val addr -- val) - Decrement and Jump if not Zero
		addr := v.dstack.pop()
		val := v.dstack.peek()
		val = val - 1
		v.dstack.replace(val)
		if val != 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 9 << 24: // JMP
		v.pc = v.dstack.pop()
	case 10 << 24: // SHL
		v.dstack.replace(v.dstack.peek() << 1)
		v.pc++
	case 11 << 24: // LIT - Put the 24-bit operand on the stack
		if operand == 0 {
			v.dstack.push(0)
		}
		// else operand is pushed to TOS at start
		// of this function
		v.pc++
	case 12 << 24: // DROP - (n --)
		v.dstack.drop()
		v.pc++
	case 13 << 24: // SWAP - (a b -- b a)
		v.dstack.swap()
		v.pc++
	case 14 << 24: // FETCHBI - (base index -- n)
		addr := v.dstack.pop() + v.dstack.peek()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		v.dstack.replace(v.mem[addr])
		v.pc++
	case 15 << 24: // ADDBI - (n base index -- n)
		addr := v.dstack.pop() + v.dstack.pop()
		if addr >= memSize {
			return false, fmt.Errorf("outside memory range: %d", addr)
		}
		val := v.mem[addr] + v.dstack.peek()
		v.dstack.replace(val)
		v.pc++
	case 16 << 24: // FETCHI
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
	case 17 << 24: // JSR
		v.rstack.push(v.pc + 1)
		v.pc = v.dstack.pop()
	case 18 << 24: // RET
		v.pc = v.rstack.pop()
	case 19 << 24: // DUP
		v.dstack.dup()
		v.pc++
	case 20 << 24: // OR
		v.dstack.replace(v.dstack.pop() | v.dstack.peek())
		v.pc++
	case 21 << 24: // JZ (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val == 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 22 << 24: // JGT (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val > 0 {
			v.pc = addr
		} else {
			v.pc++
		}
	case 23 << 24: // ROT (a b c -- b c a)
		v.dstack.rot()
		v.pc++
	case 24 << 24: // OVER (a b -- a b a)
		v.dstack.over()
		v.pc++
	}

	//fmt.Printf("%5s %-7d)\n", v.dstack.nos(), v.dstack.peek())
	return false, nil
}
