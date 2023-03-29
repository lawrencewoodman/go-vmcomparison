/*
 * Simple Implementation Stack Virtual Machine using a stack
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package vmstack

// TODO: Make this configurable
const memSize = 32000

type VMStack struct {
	mem   [memSize]uint // Memory
	pc    uint          // Program Counter
	stack *LStack       // 8 element limited stack
	// stack  *CStack // 8 element circular stack
	hltVal uint // A value returned by HLT
	r      uint
}

func New() *VMStack {
	return &VMStack{stack: NewLStack()}
	// return &VMStack{stack: NewCStack()}
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

func (v *VMStack) Mem() [memSize]uint {
	return v.mem
}

func (v *VMStack) LoadRoutine(routine []uint) {
	copy(v.mem[:], routine)
}

// Returns: hlt, error
func (v *VMStack) Step() (bool, error) {
	ir := v.mem[v.pc]
	opcode := (ir & 0xFF000000)
	//	fmt.Printf("PC: %d, opcode: %d (%d)\n", v.pc, opcode, opcode>>24)
	// TODO: do something with operand for LIT, STORE, FETCH, ADD?
	switch opcode {
	case 0 << 24: // HLT
		v.hltVal = v.stack.pop()
		return true, nil
	case 1 << 24: // FETCH
		// TODO: check memory in range
		v.stack.replace(v.mem[v.stack.peek()])
		v.pc = mask32(v.pc + 1)
	case 2 << 24: // STORE
		// TODO: check memory in range
		v.mem[v.stack.pop()] = v.stack.pop()
		v.pc = mask32(v.pc + 1)
	case 3 << 24: // ADD
		a := v.stack.pop()
		b := v.stack.peek()
		c := mask32(a + b)
		v.stack.replace(c)
		//fmt.Printf("PC: %d  ADD %d + %d = %d\n", v.pc, a, b, c)
		v.pc = mask32(v.pc + 1)
	case 4 << 24: // SUB
		v.stack.replace(mask32(v.stack.pop() - v.stack.peek()))
		v.pc = mask32(v.pc + 1)
	case 5 << 24: // AND
		v.stack.replace(v.stack.pop() & v.stack.peek())
		v.pc = mask32(v.pc + 1)
	case 6 << 24: // INC
		v.stack.replace(mask32(v.stack.peek() + 1))
		v.pc = mask32(v.pc + 1)
	case 7 << 24: // JNZ
		if v.stack.pop() != 0 {
			v.pc = v.stack.pop()
		} else {
			v.pc = mask32(v.pc + 1)
		}
	case 9 << 24: // STORE13 - Store least significant 13 bits
		v.mem[v.stack.pop()] = v.stack.pop() & 0o17777
		v.pc = mask32(v.pc + 1)
	case 10 << 24: // INC12 - Increment and store least significant 12 bits
		v.stack.replace((v.stack.peek() + 1) & 0o7777)
		v.pc = mask32(v.pc + 1)
	case 11 << 24: // DJNZ - (val addr -- val) - Decrement and Jump if not Zero
		addr := v.stack.pop()
		val := v.stack.peek()
		val = mask32(val - 1)
		v.stack.replace(val)
		if val != 0 {
			v.pc = addr
		} else {
			v.pc = mask32(v.pc + 1)
		}
	case 12 << 24: // JMP
		v.pc = v.stack.pop()
	case 13 << 24: // SHL
		v.stack.replace(mask32(v.stack.peek() << 1))
		v.pc = mask32(v.pc + 1)
	case 14 << 24: // STORE12 - Store least significant 12 bits
		addr := v.stack.pop()
		val := v.stack.pop()
		v.mem[addr] = val & 0o7777
		//		fmt.Printf("PC: %d  STORE12 mem[%d] = %d\n", v.pc, addr, val)
		v.pc = mask32(v.pc + 1)
	case 15 << 24: // LITO - Put the 24-bit operand on the stack
		literal := ir & 0x00FFFFFF
		v.stack.push(literal)
		v.pc = mask32(v.pc + 1)
	case 16 << 24: // FETCHO
		addr := ir & 0x00FFFFFF
		v.stack.push(v.mem[addr])
		v.pc = mask32(v.pc + 1)
	case 17 << 24: // DJNZO - op: addr (val -- val) - Decrement and Jump if not Zero
		addr := ir & 0x00FFFFFF
		val := v.stack.peek()
		//		fmt.Printf("PC: %d  DJNZO addr: %d, val: %d\n", v.pc, addr, val)

		val = mask32(val - 1)
		v.stack.replace(val)
		if val != 0 {
			v.pc = addr
		} else {
			v.pc = mask32(v.pc + 1)
		}
	case 18 << 24: // DROP - (n --)
		v.stack.pop()
		//fmt.Printf("PC: %d  DROP %d\n", v.pc, a)
		v.pc = mask32(v.pc + 1)
	case 19 << 24: // SWAP - (a b -- b a)
		a := v.stack.pop()
		b := v.stack.peek()
		v.stack.replace(a)
		v.stack.push(b)
		// TODO: see if we can use peek/replace here
		//fmt.Printf("PC: %d  SWAP (%d %d -- %d %d\n", v.pc, a, b, b, a)
		v.pc = mask32(v.pc + 1)
	case 20 << 24: // FETCHBI - (base index -- n)
		// TODO: check memory in range
		addr := v.stack.pop() + v.stack.peek()
		v.stack.replace(v.mem[addr])
		v.pc = mask32(v.pc + 1)
	case 21 << 24: // STOREO - op: addr (val --)
		// TODO: check memory in range
		addr := ir & 0x00FFFFFF
		v.mem[addr] = v.stack.pop()
		v.pc = mask32(v.pc + 1)
	case 22 << 24: // DSZO - op: addr (--)
		addr := ir & 0x00FFFFFF
		v.mem[addr] = mask32(v.mem[addr] - 1)
		if v.mem[addr] == 0 {
			v.pc = mask32(v.pc + 2)
		} else {
			v.pc = mask32(v.pc + 1)
		}
	case 23 << 24: // JMPO - op: addr (--)
		addr := ir & 0x00FFFFFF
		v.pc = addr
	case 24 << 24: // ADDBI - (n base index -- n)
		// TODO: check memory in range
		addr := v.stack.pop() + v.stack.pop()
		val := mask32(v.mem[addr] + v.stack.peek())
		//		fmt.Printf("PC: %d  ADDBI addr: %d, newVal: %d\n", v.pc, addr, val)
		v.stack.replace(val)
		v.pc = mask32(v.pc + 1)
	case 25 << 24: // R_PUSH - push TOS to R
		v.r = v.stack.pop()
		v.pc = mask32(v.pc + 1)
	case 26 << 24: // R_POP - POP R to TOS
		v.stack.push(v.r)
		v.pc = mask32(v.pc + 1)
	case 27 << 24: // FETCHI
		// TODO: check memory in range
		addr := v.mem[v.stack.peek()]
		v.stack.replace(v.mem[addr])
		v.pc = mask32(v.pc + 1)
	}
	return false, nil
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}
