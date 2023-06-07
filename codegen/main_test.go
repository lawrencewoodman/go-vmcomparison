package codegen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	filesWritten, err := createTestFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if filesWritten > 0 {
		fmt.Fprintf(os.Stderr, "%d test files generated, rerun test\n", filesWritten)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

type TestFile struct {
	testFilename string
	asmFilename  string
	want         map[uint]uint
}

// Create test files if they don't exist or the source has changed
func createTestFiles() (int, error) {
	expectedTestFiles := []TestFile{
		{"tad_v1_test.go", "tad_v1.asm", map[uint]uint{3: 32}},
		{"subleq_v1_test.go", "subleq_v1.asm", map[uint]uint{22: 5000}},
		{"loopuntil_v1_test.go", "loopuntil_v1.asm", map[uint]uint{0: 5000}},
	}
	filesWritten := 0
	for _, testFile := range expectedTestFiles {
		writeFile := false
		asmFilename := filepath.Join("fixtures", testFile.asmFilename)
		source, err := asm(asmFilename, testFile.want)
		if err != nil {
			return filesWritten, err
		}

		if _, err := os.Stat(testFile.testFilename); err == nil {
			currentSource, err := os.ReadFile(testFile.testFilename)
			if err != nil {
				return filesWritten, err
			}
			if !bytes.Equal([]byte(source), currentSource) {
				writeFile = true
			}
		} else {
			writeFile = true
		}
		if writeFile {
			err = os.WriteFile(testFile.testFilename, []byte(source), 0644)

			if err != nil {
				return filesWritten, err
			}
			filesWritten++
		}
	}

	return filesWritten, nil
}
