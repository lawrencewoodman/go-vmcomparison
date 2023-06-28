/*
 * Simple Implementation Virtual Machine v2 using Big Numbers
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package bvm2

import (
	"fmt"
	"math/big"
)

// TODO: Make this configurable
const memSize = 32000

type VM2 struct {
	code        [memSize]int64    // Code / Program
	mem         [memSize]*big.Int // Memory
	pc          int64             // Program Counter
	hltVal      *big.Int          // A value returned by HLT
	codeSymbols map[string]int64  // The code symbols table from the assembler - to aid debugging
	dataSymbols map[string]int64  // The data symbols table from the assembler - to aid debugging
}

func New() *VM2 {
	var mem [memSize]*big.Int
	for i := 0; i < memSize; i++ {
		mem[i] = big.NewInt(0)

	}
	return &VM2{mem: mem}
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

func (v *VM2) Mem() [memSize]*big.Int {
	return v.mem
}

func (v *VM2) LoadRoutine(code []int64, data []*big.Int, codeSymbols, dataSymbols map[string]int64) {
	// Need to copy the individual data points of the routine because they are pointers
	for i, d := range data {
		v.mem[i].Set(d)
	}
	copy(v.code[:], code)
	v.codeSymbols = codeSymbols
	v.dataSymbols = dataSymbols
}

// fetch gets the next instruction from memory
// Returns: opcode, operandA, operandB
// TODO: describe instruction format
func (v *VM2) fetch() (int64, int64, int64, error) {
	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, v.pc+1)
	}
	opcode := v.code[v.pc]
	operandA := v.code[v.pc+1]
	operandB := v.code[v.pc+2]

	// If addressing mode: operand A indirect
	if operandA < 0 {
		operandA = -operandA
		if operandA >= memSize {
			return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
		}
		iOperandA := v.mem[operandA]
		if !iOperandA.IsInt64() {
			return 0, 0, 0, fmt.Errorf("PC: %d, operand A not int64", v.pc)
		}
		operandA = iOperandA.Int64()
	}
	if operandA >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
	}

	// If addressing mode: operand B indirect
	if operandB < 0 {
		operandB = -operandB
		if operandB >= memSize {
			return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
		}
		iOperandB := v.mem[operandB]
		if !iOperandB.IsInt64() {
			return 0, 0, 0, fmt.Errorf("PC: %d, operand B not int64", v.pc)
		}
		operandB = iOperandB.Int64()
	}
	if operandB >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
	}
	return opcode, operandA, operandB, nil
}

func (v *VM2) addr2symbol(addr int64, onlyCode ...bool) string {
	if len(onlyCode) == 0 {
		for k, v := range v.dataSymbols {
			if v == addr {
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

func (v *VM2) opcode2mnemonic(opcode int64) string {
	for m, o := range instructions {
		if o == opcode {
			return m
		}
	}
	panic("opcode not found")
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *VM2) execute(opcode int64, operandA int64, operandB int64) (bool, error) {
	//fmt.Printf("%7s:    %s   %s, %s\n", v.addr2symbol(v.pc, true), v.opcode2mnemonic(opcode), v.addr2symbol(operandA), v.addr2symbol(operandB))
	//fmt.Printf("            pre:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])
	var one = big.NewInt(1)
	switch opcode {
	case 0: // HLT
		v.hltVal = v.mem[operandA]
		// TODO: this wastes the following memory location, should it?
		return true, nil
	case 1: // MOV
		v.mem[operandB].Set(v.mem[operandA])
		v.pc += 3
	case 2: // JSR
		v.mem[operandB] = big.NewInt(v.pc + 3)
		v.pc = operandA
	case 3: // ADD
		v.mem[operandB].Add(v.mem[operandB], v.mem[operandA])
		v.pc += 3
	case 4: // DJNZ
		v.mem[operandA].Sub(v.mem[operandA], one)
		if v.mem[operandA].Sign() != 0 {
			v.pc = operandB
		} else {
			v.pc += 3
		}
	case 5: // JMP
		v.pc = operandA + operandB
	case 6: // AND
		v.mem[operandB].And(v.mem[operandB], v.mem[operandA])
		v.pc += 3
	case 7: // OR
		v.mem[operandB].Or(v.mem[operandB], v.mem[operandA])
		v.pc += 3
	case 8: // SHL
		if !v.mem[operandA].IsUint64() {
			return false, fmt.Errorf("PC: %d, mem[A] is not uint64", v.pc)
		}
		v.mem[operandB].Lsh(v.mem[operandB], uint(v.mem[operandA].Uint64()))
		v.pc += 3
	case 9: // JNZ
		if v.mem[operandA].Sign() != 0 {
			v.pc = operandB
		} else {
			v.pc += 3
		}
	case 10: // SNE
		if v.mem[operandA].Cmp(v.mem[operandB]) != 0 {
			v.pc += 6
		} else {
			v.pc += 3
		}
	case 11: // SLE
		if v.mem[operandA].Cmp(v.mem[operandB]) <= 0 {
			v.pc += 6
		} else {
			v.pc += 3
		}
	case 12: // SUB
		v.mem[operandB].Sub(v.mem[operandB], v.mem[operandA])
		v.pc += 3
	case 13: // JGT
		if v.mem[operandA].Sign() == 1 {
			v.pc = operandB
		} else {
			v.pc += 3
		}

	default:
		return false, fmt.Errorf("unknown opcode: %d (%d)", opcode, (opcode&0x3f000000)>>24)
	}

	//	fmt.Printf("            post:  [%s]: %d, [%s]: %d\n", v.addr2symbol(operandA), v.mem[operandA], v.addr2symbol(operandB), v.mem[operandB])
	return false, nil
}
