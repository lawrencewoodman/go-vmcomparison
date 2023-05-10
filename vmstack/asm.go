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

var instructions = map[string]uint{
	"HLT":     0 << 24,
	"FETCH":   1 << 24,
	"STORE":   2 << 24,
	"ADD":     3 << 24,
	"AND":     5 << 24,
	"JNZ":     7 << 24,
	"DJNZ":    11 << 24,
	"JMP":     12 << 24,
	"SHL":     13 << 24,
	"LIT":     15 << 24,
	"DROP":    18 << 24,
	"SWAP":    19 << 24,
	"FETCHBI": 20 << 24,
	"ADDBI":   24 << 24,
	"FETCHI":  27 << 24,
	"JSR":     28 << 24,
	"RET":     29 << 24,
	"DUP":     30 << 24,
	"OR":      31 << 24,
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
var reLiteral = regexp.MustCompile(`^\s*([0-9]+)\s*`)
var reSymbol = regexp.MustCompile(`^\s*!([a-zA-Z][0-9a-zA-Z]*).*`)
var reComment = regexp.MustCompile(`^\s*(;.*)$`)

// Build symbol table
func pass1(srcLines []string) map[string]uint {
	var pos uint = 0
	symbols := make(map[string]uint, 0)
	for _, line := range srcLines {
		//		fmt.Printf("pos: %2d, line: %s\n", pos, line)
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

func pass2(srcLines []string, symbols map[string]uint) []uint {
	code := make([]uint, 0)
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
			ui64, err := strconv.ParseUint(lit, 10, 64)
			if err != nil {
				panic(err)
			}
			code = append(code, uint(ui64))
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

func resolveOperand(symbols map[string]uint, operand string) uint {
	if operand == "" {
		return 0
	}
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		ui64, err := strconv.ParseUint(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		//		fmt.Printf("lit: %d\n", ui64)
		return uint(ui64)
		// If operand is an indexed address
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func asmInstr(symbols map[string]uint, instr string, operand string) uint {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}
	//	fmt.Printf("asmInstr - opcode: %d, operand: %s\n", opcode, operand)

	code := opcode + resolveOperand(symbols, operand)
	//	fmt.Printf("code: %v\n", code)
	return code
}

func asm(filename string) ([]uint, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []uint{}, err
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
