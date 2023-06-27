/*
 * Simple Implementation Stack Virtual Machine using a stack
 * This version uses the bottom 24 bits of the word as if it
 * is the TOS if > 0
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package bvmstack

import (
	"fmt"
	"math/big"
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
	var zero = big.NewInt(0)
	var one = big.NewInt(1)

	if v.pc >= memSize {
		return false, fmt.Errorf("outside memory range: %d", v.pc)
	}
	ir := v.mem[v.pc]
	opcode := (ir & 0xFF000000)
	operand := (ir & 0x00FFFFFF)
	//fmt.Printf("PC: %3d, instruction: %5s, operand: %4d, pre NOS: %5s, pre TOS: %5d, post TOS:", v.pc, v.opcode2mnemonic(opcode), operand, v.dstack.nos(), v.dstack.peek())

	if operand > 0 {
		v.dstack.push(big.NewInt(operand))
	}

	switch opcode {
	case 0 << 24: // HLT
		t := v.dstack.pop()
		if !t.IsInt64() {
			panic("halt value not in64")
		}
		v.hltVal = t.Int64()
		return true, nil
	case 1 << 24: // FETCH
		t := v.dstack.peek()
		if !t.IsInt64() {
			panic("address is not in64")
		}
		addr := t.Int64()
		if addr >= memSize {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		v.dstack.replace(big.NewInt(v.mem[addr]))
		v.pc++
	case 2 << 24: // STORE (n addr --)
		t := v.dstack.pop()
		if !t.IsInt64() {
			panic("address is not in64")
		}
		addr := t.Int64()
		if addr >= memSize {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}

		t = v.dstack.pop()
		if !t.IsInt64() {
			panic("value is not in64")
		}

		v.mem[addr] = t.Int64()
		v.pc++
	case 3 << 24: // ADD
		a := v.dstack.pop()
		b := v.dstack.peek()
		a.Add(a, b)
		v.dstack.replace(a)
		v.pc++
	case 4 << 24: // SUB (a b -- a-b)
		b := v.dstack.pop()
		a := v.dstack.peek()
		// TODO: Work out why can't just use a.Sub(a,b)  and replace(a)
		//a.Sub(a, b)
		c := big.NewInt(0).Set(a)
		c = c.Sub(a, b)
		v.dstack.replace(c)
		v.pc++
	case 5 << 24: // AND
		a := v.dstack.pop()
		b := v.dstack.peek()
		v.dstack.replace(a.And(a, b))
		v.pc++
	case 6 << 24: // INC
		a := v.dstack.peek()
		v.dstack.replace(a.Add(a, one))
		v.pc++
	case 7 << 24: // JNZ (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val.Sign() != 0 {
			if !addr.IsInt64() {
				panic("address is not in64")
			}
			v.pc = addr.Int64()
		} else {
			v.pc++
		}
	case 8 << 24: // DJNZ - (val addr -- val) - Decrement and Jump if not Zero
		addr := v.dstack.pop()
		val := v.dstack.peek()
		v.dstack.replace(val.Sub(val, one))
		if val.Sign() != 0 {
			if !addr.IsInt64() {
				panic("address is not in64")
			}

			v.pc = addr.Int64()
		} else {
			v.pc++
		}
	case 9 << 24: // JMP
		addr := v.dstack.pop()
		if !addr.IsInt64() {
			panic("address is not in64")
		}
		v.pc = addr.Int64()
	case 10 << 24: // SHL
		// TODO: Be able to supply number of bits to shift on stack?
		a := v.dstack.peek()
		v.dstack.replace(a.Lsh(a, 1))
		v.pc++
	case 11 << 24: // LIT - Put the 24-bit operand on the stack
		if operand == 0 {
			v.dstack.push(zero)
		}
		// else operand is pushed to TOS at start
		// of this function
		v.pc++
	case 12 << 24: // DROP - (n --)
		v.dstack.pop()
		v.pc++
	case 13 << 24: // SWAP - (a b -- b a)
		b := v.dstack.pop()
		a := v.dstack.peek()
		v.dstack.replace(b)
		v.dstack.push(a)
		v.pc++
	case 17 << 24: // JSR
		pc := big.NewInt(v.pc)
		v.rstack.push(pc.Add(pc, one))
		addr := v.dstack.pop()
		if !addr.IsInt64() {
			panic("address is not in64")
		}
		v.pc = addr.Int64()
	case 18 << 24: // RET
		addr := v.rstack.pop()
		if !addr.IsInt64() {
			panic("address is not in64")
		}
		v.pc = addr.Int64()
	case 19 << 24: // DUP
		v.dstack.push(v.dstack.peek())
		v.pc++
	case 20 << 24: // OR
		a := v.dstack.pop()
		b := v.dstack.peek()
		v.dstack.replace(a.Or(a, b))
		v.pc++
	case 21 << 24: // JZ (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val.Sign() == 0 {
			if !addr.IsInt64() {
				panic("address is not in64")
			}
			v.pc = addr.Int64()
		} else {
			v.pc++
		}
	case 22 << 24: // JGT (val addr --)
		addr := v.dstack.pop()
		val := v.dstack.pop()
		if val.Sign() > 0 {
			if !addr.IsInt64() {
				panic("address is not in64")
			}

			v.pc = addr.Int64()
		} else {
			v.pc++
		}
	case 23 << 24: // ROT (a b c -- b c a)
		c := v.dstack.pop()
		b := v.dstack.pop()
		a := v.dstack.peek()
		v.dstack.replace(b)
		v.dstack.push(c)
		v.dstack.push(a)
		v.pc++
	case 24 << 24: // OVER (a b -- a b a)
		b := v.dstack.pop()
		a := v.dstack.peek()
		v.dstack.push(b)
		v.dstack.push(a)
		v.pc++
	default:
		panic("unknown opcode")
	}
	//fmt.Printf("post TOS: %d\n", v.dstack.peek())
	return false, nil
}
