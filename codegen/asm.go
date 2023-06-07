/*
 * A simple assembler for this VM
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package codegen

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

var instructions = map[string]uint{
	"HLT":   0 << 24,
	"LDA":   1 << 24,
	"STA":   2 << 24,
	"ADD":   3 << 24,
	"SUB":   4 << 24,
	"AND":   5 << 24,
	"INC":   6 << 24,
	"JNZ":   7 << 24,
	"DSZ":   8 << 24,
	"JMP":   9 << 24,
	"SHL":   10 << 24,
	"LDX":   11 << 24,
	"LDY":   12 << 24,
	"DYJNZ": 13 << 24,
	"JSR":   14 << 24,
	"RET":   15 << 24,
	"TAY":   16 << 24,
	"STY":   17 << 24,
	"OR":    18 << 24,
	"JEQ":   19 << 24,
	"JGT":   20 << 24,
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
var reLiteral = regexp.MustCompile(`^\s*([\-]?)([0-9]+).*`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)
var reIndexOperand = regexp.MustCompile(`^\s*([0-9a-zA-Z]+),([0-9a-zA-z]+).*`)
var reDirective = regexp.MustCompile(`^\.([a-zA-Z]+)$`)
var reFilename = regexp.MustCompile(`^([0-9a-zA-z_]+).*`)

// Build symbol table
func pass1(srcLines []string) (map[string]uint, map[string]uint) {
	symbolType := "p"
	var memPos uint = 0
	var progPos uint = 0
	progSymbols := make(map[string]uint, 0)
	memSymbols := make(map[string]uint, 0)
	for _, line := range srcLines {
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				symbolType = "m"
			} else {
				panic(fmt.Sprintf("unknown directive: .%s", directive))
			}
		}
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			if symbolType == "p" {
				progSymbols[label] = progPos
			} else {
				memSymbols[label] = memPos
			}
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr.MatchString(line) {
			progPos++
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

func pass2(srcLines []string, progSymbols, memSymbols map[string]uint) string {
	code := "\tprogram := []func(v *CGVM){\n"
	for lineNum, line := range srcLines {
		// If there is a directive
		if reDirective.MatchString(line) {
			directive := reDirective.FindStringSubmatch(line)[1]
			if directive == "data" {
				code += "\t}"
				code += "\n"
				code += "\tmemory := []uint{\n"
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
			code += asmInstr(progSymbols, memSymbols, instr, addrMode, operand)

		} else if reSymbol.MatchString(line) {
			// If there is a symbol
			sym := reSymbol.FindStringSubmatch(line)[1]
			_, ok := memSymbols[sym]
			if !ok {
				panic(fmt.Sprintf("%d: unknown symbol: %s", lineNum, sym))
			}

			code += fmt.Sprintf("\t\tm_%s,\n", sym)

		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			sign := reLiteral.FindStringSubmatch(line)[1]
			lit := reLiteral.FindStringSubmatch(line)[2]
			if sign == "-" {
				ui64, err := strconv.ParseUint(lit, 10, 64)
				if err != nil {
					panic(err)
				}
				// Roll the number around to represent a negative number
				ui64 = math.MaxUint64 - (ui64 - 1)
				lit = strconv.FormatUint(ui64, 10)
			}

			code += fmt.Sprintf("\t\t%s,\n", lit)
		}
	}
	code += "\t}\n"
	return code
}

func resolveOperand(progSymbols, memSymbols map[string]uint, addrMode, operand string) string {
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		return operand
		// If operand is an indexed address
	} else if reIndexOperand.MatchString(operand) {
		if addrMode != "II" {
			panic("operand doesn't match addressing mode")
		}
		base := reIndexOperand.FindStringSubmatch(operand)[1]
		base = resolveOperand(progSymbols, memSymbols, "", base)
		index := reIndexOperand.FindStringSubmatch(operand)[2]
		index = resolveOperand(progSymbols, memSymbols, "", index)
		return fmt.Sprintf("calcBaseIndexAddr(v, %s, %s)", base, index)
	}

	if _, ok := progSymbols[operand]; ok {
		return fmt.Sprintf("p_%s", operand)
	}

	if _, ok := memSymbols[operand]; ok {
		return fmt.Sprintf("m_%s", operand)
	}
	panic(fmt.Sprintf("unknown operand: %s", operand))
}

func asmInstr(progSymbols, memSymbols map[string]uint, instr string, addrMode string, operand string) string {
	// TODO: don't need map for instructions as opcode value isn't needed
	_, ok := instructions[instr]
	if !ok {
		panic(fmt.Sprintf("unknown instruction: %s", instr))
	}

	return fmt.Sprintf("\t\tfunc(v *CGVM) { op_%s(v, %s) },\n", instr, resolveOperand(progSymbols, memSymbols, addrMode, operand))
}

func createConsts(progSymbols, memSymbols map[string]uint) string {
	code := "\tconst (\n"
	pkeys := make([]string, 0, len(progSymbols))
	mkeys := make([]string, 0, len(memSymbols))

	for k := range progSymbols {
		pkeys = append(pkeys, k)
	}
	for k := range memSymbols {
		mkeys = append(mkeys, k)
	}

	sort.Strings(pkeys)
	sort.Strings(mkeys)

	for _, k := range pkeys {
		code += fmt.Sprintf("\t\tp_%s = %d\n", k, progSymbols[k])
	}

	code += "\t)\n"
	if len(memSymbols) > 0 {
		code += "\tconst (\n"
		for _, k := range mkeys {
			code += fmt.Sprintf("\t\tm_%s = %d\n", k, memSymbols[k])
		}
		code += "\t)\n"
	}
	return code
}

func makeWantMapStr(want map[uint]uint) string {
	str := "map[uint]uint{"
	for k, v := range want {
		str += fmt.Sprintf("%d: %d,", k, v)
	}
	str += "}"
	return str
}

func asm(filename string, want map[uint]uint) (string, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return "", err
	}

	progSymbols, memSymbols := pass1(srcLines)
	header := "// Generated test file by main_test.go\n\n"
	header += "package codegen\n\n"
	cmd_name := reFilename.FindStringSubmatch(filepath.Base(filename))[1]
	code := header
	code += fmt.Sprintf("func init%s() ([]uint, []func(*CGVM)) {\n", cmd_name)
	code += createConsts(progSymbols, memSymbols)
	code += pass2(srcLines, progSymbols, memSymbols)
	code += "\treturn memory, program\n"
	code += "}\n\n"
	code += "func init() {\n"
	code += fmt.Sprintf("\taddTest(\"%s\", init%s, %s)\n", cmd_name, cmd_name, makeWantMapStr(want))
	code += "}\n"
	return code, nil
}
