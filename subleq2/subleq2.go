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
	mem     [memSize]int64   // Memory
	pc      int64            // Program Counter
	hltVal  int64            // A value returned by HLT
	symbols map[string]int64 // The symbols table from the assembler - added because of difficulty debugging
}

func New() *SUBLEQ {
	return &SUBLEQ{}
}

func (v *SUBLEQ) Step() (bool, error) {
	operandA, operandB, operandC, err := v.fetch()
	if err != nil {
		return false, err
	}
	return v.execute(operandA, operandB, operandC), nil
}

func (v *SUBLEQ) Run() error {
	var err error
	hlt := false
	for !hlt {
		hlt, err = v.Step()
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *SUBLEQ) LoadRoutine(routine []int64, symbols map[string]int64) {
	copy(v.mem[:], routine)
	v.symbols = symbols
}

// getOperand returns the operand as supplied unless it is negative in which
// case it returns the value at the location in memory pointed to by the
// operand
func (v *SUBLEQ) getOperand(operand int64) (int64, error) {
	if operand < 0 {
		operand = 0 - operand
		if operand >= memSize {
			return 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operand)
		}
		operand = v.mem[operand]
		if operand < 0 {
			return 0, fmt.Errorf("PC: %d, double indirect not supported", v.pc)
		}
	}
	if operand >= memSize {
		return 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operand)
	}

	return operand, nil
}

// fetch gets the next instruction from memory
// Returns: A, B, C, error
func (v *SUBLEQ) fetch() (int64, int64, int64, error) {
	var err error

	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, v.pc)
	}
	operandA := v.mem[v.pc]
	operandB := v.mem[v.pc+1]
	operandC := v.mem[v.pc+2]

	operandA, err = v.getOperand(operandA)
	if err != nil {
		return 0, 0, 0, err
	}
	operandB, err = v.getOperand(operandB)
	if err != nil {
		return 0, 0, 0, err
	}
	operandC, err = v.getOperand(operandC)
	if err != nil {
		return 0, 0, 0, err
	}

	return operandA, operandB, operandC, nil
}

func (v *SUBLEQ) addr2symbol(addr int64) string {
	for k, v := range v.symbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *SUBLEQ) execute(operandA int64, operandB int64, operandC int64) bool {
	//fmt.Printf("PC: %7s    SUBLEQ %s, %s, %s\n", v.addr2symbol(v.pc), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC))
	//fmt.Printf("                      %d - %d = ", v.mem[operandB], v.mem[operandA])
	v.mem[operandB] = v.mem[operandB] - v.mem[operandA]
	//fmt.Printf("%d\n", v.mem[operandB])
	if operandB == hltLoc {
		v.hltVal = v.mem[operandB]
		return true
	}
	if v.mem[operandB] <= 0 {
		v.pc = operandC
	} else {
		v.pc = v.pc + 3
	}
	return false
}
