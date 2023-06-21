package bvm2

import (
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[int64]int64 // [memloc]value
}{
	{"add12_v1.asm", map[int64]int64{2: 4}},
	{"and_v1.asm", map[int64]int64{4: 4499}},
	{"tad_v1.asm", map[int64]int64{4: 32}},
	{"isz_v1.asm", map[int64]int64{4: 9, 6: 24}},
	{"jsr_v1.asm", map[int64]int64{2: 50}},
	{"loopuntil_v1.asm", map[int64]int64{0: 5000}},
	// TODO: reinstate subleq_v1?
	//{"subleq_v1.asm", map[int64]int64{101: 5000}},
	{"subleq_v2.asm", map[int64]int64{27: 5000}},
	{"subleq_v3.asm", map[int64]int64{25: 5000}},
	//TODO: reinstate switch_v1
	//{"switch_v1.asm", map[int64]int64{44: 2255}},
	{"switch_v2.asm", map[int64]int64{0: 2255}},
}

func TestRun(t *testing.T) {
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			code, data, codeSymbols, dataSymbols, err := asm(filepath.Join("fixtures", test.filename))
			if err != nil {
				t.Fatalf("asm() err: %v", err)
			}
			v := New()
			v.LoadRoutine(code, data, codeSymbols, dataSymbols)
			_, err = v.Run()
			if err != nil {
				t.Fatalf("Run() err: %v", err)
			}
			for memLoc, wantValue := range test.want {
				if v.mem[memLoc].Cmp(big.NewInt(wantValue)) != 0 {
					t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
				}
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	for _, test := range tests {
		code, data, codeSymbols, dataSymbols, err := asm(filepath.Join("fixtures", test.filename))
		if err != nil {
			b.Fatalf("asm() err: %v", err)
		}
		b.Run(test.filename, func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				v := New()
				v.LoadRoutine(code, data, codeSymbols, dataSymbols)

				b.StartTimer()
				_, err = v.Run()
				b.StopTimer()

				if err != nil {
					b.Errorf("Run() err: %v", err)
				}
				for memLoc, wantValue := range test.want {
					if v.mem[memLoc].Cmp(big.NewInt(wantValue)) != 0 {
						b.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
					}
				}
			}
		})
		fmt.Printf("Routine: %s size: %d\n", test.filename, len(code)+len(data))
	}
}
