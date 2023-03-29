package vm2

import (
	"path/filepath"
	"testing"
)

var tests = []struct {
	filename string
	want     map[uint]uint // [memloc]value
}{
	{"jsr_v1.asm", map[uint]uint{12: 607500, 27: 4050}},
	{"jsr_v2.asm", map[uint]uint{12: 607500, 24: 4050}},
	{"loopuntil_v1.asm", map[uint]uint{16: 3459}},
	{"switch_v1.asm", map[uint]uint{51: 2255}},
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
		b.Run(test.filename, func(b *testing.B) {
			b.StopTimer()

			routine, err := asm(filepath.Join("fixtures", test.filename))
			if err != nil {
				b.Fatalf("asm() err: %v", err)
			}

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
	}
}
