package codegen

import (
	"testing"
)

/*
			; PDP-8 TAD
			LDA II memBase:opAddr
			ADD     lac
            STA13   lac
done:    	HLT     ok
memBase: 	4
opAddr: 	4
lac:    	9
ok:     	0
val:     	23
*/

// A compiled TAD
func initTAD() ([]uint, []func(*CGVM)) {
	mem := []uint{
		555, // Just to create an offset from 0
		1,   // memBase
		4,   // opAddr
		9,   // lac
		1,   // ok for HLT - TODO: using to indicate hlt
		23,  // value to add to AC
	}

	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, 1, 2)) },
		func(v *CGVM) { op_ADD(v, 3) },
		func(v *CGVM) { op_STA13(v, 3) },
		func(v *CGVM) { op_HLT(v, 4) },
	}
	return mem, program
}

/*
        ; Version 1
        ; LOOP UNTIL
        LDA     l150
        STA     cnt
        LDA     l0
loop:   ADD     l1
        DSZ     cnt
        JMP     loop
        STA     sum
done:   HLT     ok
sum:    0
ok:     0
cnt:    0
l0:     0
l1:     1
l150:   150
*/

// A compiled LoopUntil
func initLoopUntil() ([]uint, []func(*CGVM)) {
	mem := []uint{
		555,  // Just to create an offset from 0
		0,    // sum
		1,    // ok for HLT - TODO: using to indicate hlt
		0,    // cnt
		0,    // l0
		1,    // l1
		5000, // l5000
	}
	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, 6) },
		func(v *CGVM) { op_STA(v, 3) },
		func(v *CGVM) { op_LDA(v, 4) },
		func(v *CGVM) { op_ADD(v, 5) },
		func(v *CGVM) { op_DSZ(v, 3) },
		func(v *CGVM) { op_JMP(v, 3) },
		func(v *CGVM) { op_STA(v, 1) },
		func(v *CGVM) { op_HLT(v, 2) },
	}
	return mem, program
}

var tests = []struct {
	name string
	init func() ([]uint, []func(*CGVM))
	want map[uint]uint // [memloc]value
}{
	{"tad", initTAD, map[uint]uint{3: 32}},
	{"loopuntil", initLoopUntil, map[uint]uint{1: 5000}},
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
