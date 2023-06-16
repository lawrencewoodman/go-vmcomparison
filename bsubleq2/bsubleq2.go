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
func (v *SUBLEQ) getOperand(operand *big.Int) (*big.Int, error) {
	addr := new(big.Int)
	if operand.Sign() < 0 {
		addr.Abs(operand)
		// if operand >= memSize
		if addr.Cmp(big.NewInt(memSize)) >= 0 {
			return big.NewInt(0), fmt.Errorf("PC: %d, outside memory range: %d", v.pc, addr)
		}
		addr = v.mem[addr.Int64()]
		if addr.Sign() < 0 {
			return big.NewInt(0), fmt.Errorf("PC: %d, double indirect not supported", v.pc)
		}
	} else {
		addr.Set(operand)
	}
	// if addr >= memSize
	if addr.Cmp(big.NewInt(memSize)) >= 0 {
		return big.NewInt(0), fmt.Errorf("PC: %d, outside memory range: %d", v.pc, operand)
	}

	return addr, nil
}

// fetch gets the next instruction from memory
// Returns: A, B, C, error
func (v *SUBLEQ) fetch() (*big.Int, *big.Int, *big.Int, error) {
	var err error

	if v.pc+2 >= memSize {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0),
			fmt.Errorf("PC: %d, outside memory range: %d", v.pc, v.pc)
	}
	operandA := v.mem[v.pc]
	operandB := v.mem[v.pc+1]
	operandC := v.mem[v.pc+2]

	operandA, err = v.getOperand(operandA)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}
	operandB, err = v.getOperand(operandB)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}
	operandC, err = v.getOperand(operandC)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0), err
	}

	return operandA, operandB, operandC, nil
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
func (v *SUBLEQ) execute(operandA *big.Int, operandB *big.Int, operandC *big.Int) bool {
	//fmt.Printf("PC: %7s    SUBLEQ %s, %s, %s\n", v.addr2symbol(big.NewInt(v.pc)), v.addr2symbol(operandA), v.addr2symbol(operandB), v.addr2symbol(operandC))
	//fmt.Printf("                      %d - %d = ", v.mem[operandB.Int64()], v.mem[operandA.Int64()])
	b := operandB.Int64()
	v.mem[b].Sub(v.mem[b], v.mem[operandA.Int64()])
	//fmt.Printf("%d\n", v.mem[operandB.Int64()])
	if b == hltLoc {
		v.hltVal = v.mem[b]
		return true
	}

	if v.mem[b].Sign() <= 0 {
		v.pc = operandC.Int64()
	} else {
		v.pc += 3
	}
	return false
}
