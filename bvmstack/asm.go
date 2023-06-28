/*
 * A simple assembler for vmstack using Big numbers
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package bvmstack

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

var instructions = map[string]int64{
	"HLT":   0 << 24,
	"FETCH": 1 << 24,
	"STORE": 2 << 24,
	"ADD":   3 << 24,
	"SUB":   4 << 24,
	"AND":   5 << 24,
	"INC":   6 << 24,
	"JNZ":   7 << 24,
	"DJNZ":  8 << 24,
	"JMP":   9 << 24,
	"SHL":   10 << 24,
	"LIT":   11 << 24,
	"DROP":  12 << 24,
	"SWAP":  13 << 24,
	"JSR":   17 << 24,
	"RET":   18 << 24,
	"DUP":   19 << 24,
	"OR":    20 << 24,
	"JZ":    21 << 24,
	"JGT":   22 << 24,
	"ROT":   23 << 24,
	"OVER":  24 << 24,
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
var reInstr = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]+)\s*`)
var reOperand = regexp.MustCompile(`^\s*([0-9a-zA-Z]+)\s*`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+)\s*`)
var reSymbol = regexp.MustCompile(`^\s*!([a-zA-Z][0-9a-zA-Z]*).*`)
var reDirective = regexp.MustCompile(`^\.([a-zA-Z]+)$`)
var reComment = regexp.MustCompile(`^\s*(;.*)$`)

// pass1 returns code and data symbol tables
func pass1(srcLines []string) (map[string]int64, map[string]int64) {
	symbolType := "c"
	var dataPos int64 = 0
	var codePos int64 = 0
	codeSymbols := make(map[string]int64, 0)
	dataSymbols := make(map[string]int64, 0)
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
				codeSymbols[label] = codePos
			} else {
				dataSymbols[label] = dataPos
			}
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr.MatchString(line) {
			codePos++
			continue
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			dataPos++
		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			dataPos++
		}
	}
	return codeSymbols, dataSymbols
}

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

			operand := ""
			if reOperand.MatchString(line) {
				operand = reOperand.FindStringSubmatch(line)[1]
				matchIndices := reOperand.FindStringSubmatchIndex(line)
				line = line[matchIndices[3]:]
			}

			if reComment.MatchString(line) {
				matchIndices := reComment.FindStringSubmatchIndex(line)
				line = line[matchIndices[3]:]
			}

			if len(line) > 0 {
				panic(fmt.Sprintf("remaining line: %s", line))
			}
			code = append(code, asmInstr(codeSymbols, dataSymbols, instr, operand))
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
	if operand == "" {
		return 0
	}
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

func asmInstr(codeSymbols, dataSymbols map[string]int64, instr string, operand string) int64 {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	code := opcode + resolveOperand(codeSymbols, dataSymbols, operand)
	return code
}

func printSymbols(codeSymbols, dataSymbols map[string]int64) {
	fmt.Printf("Symbols\n=======\n\ncode:\n")
	for k, v := range codeSymbols {
		fmt.Printf("  %s: %d\n", k, v)
	}
	fmt.Printf("\ndata:\n")
	for k, v := range dataSymbols {
		fmt.Printf("  %s: %d\n", k, v)
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
	//printSymbols(codeSymbols, dataSymbols)
	//fmt.Printf("%v\n", code)
	return code, data, codeSymbols, dataSymbols, nil
}
