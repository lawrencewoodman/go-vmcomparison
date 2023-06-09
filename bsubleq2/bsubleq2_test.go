package bsubleq2

import (
	"fmt"
	"math"
	"math/big"
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[int]int64 // [memloc]value
}{
	{"loopuntil_v1.asm", map[int]int64{2: 5000}},
	{"add12_v1.asm", map[int]int64{2: 4}},
	{"and_v1.asm", map[int]int64{9: 4499}},
	{"and_v2.asm", map[int]int64{8: 4499}},
	{"isz_v1.asm", map[int]int64{9: 9, 10: 24}},
	{"jsr_v1.asm", map[int]int64{4: 50}},
	{"tad_v1.asm", map[int]int64{8: 32}},
	{"subleq_v1.asm", map[int]int64{27: 5000}},
	{"subleq_v2.asm", map[int]int64{25: 5000}},
	// TODO: reimplement switch_v1?
	// {"switch_v1.asm", map[int]int64{3: 2255}},
	{"switch_v2.asm", map[int]int64{0: 2255}},
	{"switch_v3.asm", map[int]int64{0: 2255}},
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
			if err := v.Run(); err != nil {
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
				err := v.Run()
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

func TestAND(t *testing.T) {
	for a := 0; a <= math.MaxUint8; a++ {
		for b := 0; b <= math.MaxUint8; b++ {
			want := a & b
			got := op_AND(a, b, 8)
			if want != got {
				t.Errorf("op_AND  a: %8b, b: %8b, got: %8b, want: %8b", a, b, got, want)
			}
		}
	}
}

func TestOR(t *testing.T) {
	for a := 0; a <= math.MaxUint8; a++ {
		for b := 0; b <= math.MaxUint8; b++ {
			want := a | b
			got := op_OR(a, b, 8)
			if want != got {
				t.Errorf("op_OR  a: %8b, b: %8b, got: %8b, want: %8b", a, b, got, want)
			}
		}
	}
}

// This is just here to test logic of routine suitable for running under SUBLEQ
func op_AND(a, b, n int) int {
	hbitval := int(math.Pow(2, float64(n-1)))
	res := 0
	for x := 0; x < n; x++ {
		m := 0
		res += res
		if a >= hbitval {
			m++
			a -= hbitval
		}
		if b >= hbitval {
			b -= hbitval
			if m == 1 {
				res++
			}
		}

		a += a
		b += b
	}
	return res
}

// This is just here to test logic of routine suitable for running under SUBLEQ
func op_OR(a, b, n int) int {
	hbitval := int(math.Pow(2, float64(n-1)))
	res := 0
	for x := 0; x < n; x++ {
		m := 0
		res += res
		if a >= hbitval {
			m++
			a -= hbitval
		}
		if b >= hbitval {
			m++
			b -= hbitval
		}

		if m > 0 {
			res++
		}

		a += a
		b += b
	}
	return res
}
