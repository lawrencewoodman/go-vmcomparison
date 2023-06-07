package codegen

import (
	"testing"
)

type Test struct {
	name string
	init func() ([]uint, []func(*CGVM))
	want map[uint]uint // [memloc]value
}

var tests = []Test{}

func addTest(name string, init func() ([]uint, []func(*CGVM)), want map[uint]uint) {
	tests = append(tests, Test{name, init, want})
}

func TestRun(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem, program := test.init()
			v := New()
			v.LoadMem(mem)
			v.Run(program)
			for memLoc, wantValue := range test.want {
				if v.mem[memLoc] != wantValue {
					t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
				}
			}
		})
	}
}

func BenchmarkRun(t *testing.B) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.B) {
			t.StopTimer()
			mem, program := test.init()
			for n := 0; n < t.N; n++ {
				v := New()
				v.LoadMem(mem)

				t.StartTimer()
				v.Run(program)
				t.StopTimer()

				for memLoc, wantValue := range test.want {
					if v.mem[memLoc] != wantValue {
						t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
					}
				}
			}
		})
	}
}
