/*
 * A simple assembler for vmstack
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vmstack

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var instructions = map[string]int64{
	"HLT":     0 << 24,
	"FETCH":   1 << 24,
	"STORE":   2 << 24,
	"ADD":     3 << 24,
	"SUB":     4 << 24,
	"AND":     5 << 24,
	"INC":     6 << 24,
	"JNZ":     7 << 24,
	"DJNZ":    8 << 24,
	"JMP":     9 << 24,
	"SHL":     10 << 24,
	"LIT":     11 << 24,
	"DROP":    12 << 24,
	"SWAP":    13 << 24,
	"FETCHBI": 14 << 24,
	"ADDBI":   15 << 24,
	"FETCHI":  16 << 24,
	"JSR":     17 << 24,
	"RET":     18 << 24,
	"DUP":     19 << 24,
	"OR":      20 << 24,
	"JZ":      21 << 24,
	"JGT":     22 << 24,
	"ROT":     23 << 24,
	"OVER":    24 << 24,
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
var reComment = regexp.MustCompile(`^\s*(;.*)$`)

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
			pos++
			continue
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
			code = append(code, asmInstr(symbols, instr, operand))
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			i64, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				panic(err)
			}

			code = append(code, i64)
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

func resolveOperand(symbols map[string]int64, operand string) int64 {
	if operand == "" {
		return 0
	}
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		i64, err := strconv.ParseInt(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return i64
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func asmInstr(symbols map[string]int64, instr string, operand string) int64 {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	code := opcode + resolveOperand(symbols, operand)
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
