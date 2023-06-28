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
	code   [memSize]int64    // Code / Program
	mem    [memSize]*big.Int // Memory
	pc     int64             // Program Counter
	dstack *LStack           // 8 element limited data stack
	// stack  *CStack // 8 element circular data stack
	rstack      *LStack          // 8 element limited return
	hltVal      *big.Int         // A value returned by HLT
	codeSymbols map[string]int64 // The code symbols table from the assembler - to aid debugging
	dataSymbols map[string]int64 // The data symbols table from the assembler - to aid debugging
}

func New() *VMStack {
	var mem [memSize]*big.Int
	for i := 0; i < memSize; i++ {
		mem[i] = big.NewInt(0)

	}
	return &VMStack{mem: mem, dstack: NewLStack(), rstack: NewLStack()}
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

func (v *VMStack) Mem() [memSize]*big.Int {
	return v.mem
}

func (v *VMStack) LoadRoutine(code []int64, data []*big.Int, codeSymbols, dataSymbols map[string]int64) {
	// Need to copy the individual data points of the routine because they are pointers
	for i, d := range data {
		v.mem[i].Set(d)
	}
	copy(v.code[:], code)
	v.codeSymbols = codeSymbols
	v.dataSymbols = dataSymbols
}

func (v *VMStack) addr2symbol(addr int64, onlyCode ...bool) string {
	if len(onlyCode) == 0 {
		for k, v := range v.dataSymbols {
			if v == addr && addr != 0 {
				return k
			}
		}
	}

	for k, v := range v.codeSymbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}

func (v *VMStack) opcode2mnemonic(opcode int64) string {
	for m, o := range instructions {
		if o == opcode {
			return m
		}
	}
	panic("opcode not found")
}

var zero = big.NewInt(0)
var one = big.NewInt(1)

// TODO: Enforce memSize  being int64
var bmemSize = big.NewInt(memSize)

// Returns: hlt, error
func (v *VMStack) Step() (bool, error) {
	if v.pc >= memSize {
		return false, fmt.Errorf("outside memory range: %d", v.pc)
	}
	ir := v.code[v.pc]
	opcode := (ir & 0xFF000000)
	operand := (ir & 0x00FFFFFF)

	//fmt.Printf("%6s: %5s %-7s   (%5s %5d -- ", v.addr2symbol(v.pc, true), v.opcode2mnemonic(opcode), v.addr2symbol(operand), v.dstack.nos(), v.dstack.peek())

	if operand > 0 {
		v.dstack.push(big.NewInt(operand))
	}

	switch opcode {
	case 0 << 24: // HLT
		v.hltVal = v.dstack.pop()
		return true, nil
	case 1 << 24: // FETCH
		addr := v.dstack.peek()
		// if addr > memSize
		if addr.Cmp(bmemSize) >= 0 {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		addr.Set(v.mem[addr.Int64()])
		v.pc++
	case 2 << 24: // STORE (n addr --)
		addr := v.dstack.pop()
		// if addr > memSize
		if addr.Cmp(bmemSize) >= 0 {
			return false, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		v.mem[addr.Int64()] = v.dstack.pop()
		v.pc++
	case 3 << 24: // ADD
		a := v.dstack.pop()
		b := v.dstack.peek()
		b.Add(a, b)
		v.pc++
	case 4 << 24: // SUB (a b -- a-b)
		b := v.dstack.pop()
		a := v.dstack.peek()
		a.Sub(a, b)
		v.pc++
	case 5 << 24: // AND
		a := v.dstack.pop()
		b := v.dstack.peek()
		b.And(a, b)
		v.pc++
	case 6 << 24: // INC
		a := v.dstack.peek()
		a.Add(a, one)
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
		val.Sub(val, one)
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
		a.Lsh(a, 1)
		v.pc++
	case 11 << 24: // LIT - Put the 24-bit operand on the stack
		if operand == 0 {
			v.dstack.push(zero)
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
		v.dstack.dup()

		v.pc++
	case 20 << 24: // OR
		a := v.dstack.pop()
		b := v.dstack.peek()
		b.Or(a, b)
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
		v.dstack.rot()
		v.pc++
	case 24 << 24: // OVER (a b -- a b a)
		v.dstack.over()
		v.pc++
	default:
		panic("unknown opcode")
	}

	//fmt.Printf("%5s %5d)\n", v.dstack.nos(), v.dstack.peek())
	return false, nil
}
