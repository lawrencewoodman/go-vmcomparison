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
	mem     [memSize]int   // Memory
	pc      int            // Program Counter
	hltVal  int            // A value returned by HLT
	symbols map[string]int // The symbols table from the assembler - added because of difficulty debugging
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

func (v *SUBLEQ) LoadRoutine(routine []int, symbols map[string]int) {
	copy(v.mem[:], routine)
	v.symbols = symbols
}

// fetch gets the next instruction from memory
// Returns: A, B, C
// NOTE: this routine doesn't use mask32
// NOTE: this causes a problem with debugging line as it won't be
// NOTE: obvious which operands are using indirect addressing
func (v *SUBLEQ) fetch() (int, int, int, error) {
	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", v.pc)
	}
	operandA := v.mem[v.pc]
	operandB := v.mem[v.pc+1]
	operandC := v.mem[v.pc+2]

	if operandA < 0 {
		operandA = 0 - operandA
		if operandA >= memSize {
			return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandA)
		}
		operandA = v.mem[operandA]
	} else if operandA >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandA)
	}

	if operandB < 0 {
		operandB = 0 - operandB
		if operandB >= memSize {
			return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandB)
		}
		operandB = v.mem[operandB]
	} else if operandB >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandB)
	}

	if operandC < 0 {
		operandC = 0 - operandC
		if operandC >= memSize {
			return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandC)
		}
		operandC = v.mem[operandC]
	} else if operandC >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandC)
	}

	return operandA, operandB, operandC, nil
}

// Maintain 32 bits
// Used rather than basing on int32 to maintain parity across other language platforms
func maintain32(n int) int {
	return int(int32(n))
}

func (v *SUBLEQ) addr2symbol(addr int) string {
	for k, v := range v.symbols {
		if v == addr {
			return k
		}
	}
	return fmt.Sprintf("%d", addr)
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *SUBLEQ) execute(operandA int, operandB int, operandC int) bool {
	//	fmt.Printf("PC: %7s    SUBLEQ %7s, %7s, %7s\n", v.addr2symbol(v.pc), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC))
	//	fmt.Printf("           %d (%b) - %d (%b) = ", v.mem[operandB], v.mem[operandB], v.mem[operandA], v.mem[operandA])
	v.mem[operandB] = maintain32(v.mem[operandB] - v.mem[operandA])
	//	fmt.Printf("%d (%b)\n", v.mem[operandB], v.mem[operandB])
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
