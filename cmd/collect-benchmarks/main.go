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
)

type stat struct {
	pkg  string
	name string
	ns   uint
}

func usage(errMsg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
	fmt.Fprintf(os.Stderr, "Usage: %s filename\n", os.Args[0])
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
var reNameStub = regexp.MustCompile(`^([^_]*).*$`)

func parse(lines []string) []stat {
	var pkg string
	var stats []stat

	for _, line := range lines {
		// If there is a pkg: line
		if rePkg.MatchString(line) {
			pkg = rePkg.FindStringSubmatch(line)[1]
		} else if reBenchmark.MatchString(line) {
			testName := reBenchmark.FindStringSubmatch(line)[1]
			ns := reBenchmark.FindStringSubmatch(line)[2]
			ui64, err := strconv.ParseUint(ns, 10, 64)
			if err != nil {
				panic(err)
			}
			benchNS := uint(ui64)
			stats = append(stats, stat{pkg, testName, benchNS})
		}
	}
	return stats
}

func printCSV(stats []stat) {
	for _, s := range stats {
		fmt.Printf("%s,%s,%d\n", s.pkg, s.name, s.ns)
	}
}

func groupSort(stats []stat) []stat {
	sort.SliceStable(stats, func(i, j int) bool {
		nameI := reNameStub.FindStringSubmatch(stats[i].name)[1]
		nameJ := reNameStub.FindStringSubmatch(stats[j].name)[1]
		return nameI < nameJ
	})
	return stats
}

func main() {
	if len(os.Args) < 2 {
		usage("no filename")
		os.Exit(1)
	}
	filename := os.Args[1]
	lines, err := readFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
	stats := parse(lines)
	stats = groupSort(stats)
	printCSV(stats)
}
