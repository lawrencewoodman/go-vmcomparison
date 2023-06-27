/*
 * A simple assembler for this VM
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm1

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var instructions = map[string]int64{
	"HLT":   0,
	"LDA":   1,
	"STA":   2,
	"ADD":   3,
	"SUB":   4,
	"AND":   5,
	"INC":   6,
	"JNZ":   7,
	"DSZ":   8,
	"JMP":   9,
	"SHL":   10,
	"LDX":   11,
	"LDY":   12,
	"DYJNZ": 13,
	"JSR":   14,
	"RET":   15,
	"TAY":   16,
	"STY":   17,
	"OR":    18,
	"JEQ":   19,
	"JGT":   20,
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
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z]+[0-9a-zA-Z]*):`)
var reInstr = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]+)\s+`)
var reAddrMode = regexp.MustCompile(`^\s*([iI]{1,2})\s+`)
var reOperand = regexp.MustCompile(`^\s*([0-9a-zA-Z,]+).*`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)
var reIndexOperand = regexp.MustCompile(`^\s*([0-9a-zA-Z]+),([0-9a-zA-z]+).*`)

// Build symbol table
func pass1(srcLines []string) map[string]int64 {
	var pos int64 = 0
	symbols := make(map[string]int64, 0)
	for _, line := range srcLines {
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			symbols[label] = pos
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr.MatchString(line) {
			pos += 2
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

func pass2(srcLines []string, symbols map[string]int64) []int64 {
	code := make([]int64, 0)
	for lineNum, line := range srcLines {
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
			operand := ""
			// If there is a possible address mode
			if reAddrMode.MatchString(line) {
				addrMode = reAddrMode.FindStringSubmatch(line)[1]
				matchIndices := reAddrMode.FindStringSubmatchIndex(line)
				line = line[matchIndices[3]:]
			}
			// If there is an operand
			if reOperand.MatchString(line) {
				operand = reOperand.FindStringSubmatch(line)[1]
				matchIndices := reOperand.FindStringSubmatchIndex(line)
				line = line[matchIndices[3]:]
			} else {
				if addrMode != "" {
					// operand mistaken for address mode
					operand = addrMode
				} else {
					// TODO: replace panic addition to errors
					panic(fmt.Sprintf("%d: no operand found for instruction: %s", lineNum, line))
				}
			}
			if len(line) > 0 {
				panic(fmt.Sprintf("%d: remaining line: %s", lineNum, line))
			}
			code = append(code, asmInstr(symbols, instr, addrMode, operand)...)

		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			sym := reSymbol.FindStringSubmatch(line)[1]
			v, ok := symbols[sym]
			if !ok {
				panic(fmt.Sprintf("%d: unknown symbol: %s", lineNum, sym))
			}

			code = append(code, v)

		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			i64, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				panic(err)
			}
			code = append(code, i64)
		}
	}
	return code
}

func resolveOperand(symbols map[string]int64, operand string) int64 {
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		i64, err := strconv.ParseInt(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return i64
		// If operand is an indexed address
	} else if reIndexOperand.MatchString(operand) {
		base := reIndexOperand.FindStringSubmatch(operand)[1]
		index := reIndexOperand.FindStringSubmatch(operand)[2]
		baseAddr := resolveOperand(symbols, base)
		indexAddr := resolveOperand(symbols, index)
		return (baseAddr << 12) + indexAddr
		// TODO: error if > 4095
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func asmInstr(symbols map[string]int64, instr string, addrMode string, operand string) []int64 {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	opA := resolveOperand(symbols, operand)

	if addrMode == "I" {
		opA = -opA
	}
	code := []int64{opcode, opA}
	return code
}

func asm(filename string) ([]int64, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []int64{}, err
	}
	symbols := pass1(srcLines)
	/*
		fmt.Printf("Symbols\n=======\n")
		for k, v := range symbols {
			fmt.Printf("%s: %d\n", k, v)
		}
	*/
	code := pass2(srcLines, symbols)
	//fmt.Printf("%v\n", code)
	return code, nil
}
