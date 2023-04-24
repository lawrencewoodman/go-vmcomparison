package subleq

import (
	"fmt"
	"math"
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[int]int // [memloc]value
}{
	{"loopuntil_v1.asm", map[int]int{14: 5000}},
	{"add12_v1.asm", map[int]int{98: 4}},
	{"jsr_v1.asm", map[int]int{22: 50}},
}

func TestRun(t *testing.T) {
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			routine, symbols, err := asm(filepath.Join("fixtures", test.filename))
			if err != nil {
				t.Fatalf("asm() err: %v", err)
			}
			v := New()
			v.LoadRoutine(routine, symbols)
			if err := v.Run(); err != nil {
				t.Fatalf("Run() err: %v", err)
			}
			for memLoc, wantValue := range test.want {
				if v.mem[memLoc] != wantValue {
					t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
				}
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	for _, test := range tests {
		routine, symbols, err := asm(filepath.Join("fixtures", test.filename))
		if err != nil {
			b.Fatalf("asm() err: %v", err)
		}

		b.Run(test.filename, func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				v := New()
				v.LoadRoutine(routine, symbols)

				b.StartTimer()
				err := v.Run()
				b.StopTimer()

				if err != nil {
					b.Errorf("Run() err: %v", err)
				}
				for memLoc, wantValue := range test.want {
					if v.mem[memLoc] != wantValue {
						b.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
					}
				}
			}
		})
		fmt.Printf("Routine: %s size: %d\n", test.filename, len(routine))
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
				res += 1
			}
		}

		a += a
		b += b
	}
	return res
}
