/*
 * SUBLEQ Virtual Machine
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package subleq

import "fmt"

// TODO: Make this configurable
const memSize = 32000

// Location in memory of hltVal
// If this is used as a destination location then a HLT is executed
const hltLoc = 1000

type SUBLEQ struct {
	mem    [memSize]int // Memory
	pc     int          // Program Counter
	hltVal int          // A value returned by HLT
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

func (v *SUBLEQ) LoadRoutine(routine []int) {
	copy(v.mem[:], routine)
}

// fetch gets the next instruction from memory
// Returns: A, B, C
// NOTE: this routine doesn't use mask32
func (v *SUBLEQ) fetch() (int, int, int, error) {
	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", v.pc)
	}
	operandA := v.mem[v.pc]
	operandB := v.mem[v.pc+1]
	operandC := v.mem[v.pc+2]

	if operandA >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandA)
	}
	if operandB >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandB)
	}
	if operandC >= memSize {
		return 0, 0, 0, fmt.Errorf("outside memory range: %d", operandC)
	}

	return operandA, operandB, operandC, nil
}

// Maintain 32 bits
// Used rather than basing on int32 to maintain parity across other language platforms
func maintain32(n int) int {
	return int(int32(n))
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *SUBLEQ) execute(operandA int, operandB int, operandC int) bool {
	v.mem[operandB] = maintain32(v.mem[operandB] - v.mem[operandA])
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
