/*
 * A simple assembler for SUBLEQ virtual machine
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package subleq

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*):`)
var reInstr2 = regexp.MustCompile(`^\s*([0-9a-zA-Z\-\+]+)\s+([0-9a-zA-Z\-\+]+)`)
var reInstr3 = regexp.MustCompile(`^\s*([0-9a-zA-Z\-\+]+)\s+([0-9a-zA-Z\-\+]+)\s+([0-9a-zA-Z\-\+]+)`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)
var reSymbol = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-Z]*).*`)
var reExpr = regexp.MustCompile(`^\s*([0-9a-zA-z]+)([\-\+])([0-9a-zA-Z]+)`)

// Build symbol table
func pass1(srcLines []string) map[string]int {
	pos := 0
	symbols := make(map[string]int, 0)
	for _, line := range srcLines {
		//		fmt.Printf("%s\n", line)
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			symbols[label] = pos
			line = line[matchIndices[1]:]
		}
		if reInstr3.MatchString(line) || reInstr2.MatchString(line) {
			// If there is an instruction
			pos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			pos++
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

func pass2(srcLines []string, symbols map[string]int) []int {
	pos := 0
	lineNum := 0
	code := make([]int, 0)
	for _, line := range srcLines {
		lineNum++
		//		fmt.Printf("pass2 line: %s\n", line)
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
			code = append(code, asmInstr(symbols, operandA, operandB, operandC)...)
			pos += 3
		} else if reInstr2.MatchString(line) {
			// If there is a 2 operand instruction
			operandA := reInstr2.FindStringSubmatch(line)[1]
			operandB := reInstr2.FindStringSubmatch(line)[2]
			operandC := strconv.Itoa(pos + 3)
			matchIndices := reInstr2.FindStringSubmatchIndex(line)
			line = line[matchIndices[5]:]
			if len(line) > 0 {
				panic(fmt.Sprintf("line number: %d, remaining line: %s", lineNum, line))
			}
			code = append(code, asmInstr(symbols, operandA, operandB, operandC)...)
			pos += 3
		} else if reExpr.MatchString(line) {
			// If there is an expression
			a := reExpr.FindStringSubmatch(line)[1]
			op := reExpr.FindStringSubmatch(line)[2]
			b := reExpr.FindStringSubmatch(line)[3]
			code = append(code, resolveExpr(symbols, a, op, b))
			pos++
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			i64, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				panic(err)
			}
			code = append(code, int(i64))
			pos++
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

func resolveOperand(symbols map[string]int, operand string) int {
	if reExpr.MatchString(operand) {
		// If operand is an expression
		a := reExpr.FindStringSubmatch(operand)[1]
		op := reExpr.FindStringSubmatch(operand)[2]
		b := reExpr.FindStringSubmatch(operand)[3]
		return resolveExpr(symbols, a, op, b)
	} else if reLiteral.MatchString(operand) {
		// If operand is a literal value
		i64, err := strconv.ParseUint(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		return int(i64)
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
}

func resolveExpr(symbols map[string]int, a, op, b string) int {
	aVal := resolveOperand(symbols, a)
	bVal := resolveOperand(symbols, b)
	switch op {
	case "+":
		return aVal + bVal
	case "-":
		return aVal - bVal
	default:
		panic(fmt.Sprintf("unknown op: %s", op))
	}
}

func asmInstr(symbols map[string]int, operandA string, operandB string, operandC string) []int {
	code := []int{
		resolveOperand(symbols, operandA),
		resolveOperand(symbols, operandB),
		resolveOperand(symbols, operandC),
	}
	//	fmt.Printf("code: %v\n", code)
	return code
}

func printCode(code []int) {
	fmt.Printf("\nCODE\n====")
	for i, v := range code {
		if i%10 == 0 {
			fmt.Printf("\n%03d: ", i)
		}
		fmt.Printf("%4d ", v)
	}
	fmt.Printf("\n")
}

func printSymbols(symbols map[string]int) {
	fmt.Printf("Symbols\n=======\n")
	for k, v := range symbols {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("\n")
}

func asm(filename string) ([]int, map[string]int, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []int{}, map[string]int{}, err
	}
	symbols := pass1(srcLines)

	code := pass2(srcLines, symbols)
	//	printSymbols(symbols)
	//	printCode(code)

	return code, symbols, nil
}
