/*
 * SUBLEQ2 Virtual Machine using Big Numbers
 *
 * SUBLEQ2 understands negative numbers passed as operands to be
 * indirect addresses.
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package bsubleq2

import (
	"fmt"
	"math/big"
)

// TODO: Make this configurable
const memSize = 32000

// Location in memory of hltVal
// If this is used as a destination location then a HLT is executed
const hltLoc = 1000

type SUBLEQ struct {
	mem     [memSize]*big.Int   // Memory
	pc      int64               // Program Counter
	hltVal  *big.Int            // A value returned by HLT
	symbols map[string]*big.Int // The symbols table from the assembler - added because of difficulty debugging
}

func New() *SUBLEQ {
	var mem [memSize]*big.Int
	for i := 0; i < memSize; i++ {
		mem[i] = big.NewInt(0)

	}
	return &SUBLEQ{mem: mem}
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

func (v *SUBLEQ) LoadRoutine(routine []*big.Int, symbols map[string]*big.Int) {
	// Need to copy the individual data points of the routine because they are pointers
	for i, d := range routine {
		v.mem[i].Set(d)
	}
	v.symbols = symbols
}

// getOperand returns the operand as supplied unless it is negative in which
// case it returns the value at the location in memory pointed to by the
// operand
func (v *SUBLEQ) getOperand(operand *big.Int) (int64, error) {
	if !operand.IsInt64() {
		return 0, fmt.Errorf("PC: %d, outside memory range", v.pc)
	}
	addr := operand.Int64()
	if addr < 0 {
		addr = 0 - addr
		if addr >= memSize {
			return 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		iaddr := v.mem[addr]
		if !iaddr.IsInt64() {
			return 0, fmt.Errorf("PC: %d, outside memory range", v.pc)
		}
		addr = iaddr.Int64()
		if addr < 0 {
			return 0, fmt.Errorf("PC: %d, double indirect not supported", v.pc)
		}
	}
	if addr >= memSize {
		return 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
	}

	return addr, nil
}

// fetch gets the next instruction from memory
// Returns: A, B, C, error
func (v *SUBLEQ) fetch() (int64, int64, int64, error) {
	var a, b, c int64
	var err error

	if v.pc+2 >= memSize {
		return 0, 0, 0, fmt.Errorf("PC: %d, outside memory range: %d", v.pc, v.pc)
	}
	operandA := v.mem[v.pc]
	operandB := v.mem[v.pc+1]
	operandC := v.mem[v.pc+2]

	a, err = v.getOperand(operandA)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err = v.getOperand(operandB)
	if err != nil {
		return 0, 0, 0, err
	}
	c, err = v.getOperand(operandC)
	if err != nil {
		return 0, 0, 0, err
	}

	return a, b, c, nil
}

func (v *SUBLEQ) addr2symbol(addr *big.Int) string {
	for s, a := range v.symbols {
		if a.Cmp(addr) == 0 {
			return s
		}
	}
	return fmt.Sprintf("%d", addr)
}

// execute executes the supplied instruction
// Returns: hlt, error
func (v *SUBLEQ) execute(operandA int64, operandB int64, operandC int64) bool {
	//fmt.Printf("PC: %7s    SUBLEQ %s, %s, %s\n", v.addr2symbol(big.NewInt(v.pc)), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC))
	//fmt.Printf("                      %d - %d = ", v.mem[operandB.Int64()], v.mem[operandA.Int64()])
	v.mem[operandB].Sub(v.mem[operandB], v.mem[operandA])
	//fmt.Printf("%d\n", v.mem[operandB.Int64()])
	if operandB == hltLoc {
		v.hltVal = v.mem[operandB]
		return true
	}

	if v.mem[operandB].Sign() <= 0 {
		v.pc = operandC
	} else {
		v.pc += 3
	}
	return false
}
