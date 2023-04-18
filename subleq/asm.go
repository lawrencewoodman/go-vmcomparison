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
var reLabel = regexp.MustCompile(`^\s*([a-zA-Z][0-9a-zA-z]*):`)
var reInstr2 = regexp.MustCompile(`^\s*([0-9a-zA-Z]+)\s+([0-9a-zA-Z]+)`)
var reInstr3 = regexp.MustCompile(`^\s*([0-9a-zA-Z]+)\s+([0-9a-zA-Z]+)\s+([0-9a-zA-Z]+)`)
var reLiteral = regexp.MustCompile(`^\s*([\-]?[0-9]+).*`)

// Build symbol table
func pass1(srcLines []string) map[string]int {
	pos := 0
	symbols := make(map[string]int, 0)
	for _, line := range srcLines {
		//fmt.Printf("%s\n", line)
		// If there is a label
		if reLabel.MatchString(line) {
			label := reLabel.FindStringSubmatch(line)[1]
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			symbols[label] = pos
			line = line[matchIndices[1]:]
		}

		// If there is an instruction
		if reInstr3.MatchString(line) || reInstr2.MatchString(line) {
			pos += 3
		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			pos++
		}
	}
	return symbols
}

func pass2(srcLines []string, symbols map[string]int) []int {
	pos := 0
	code := make([]int, 0)
	for _, line := range srcLines {
		// If there is a label
		if reLabel.MatchString(line) {
			// Remove from line
			matchIndices := reLabel.FindStringSubmatchIndex(line)
			line = line[matchIndices[1]:]
		}
		// If there is a 3 operand instruction
		if reInstr3.MatchString(line) {
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
				panic(fmt.Sprintf("remaining line: %s", line))
			}
			code = append(code, asmInstr(symbols, operandA, operandB, operandC)...)
			pos += 3

		} else if reLiteral.MatchString(line) {
			// If there is a literal value
			lit := reLiteral.FindStringSubmatch(line)[1]
			i64, err := strconv.ParseInt(lit, 10, 64)
			if err != nil {
				panic(err)
			}
			code = append(code, int(i64))
			pos++
		}
	}
	return code
}

func resolveOperand(symbols map[string]int, operand string) int {
	// If operand is a literal value
	if reLiteral.MatchString(operand) {
		i64, err := strconv.ParseUint(operand, 10, 64)
		if err != nil {
			panic(err)
		}
		//		fmt.Printf("lit: %d\n", ui64)
		return int(i64)
		// If operand is an indexed address
	}
	v, ok := symbols[operand]
	if !ok {
		panic(fmt.Sprintf("unknown operand: %s", operand))
	}
	return v
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

func asm(filename string) ([]int, error) {
	srcLines, err := readFile(filename)
	if err != nil {
		return []int{}, err
	}
	//fmt.Printf("before pass1\n")
	symbols := pass1(srcLines)
	/*
		fmt.Printf("Symbols\n=======\n")
		for k, v := range symbols {
			fmt.Printf("%s: %d\n", k, v)
		}
	*/
	code := pass2(srcLines, symbols)
	// fmt.Printf("%v\n", code)
	return code, nil
}
