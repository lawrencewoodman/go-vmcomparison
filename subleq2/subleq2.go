/*
 * SUBLEQ2 Virtual Machine
 *
 * SUBLEQ2 understands negative numbers passed as operands to be
 * indirect addresses.
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package subleq2

import (
	"fmt"
)

// TODO: Make this configurable
const memSize = 32000

// Location in memory of hltVal
// If this is used as a destination location then a HLT is executed
const hltLoc = 1000

type SUBLEQ struct {
	code        [memSize]int64   // Code / Program
	mem         [memSize]int64   // Memory
	pc          int64            // Program Counter
	hltVal      int64            // A value returned by HLT
	codeSymbols map[string]int64 // The code symbols table from the assembler - to aid debugging
	dataSymbols map[string]int64 // The data symbols table from the assembler - to aid debugging
	codeSize    int64            // The size of the code / program

}

func New() *SUBLEQ {
	return &SUBLEQ{}
}

func (v *SUBLEQ) Run() error {
	hlt := false
	for !hlt {
		if v.pc+2 >= v.codeSize {
			return fmt.Errorf("PC: %d, outside code range", v.pc)
		}
		operandA := v.code[v.pc]
		operandB := v.code[v.pc+1]
		operandC := v.code[v.pc+2]

		if operandA < 0 {
			operandA = v.mem[-operandA]
			if operandA < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
			if operandA >= memSize {
				return fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandA)
			}
		}

		if operandB < 0 {
			operandB = v.mem[-operandB]
			if operandB < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
			if operandB >= memSize {
				return fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operandB)
			}
		}
		if operandC < 0 {
			operandC = v.mem[-operandC]
			if operandC < 0 {
				return fmt.Errorf("PC: %d, double indirect not supported", v.pc)
			}
		}

		//fmt.Printf("PC: %7s    SUBLEQ %s, %s, %s\n", v.addr2symbol(v.pc, true), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC, true))
		//fmt.Printf("                      %d - %d = ", v.mem[operandB], v.mem[operandA])

		if operandB == hltLoc {
			v.hltVal = -v.mem[operandA]
			hlt = true
			break
		} else {
			v.mem[operandB] -= v.mem[operandA]
		}
		//fmt.Printf("%d\n", v.mem[operandB])

		if v.mem[operandB] <= 0 {
			v.pc = operandC
		} else {
			v.pc += 3
		}
	}
	return nil
}

func (v *SUBLEQ) LoadRoutine(code []int64, data []int64, codeSymbols map[string]int64, dataSymbols map[string]int64) {
	copy(v.code[:], code)
	copy(v.mem[:], data)
	v.codeSize = int64(len(v.code))
	v.codeSymbols = codeSymbols
	v.dataSymbols = dataSymbols
}

func (v *SUBLEQ) addr2symbol(addr int64, onlyCode ...bool) string {
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
