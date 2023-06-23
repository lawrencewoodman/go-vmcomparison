/*
 * A simple assembler for SUBLEQ2 virtual machine using Big Numbers
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package bsubleq2

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
)

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
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-z]*):`)
var reInstr2 = regexp.MustCompile(`^\s*([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)`)
var reInstr3 = regexp.MustCompile(`^\s*([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)\s+([\[\]0-9a-zA-Z\-\+]+)`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)
var reExpr = regexp.MustCompile(`^\s*([0-9a-zA-z]+)([\-\+])([0-9a-zA-Z]+)`)
var reDirective = regexp.MustCompile(`^\.([a-zA-Z]+)$`)
var reIndirect = regexp.MustCompile(`^\s*(\[([0-9a-zA-z\-\+]+)\])`)

// Build symbol table
func pass1(srcLines []string) (map[string]int64, map[string]int64) {
	symbolType := "c"
	var memPos int64 = 0
	var progPos int64 = 0
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
				codeSymbols[label] = progPos
			} else {
				dataSymbols[label] = memPos
			}
			line = line[matchIndices[1]:]
		}
		if reInstr3.MatchString(line) || reInstr2.MatchString(line) {
			// If there is an instruction
			progPos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			memPos++
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			memPos++
		}

	}
	return codeSymbols, dataSymbols
}

func pass2(srcLines []string, codeSymbols, dataSymbols map[string]int64) ([]int64, []*big.Int) {
	outputType := "c"
	code := make([]int64, 0)
	data := make([]*big.Int, 0)
	codePos := 0
	lineNum := 0
	for _, line := range srcLines {
		lineNum++
		//fmt.Printf("%d: %s\n", lineNum, line)
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
		if reInstr3.MatchString(line) {
			// If there is a 3 operand instruction
			operandA := reInstr3.FindStringSubmatch(line)[1]
			operandB := reInstr3.FindStringSubmatch(line)[2]
			operandC := reInstr3.FindStringSubmatch(line)[3]
			matchIndices := reInstr3.FindStringSubmatchIndex(line)
			line = line[matchIndices[7]:]
			if len(line) > 0 {
				panic(fmt.Sprintf("remaining line: %s", line))
			}
			code = append(code, asmInstr(codeSymbols, dataSymbols, operandA, operandB, operandC)...)
			codePos += 3
		} else if reInstr2.MatchString(line) {
			// If there is a 2 operand instruction
			operandA := reInstr2.FindStringSubmatch(line)[1]
			operandB := reInstr2.FindStringSubmatch(line)[2]
			operandC := strconv.Itoa(codePos + 3)
			matchIndices := reInstr2.FindStringSubmatchIndex(line)
			line = line[matchIndices[5]:]
			if len(line) > 0 {
				panic(fmt.Sprintf("line number: %d, remaining line: %s", lineNum, line))
			}
			code = append(code, asmInstr(codeSymbols, dataSymbols, operandA, operandB, operandC)...)
			codePos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			a := reExpr.FindStringSubmatch(line)[1]
			op := reExpr.FindStringSubmatch(line)[2]
			b := reExpr.FindStringSubmatch(line)[3]
			if outputType == "d" {
				data = append(data, big.NewInt(resolveExpr(codeSymbols, dataSymbols, a, op, b)))
			} else {
				panic("expression shouldn't appear here")
			}
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
	if reIndirect.MatchString(operand) {
		s := reIndirect.FindStringSubmatch(operand)[2]
		return -resolveOperand(codeSymbols, dataSymbols, s)
	} else if reExpr.MatchString(operand) {
		// If operand is an expression
		a := reExpr.FindStringSubmatch(operand)[1]
		op := reExpr.FindStringSubmatch(operand)[2]
		b := reExpr.FindStringSubmatch(operand)[3]
		return resolveExpr(codeSymbols, dataSymbols, a, op, b)
	} else if reLiteral.MatchString(operand) {
		// If operand is a literal value
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

func resolveExpr(codeSymbols, dataSymbols map[string]int64, a, op, b string) int64 {
	aVal := resolveOperand(codeSymbols, dataSymbols, a)
	bVal := resolveOperand(codeSymbols, dataSymbols, b)
	switch op {
	case "+":
		return aVal + bVal
	case "-":
		return aVal - bVal
	default:
		panic(fmt.Sprintf("unknown op: %s", op))
	}
}

func asmInstr(codeSymbols, dataSymbols map[string]int64, operandA string, operandB string, operandC string) []int64 {
	code := []int64{
		resolveOperand(codeSymbols, dataSymbols, operandA),
		resolveOperand(codeSymbols, dataSymbols, operandB),
		resolveOperand(codeSymbols, dataSymbols, operandC),
	}
	return code
}

func printSymbols(symbols map[string]*big.Int) {
	labels := make([]string, 0, len(symbols))
	for l := range symbols {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	fmt.Printf("Symbols\n=======\n")
	for _, l := range labels {
		fmt.Printf("%s: %d\n", l, symbols[l])
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
	//fmt.Printf("%v\n", code)
	return code, data, codeSymbols, dataSymbols, nil
}
