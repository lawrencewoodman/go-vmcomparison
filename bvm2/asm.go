/*
 * A simple assembler for v2 of this VM using Big Numbers
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package bvm2

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
)

var instructions = map[string]int64{
	"HLT":  0,
	"MOV":  1,
	"JSR":  2,
	"ADD":  3,
	"DJNZ": 4,
	"JMP":  5,
	"AND":  6,
	"OR":   7,
	"SHL":  8,
	"JNZ":  9,
	"SNE":  10,
	"SLE":  11,
	"SUB":  12,
	"JGT":  13,
}

func readFile(filename string) ([]string, error) {
	lines := make([]string, 0)

	f, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return []string{}, err
	}
	return lines, nil
}

// Regular expressions for parts of a line
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*):`)
var reInstr = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]+)\s+`)
var reAddrMode = regexp.MustCompile(`^\s*([diDI]{1,2})\s+`)
var reOperand = regexp.MustCompile(`^\s*([0-9a-zA-Z]+)`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)

// Build symbol table
func pass1(srcLines []string) map[string]*big.Int {
	var pos int64 = 0
	symbols := make(map[string]*big.Int, 0)
	for _, line := range srcLines {
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			symbols[label] = big.NewInt(pos)
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr.MatchString(line) {
			pos += 3
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			pos++
		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			pos++
		}
	}
	return symbols
}

// Returns: operand, restLine
func getOperand(line string, addrMode string) (string, string) {
	operand := ""
	if reOperand.MatchString(line) {
		operand = reOperand.FindStringSubmatch(line)[1]
		matchIndices := reOperand.FindStringSubmatchIndex(line)
		line = line[matchIndices[3]:]
	} else {
		if addrMode != "" {
			// operand mistaken for address mode
			operand = addrMode
		} else {
			// If operand missing then use '0'
			operand = "0"
		}
	}
	return operand, line
}

func pass2(srcLines []string, symbols map[string]*big.Int) []*big.Int {
	code := make([]*big.Int, 0)
	for _, line := range srcLines {
		// If there is a label
		if reLabel.MatchString(line) {
			// Remove from line
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			line = line[matchIndices[1]:]
		}
		// If there is an instruction
		if reInstr.MatchString(line) {
			instr := reInstr.FindStringSubmatch(line)[1]
			matchIndices := reInstr.FindStringSubmatchIndex(line)
			line = line[matchIndices[3]:]
			addrMode := ""

			// If there is a possible address mode
			if reAddrMode.MatchString(line) {
				addrMode = reAddrMode.FindStringSubmatch(line)[1]
				matchIndices := reAddrMode.FindStringSubmatchIndex(line)
				line = line[matchIndices[3]:]
			}
			// TODO: handle operand being mistaken for addrMode
			operandA, line := getOperand(line, "")
			operandB, line := getOperand(line, "")

			if len(line) > 0 {
				panic(fmt.Sprintf("remaining line: %s", line))
			}
			code = append(code, asmInstr(symbols, instr, addrMode, operandA, operandB)...)
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			num, ok := new(big.Int).SetString(lit, 10)
			if !ok {
				panic("can't create literal value")
			}
			code = append(code, num)
		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			sym := reSymbol.FindStringSubmatch(line)[1]
			v, ok := symbols[sym]
			if !ok {
				panic(fmt.Sprintf("unknown symbol: %s", sym))
			}
			code = append(code, v)
		}
	}
	return code
}

func resolveOperand(symbols map[string]*big.Int, operand string) *big.Int {
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		num, ok := new(big.Int).SetString(operand, 10)
		if !ok {
			panic("can't create literal value")
		}
		return num
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func asmInstr(symbols map[string]*big.Int, instr string, addrMode string, operandA string, operandB string) []*big.Int {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	opA := resolveOperand(symbols, operandA)
	opB := resolveOperand(symbols, operandB)

	switch addrMode {
	case "":
	case "I":

		opA = big.NewInt(0).Neg(opA)
	case "DI":
		opB = big.NewInt(0).Neg(opB)
	case "II":
		opA = big.NewInt(0).Neg(opA)
		opB = big.NewInt(0).Neg(opB)
	default:
		panic(fmt.Sprintf("unknown addressing mode: %s", addrMode))
	}

	code := []*big.Int{big.NewInt(opcode), opA, opB}
	return code
}

func printCode(code []*big.Int) {
	fmt.Printf("\nCODE\n====")
	for i, v := range code {
		if i%10 == 0 {
			fmt.Printf("\n%03d: ", i)
		}
		fmt.Printf("%4d ", v)
	}
	fmt.Printf("\n")
}

func printSymbols(symbols map[string]*big.Int) {
	fmt.Printf("Symbols\n=======\n")
	for k, v := range symbols {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("\n")
}

func asm(filename string) ([]*big.Int, map[string]*big.Int, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []*big.Int{}, map[string]*big.Int{}, err
	}
	symbols := pass1(srcLines)
	code := pass2(srcLines, symbols)
	//printSymbols(symbols)
	//printCode(code)

	return code, symbols, nil
}
