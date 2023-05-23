/*
 * A simple assembler for v2 of this VM
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package vm2

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var instructions = map[string]uint{
	"HLT":  0 << 24,
	"MOV":  1 << 24,
	"JSR":  2 << 24,
	"ADD":  3 << 24,
	"DJNZ": 4 << 24,
	"JMP":  5 << 24,
	"LIT":  6 << 24,
	"AND":  7 << 24,
	"OR":   8 << 24,
	"SHL":  9 << 24,
	"JNZ":  10 << 24,
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
var reLiteral = regexp.MustCompile(`^\s*([0-9]+).*`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)

// Build symbol table
func pass1(srcLines []string) map[string]uint {
	var pos uint = 0
	symbols := make(map[string]uint, 0)
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
			// TODO: replace panic addition to errors
			panic(fmt.Sprintf("no operand found for instruction: %s", line))
		}
	}
	return operand, line
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
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		ui64, err := strconv.ParseUint(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return uint(ui64)
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func asmInstr(symbols map[string]uint, instr string, addrMode string, operandA string, operandB string) []uint {
	opcode, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	switch addrMode {
	case "":
	case "I":
		opcode |= 0x80 << 24
	case "DI":
		opcode |= 0x40 << 24
	default:
		panic(fmt.Sprintf("unknown addressing mode: %s", addrMode))
	}

	code := []uint{opcode + resolveOperand(symbols, operandA), resolveOperand(symbols, operandB)}
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
	//	fmt.Printf("%v\n", code)
	return code, nil
}
