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
	"strconv"
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
var reLiteral = regexp.MustCompile(`^\s*(\-?[0-9]+).*`)
var reDirective = regexp.MustCompile(`^\.([a-zA-Z]+)$`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)

// pass1 returns program and data symbol tables
func pass1(srcLines []string) (map[string]int64, map[string]int64) {
	symbolType := "c"
	var memPos int64 = 0
	var progPos int64 = 0
	progSymbols := make(map[string]int64, 0)
	memSymbols := make(map[string]int64, 0)
	for _, line := range srcLines {
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				symbolType = "d"
			} else {
				panic(fmt.Sprintf("unknown directive: .%s", directive))
			}
		}
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			if symbolType == "c" {
				progSymbols[label] = progPos
			} else {
				memSymbols[label] = memPos
			}
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr.MatchString(line) {
			progPos += 3
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			memPos++
		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			memPos++
		}
	}
	return progSymbols, memSymbols
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

// pass2 returns code and data
func pass2(srcLines []string, codeSymbols, dataSymbols map[string]int64) ([]int64, []*big.Int) {
	outputType := "c"
	code := make([]int64, 0)
	data := make([]*big.Int, 0)
	for _, line := range srcLines {
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				outputType = "d"
			} else {
				panic(fmt.Sprintf("unknown directive: .%s", directive))
			}
		}
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
			code = append(code, asmInstr(codeSymbols, dataSymbols, instr, addrMode, operandA, operandB)...)
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			num, ok := new(big.Int).SetString(lit, 10)
			if !ok {
				panic("can't create literal value")
			}
			if outputType == "d" {
				data = append(data, num)
			} else {
				panic("literal shouldn't appear here")
			}
		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			sym := reSymbol.FindStringSubmatch(line)[1]
			addr, err := resolveSymbol(codeSymbols, dataSymbols, sym)
			if err != nil {
				panic(err)
			}
			if outputType == "d" {
				data = append(data, big.NewInt(addr))
			} else {
				panic("symbol shouldn't appear here")
			}
		}
	}
	return code, data
}

func resolveSymbol(codeSymbols, dataSymbols map[string]int64, sym string) (int64, error) {
	v, ok := codeSymbols[sym]
	if !ok {
		v, ok = dataSymbols[sym]
		if !ok {
			return v, fmt.Errorf("unknown symbol: %s", sym)
		}
	}
	return v, nil

}

func resolveOperand(codeSymbols, dataSymbols map[string]int64, operand string) int64 {
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		addr, err := strconv.ParseInt(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return addr
	}
	addr, err := resolveSymbol(codeSymbols, dataSymbols, operand)
	if err != nil {
		panic(err)
	}
	return addr
}

func asmInstr(codeSymbols, dataSymbols map[string]int64, instr string, addrMode string, operandA string, operandB string) []int64 {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	opA := resolveOperand(codeSymbols, dataSymbols, operandA)
	opB := resolveOperand(codeSymbols, dataSymbols, operandB)

	switch addrMode {
	case "":
	case "I":
		opA = 0 - opA
	case "DI":
		opB = 0 - opB
	case "II":
		opA = 0 - opA
		opB = 0 - opB
	default:
		panic(fmt.Sprintf("unknown addressing mode: %s", addrMode))
	}

	code := []int64{opcode, opA, opB}
	return code
}

func printCode(code []int64) {
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

func asm(filename string) ([]int64, []*big.Int, map[string]int64, map[string]int64, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []int64{}, []*big.Int{}, map[string]int64{}, map[string]int64{}, err
	}
	codeSymbols, dataSymbols := pass1(srcLines)
	code, data := pass2(srcLines, codeSymbols, dataSymbols)
	//printSymbols(symbols)
	//printCode(code)

	return code, data, codeSymbols, dataSymbols, nil
}
