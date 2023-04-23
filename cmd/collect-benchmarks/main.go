/*
 * A utility to parse the output of 'go test -bench'
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type stat struct {
	pkg  string
	name string
	ns   int64
	size int64
}

func usage(errMsg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
	fmt.Fprintf(os.Stderr, "Usage: %s [csv|tables] filename\n", os.Args[0])
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
var rePkg = regexp.MustCompile(`^pkg: .*?\/go-vmcomparison\/(.*)$`)
var reBenchmark = regexp.MustCompile(`^Benchmark.*?\/(.*?)-[^ ]+\s+\d+\s+(\d+) ns\/op$`)
var reSize = regexp.MustCompile(`^Routine: ([^ ]+) size: (\d+)$`)
var reNameStub = regexp.MustCompile(`^([^_]*).*$`)

func parse(lines []string) []stat {
	var pkg string
	var size int64
	var stats []stat
	var err error

	for _, line := range lines {
		//		fmt.Printf("line: %s\n", line)
		// If there is a pkg: line
		if rePkg.MatchString(line) {
			pkg = rePkg.FindStringSubmatch(line)[1]
			size = 0
		} else if reSize.MatchString(line) {
			// if there is a size line
			testName := reSize.FindStringSubmatch(line)[1]
			sizeS := reSize.FindStringSubmatch(line)[2]
			size, err = strconv.ParseInt(sizeS, 10, 64)
			if err != nil {
				panic(err)
			}
			if testName != stats[len(stats)-1].name {
				panic(fmt.Sprintf("test names do not match: '%s' != '%s'", testName, stats[len(stats)-1].name))
			}
			stats[len(stats)-1].size = size
		} else if reBenchmark.MatchString(line) {
			testName := reBenchmark.FindStringSubmatch(line)[1]
			nsS := reBenchmark.FindStringSubmatch(line)[2]
			ns, err := strconv.ParseInt(nsS, 10, 64)
			if err != nil {
				panic(err)
			}
			stats = append(stats, stat{pkg, testName, ns, size})
		}
	}
	return stats
}

func getStubName(s string) string {
	return reNameStub.FindStringSubmatch(s)[1]
}

func printCSV(stats []stat) {
	for _, s := range stats {
		fmt.Printf("%s,%s,%d,%d\n", s.pkg, s.name, s.ns, s.size)
	}
}

func printTables(stats []stat) {
	currentStubName := ""
	for _, s := range stats {
		stubName := getStubName(s.name)
		if stubName != currentStubName {
			// Print a title
			fmt.Printf("\n%s\n%s\n", stubName, strings.Repeat("=", len(stubName)))
			currentStubName = stubName
		}
		if s.size > 0 {
			fmt.Printf("%-8s %-17s %9d ns  %3d words\n", s.pkg, s.name, s.ns, s.size)
		} else {
			fmt.Printf("%-8s %-17s %9d ns\n", s.pkg, s.name, s.ns)
		}
	}
}

func groupSort(stats []stat) []stat {
	sort.SliceStable(stats, func(i, j int) bool {
		nameI := getStubName(stats[i].name)
		nameJ := getStubName(stats[j].name)
		return nameI < nameJ
	})
	return stats
}

func main() {
	filename := ""
	command := "csv"

	switch len(os.Args) {
	case 1:
		usage("no filename")
		os.Exit(1)
	case 2:
		filename = os.Args[1]
	case 3:
		command = os.Args[1]
		if command != "csv" && command != "tables" {
			usage(fmt.Sprintf("incorrect command: %s", command))
			os.Exit(1)
		}
		filename = os.Args[2]
	default:
		usage("wrong number of arguments")
		os.Exit(1)
	}

	lines, err := readFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
	stats := parse(lines)
	stats = groupSort(stats)
	switch command {
	case "csv":
		printCSV(stats)
	case "tables":
		printTables(stats)
	}
}
