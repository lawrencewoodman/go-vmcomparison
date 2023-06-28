package vmstack

import (
	"fmt"
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[int64]int64 // [memloc]value
}{
	{"add12_v1.asm", map[int64]int64{7: 4}},
	{"and_v1.asm", map[int64]int64{11: 4499}},
	{"tad_v1.asm", map[int64]int64{11: 32}},
	{"tad_v2.asm", map[int64]int64{10: 32}},
	{"isz_v1.asm", map[int64]int64{20: 9, 21: 24}},
	{"isz_v2.asm", map[int64]int64{19: 9, 20: 24}},
	{"isz_v3.asm", map[int64]int64{19: 9, 20: 24}},
	{"jsr_v1.asm", map[int64]int64{5: 50}},
	{"loopuntil_v1.asm", map[int64]int64{9: 5000}},
	{"loopuntil_v2.asm", map[int64]int64{8: 5000}},
	{"subleq_v1.asm", map[int64]int64{78: 5000}},
	{"subleq_v2.asm", map[int64]int64{77: 5000}},
	{"subleq_v3.asm", map[int64]int64{76: 5000}},
	{"switch_v1.asm", map[int64]int64{44: 2255}},
	{"switch_v2.asm", map[int64]int64{44: 2255}},
	{"switch_v3.asm", map[int64]int64{51: 2255}},
}

func TestRun(t *testing.T) {
	for _, test := range tests {
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
	for _, test := range tests {
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
