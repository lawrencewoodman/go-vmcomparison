package codegen

import (
	"math"
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
	// Memory locations
	const (
		m_memBase = 1
		m_opAddr  = 2
		m_lac     = 3
		m_ok      = 4
		m_val     = 5
	)
	mem := []uint{
		555, // Just to create an offset from 0
		1,   // memBase
		4,   // opAddr
		9,   // lac
		1,   // ok for HLT - TODO: using to indicate hlt
		23,  // value to add to AC
	}

	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opAddr)) },
		func(v *CGVM) { op_ADD(v, m_lac) },
		func(v *CGVM) { op_STA13(v, m_lac) },
		func(v *CGVM) { op_HLT(v, m_ok) },
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
	// Memory locations
	const (
		m_sum   = 1
		m_ok    = 2
		m_cnt   = 3
		m_l0    = 4
		m_l1    = 5
		m_l5000 = 6
	)
	mem := []uint{
		555,  // Just to create an offset from 0
		0,    // sum
		1,    // ok for HLT - TODO: using to indicate hlt
		0,    // cnt
		0,    // l0
		1,    // l1
		5000, // l5000
	}
	// Program locations
	const (
		p_loop = 3
	)
	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, m_l5000) },
		func(v *CGVM) { op_STA(v, m_cnt) },
		func(v *CGVM) { op_LDA(v, m_l0) },
		// loop:
		func(v *CGVM) { op_ADD(v, m_l1) },
		func(v *CGVM) { op_DSZ(v, m_cnt) },
		func(v *CGVM) { op_JMP(v, p_loop) },
		func(v *CGVM) { op_STA(v, m_sum) },
		func(v *CGVM) { op_HLT(v, m_ok) },
	}
	return mem, program
}

/*
; SUBLEQ emulator

; Fetch operands
fetch:  LDA II memBase,pc
STA     opA
INC     pc
LDA II memBase,pc
STA     opB
INC     pc
LDA II memBase,pc
STA     opC
INC     pc


; Execute
exec:   LDA II  memBase,opB
SUB II  memBase,opA
STA II  memBase,opB

; If opB == 1000 THEN halt
LDA     l1000
SUB     opB
JEQ     halt

; IF mem[opB] > 0 THEN jump to fetch
LDA II  memBase,opB
JGT     fetch

; ELSE jump to opC
jmpC:   LDA     opC
STA     pc
JMP     fetch


halt:   LDA II  memBase,opB
STA     hltVal
HLT     ok


l1000:   1000
pc:      0
hltVal:  0
memBase: program
opA:     0
opB:     0
opC:     0
ok:      0


; loopuntil_v1 from subleq/fixtures/
program:
15
13
3
16
14
6
16
13
3
16
1000
12
0
0
sum: 0
4999
-1
*/

// A compiled Subleq VM
func initSubleq() ([]uint, []func(*CGVM)) {
	// Memory locations
	const (
		m_l1000   = 0
		m_pc      = 1
		m_hltVal  = 2
		m_memBase = 3
		m_opA     = 4
		m_opB     = 5
		m_opC     = 6
		m_ok      = 7
		m_program = 8
	)

	mem := []uint{
		1000,      // l1000
		0,         // pc
		0,         // hltVal
		m_program, // memBase
		0,         // opA
		0,         // opB
		0,         // opC
		0,         // ok

		// loopuntil_v1 from subleq/fixtures/
		// programL:
		15,
		13,
		3,
		16,
		14,
		6,
		16,
		13,
		3,
		16,
		1000,
		12,
		0,
		0,
		0, // sum
		4999,
		math.MaxUint64, // -1

	}
	for i := 0; i < 1000; i++ {
		mem = append(mem, 0)
	}
	// Program locations
	const (
		p_fetch = 0
		p_halt  = 20
	)
	program := []func(v *CGVM){
		// Fetch operands
		// fetch:
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opA) },
		func(v *CGVM) { op_INC(v, m_pc) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opB) },
		func(v *CGVM) { op_INC(v, m_pc) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opC) },
		func(v *CGVM) { op_INC(v, m_pc) },

		// Execute
		// exec:
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_SUB(v, calcBaseIndexAddr(v, m_memBase, m_opA)) },
		func(v *CGVM) { op_STA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },

		// IF opB == 1000 THEN halt
		func(v *CGVM) { op_LDA(v, m_l1000) },
		func(v *CGVM) { op_SUB(v, m_opB) },
		func(v *CGVM) { op_JEQ(v, p_halt) },

		// IF mem[opB] > 0 THEN jump to fetch
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_JGT(v, p_fetch) },

		// ELSE jump to opC
		func(v *CGVM) { op_LDA(v, m_opC) },
		func(v *CGVM) { op_STA(v, m_pc) },
		func(v *CGVM) { op_JMP(v, p_fetch) },

		// halt:
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_STA(v, m_hltVal) },
		func(v *CGVM) { op_HLT(v, m_ok) },
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
	{"subleq", initSubleq, map[uint]uint{22: 5000}},
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
