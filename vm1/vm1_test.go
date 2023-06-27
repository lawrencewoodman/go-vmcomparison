package vm1

import (
	"fmt"
	"path/filepath"
	"testing"
)

var VMtests = []struct {
	filename string
	want     map[int64]int64 // [memloc]value
}{
	{"add12_v1.asm", map[int64]int64{12: 4}},
	{"and_v1.asm", map[int64]int64{21: 4499}},
	{"tad_v1.asm", map[int64]int64{20: 32}},
	{"isz_v1.asm", map[int64]int64{30: 9, 32: 24}},
	{"loopuntil_v1.asm", map[int64]int64{16: 5000}},
	{"loopuntil_v2.asm", map[int64]int64{12: 5000}},
	{"subleq_v1.asm", map[int64]int64{87: 5000}},
	{"switch_v1.asm", map[int64]int64{94: 2255}},
	{"switch_v2.asm", map[int64]int64{97: 2255}},
	{"jsr_v1.asm", map[int64]int64{6: 50}},
}

func TestRun(t *testing.T) {
	for _, test := range VMtests {
		t.Run(test.filename, func(t *testing.T) {
			routine, err := asm(filepath.Join("fixtures", test.filename))
			if err != nil {
				t.Fatalf("asm() err: %v", err)
			}
			v := New()
			v.LoadRoutine(routine)
			_, err = v.Run()
			if err != nil {
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
	for _, test := range VMtests {
		routine, err := asm(filepath.Join("fixtures", test.filename))
		if err != nil {
			b.Fatalf("asm() err: %v", err)
		}
		b.Run(test.filename, func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				v := New()
				v.LoadRoutine(routine)

				b.StartTimer()
				_, err = v.Run()
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
