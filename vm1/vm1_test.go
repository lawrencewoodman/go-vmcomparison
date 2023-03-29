package vm1

import (
	"path/filepath"
	"testing"
)

var VMtests = []struct {
	filename string
	want     map[uint]uint // [memloc]value
}{
	{"tad_v1.asm", map[uint]uint{12: 32}},
	{"tad_v2.asm", map[uint]uint{10: 32}},
	{"tad_v3.asm", map[uint]uint{6: 32}},
	{"isz_v1.asm", map[uint]uint{19: 24, 17: 9}},
	{"isz_v2.asm", map[uint]uint{12: 24, 10: 9}},
	{"isz_v3.asm", map[uint]uint{9: 24, 7: 9}},
	{"loopuntil_v1.asm", map[uint]uint{13: 3459}},
	{"loopuntil_v2.asm", map[uint]uint{10: 3459}},
	{"loopuntil_v3.asm", map[uint]uint{13: 3459}},
	{"loopuntil_v4.asm", map[uint]uint{11: 3459}},
	{"loopuntil_v5.asm", map[uint]uint{9: 3459}},
	{"switch_v1.asm", map[uint]uint{44: 2255}},
	{"jsr_v1.asm", map[uint]uint{8: 607500, 22: 4050}},
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
